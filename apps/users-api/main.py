import asyncio
from http.cookiejar import debug

from fastapi import FastAPI, Depends, HTTPException, status
from fastapi.security import HTTPBearer, HTTPAuthorizationCredentials
from jose import JWTError, jwt
from typing import List
import os

from models.user import UserRequest, UserResponse, LoginRequest, AuthToken
from models.home import HomeRequest, HomeResponse
from controllers.auth_controller import AuthController
from controllers.home_controller import HomeController
from database.connection import init_database, close_database, get_pool

SECRET_KEY = os.getenv("SECRET_KEY", "secret-key")
ALGORITHM = "HS256"

app = FastAPI(title="Users API", version="1.0.0", debug=True)
security = HTTPBearer()


async def get_current_user_email(credentials: HTTPAuthorizationCredentials = Depends(security)) -> str:
    """Extract user email from JWT token"""
    try:
        payload = jwt.decode(credentials.credentials, SECRET_KEY, algorithms=[ALGORITHM])
        email: str = payload.get("sub")
        if email is None:
            raise HTTPException(
                status_code=status.HTTP_401_UNAUTHORIZED,
                detail="Could not validate credentials",
                headers={"WWW-Authenticate": "Bearer"},
            )
        return email
    except JWTError:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Could not validate credentials",
            headers={"WWW-Authenticate": "Bearer"},
        )


@app.on_event("startup")
async def startup_event():
    """Initialize database connection on startup"""
    await init_database()


@app.on_event("shutdown")
async def shutdown_event():
    """Close database connection on shutdown"""
    await close_database()


@app.post("/auth/register", response_model=UserResponse)
async def register(user_data: UserRequest):
    """Register a new user"""
    pool = await get_pool()
    async with pool.acquire() as conn:
        auth_controller = AuthController(conn)
        return await auth_controller.register(user_data)


@app.post("/auth/login", response_model=AuthToken)
async def login(credentials: LoginRequest):
    """Login user and get access token"""
    pool = await get_pool()
    async with pool.acquire() as conn:
        auth_controller = AuthController(conn)
        return await auth_controller.login(credentials)


@app.post("/homes", response_model=HomeResponse)
async def create_home(home_data: HomeRequest, current_user_email: str = Depends(get_current_user_email)):
    """Create a new home"""
    pool = await get_pool()
    async with pool.acquire() as conn:
        home_controller = HomeController(conn)
        return await home_controller.create_home(home_data, current_user_email)


@app.get("/homes", response_model=List[HomeResponse])
async def get_user_homes(current_user_email: str = Depends(get_current_user_email)):
    """Get homes for a specific user"""
    pool = await get_pool()
    async with pool.acquire() as conn:
        home_controller = HomeController(conn)
        return await home_controller.get_user_homes(current_user_email)


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)
