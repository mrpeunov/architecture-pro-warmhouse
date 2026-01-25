import asyncpg
from fastapi import HTTPException, status
from models.user import UserRequest, UserResponse, LoginRequest, AuthToken
from services import HomeService
from services.user_service import UserService


class AuthController:
    def __init__(self, conn: asyncpg.Connection):
        self.user_service = UserService(conn)
        self.home_service = HomeService(conn)

    async def register(self, user_data: UserRequest) -> UserResponse:
        """Register a new user"""
        try:
            user = await self.user_service.create_user(user_data)
            return UserResponse(
                email=user.email,
                name=user.name,
                created_at=user.created_at
            )
        except ValueError as e:
            raise HTTPException(
                status_code=status.HTTP_400_BAD_REQUEST,
                detail=str(e)
            )
        except Exception as e:
            raise HTTPException(
                status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
                detail="Internal server error"
            )

    async def login(self, credentials: LoginRequest) -> AuthToken:
        """Login user and return access token"""
        user = await self.user_service.authenticate(
            str(credentials.email),
            credentials.password
        )
        
        if not user:
            raise HTTPException(
                status_code=status.HTTP_401_UNAUTHORIZED,
                detail="Incorrect email or password",
                headers={"WWW-Authenticate": "Bearer"},
            )

        homes = await self.home_service.get_homes(user.email)
        access_token = self.user_service.create_token(user, homes)
        return AuthToken(access_token=access_token, token_type="Bearer")
