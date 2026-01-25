from typing import Optional

import asyncpg

from jose import jwt
from datetime import datetime, timedelta
import os

from models import Home
from models.user import User, UserRequest
from repositories.user_repository import UserRepository

# JWT settings
SECRET_KEY = os.getenv("SECRET_KEY", "secret-key")
ALGORITHM = "HS256"
ACCESS_TOKEN_EXPIRE_MINUTES = 30


class UserService:
    def __init__(self, conn: asyncpg.Connection):
        self.user_repository = UserRepository(conn)

    def _create_access_token(self, data: dict, expires_delta: Optional[timedelta] = None):
        """Create a JWT access token"""
        to_encode = data.copy()
        if expires_delta:
            expire = datetime.now() + expires_delta
        else:
            expire = datetime.now() + timedelta(minutes=15)
        to_encode.update({"exp": expire})
        encoded_jwt = jwt.encode(to_encode, SECRET_KEY, algorithm=ALGORITHM)
        return encoded_jwt

    async def create_user(self, user_data: UserRequest) -> User:
        """Create a new user"""
        # Check if user already exists
        existing_user = await self.user_repository.find_by_email(str(user_data.email))
        if existing_user:
            raise ValueError("User with this email already exists")
        
        # Create user model
        user = User(
            email=str(user_data.email),
            password=user_data.password,
            name=user_data.name,
            created_at=datetime.utcnow()
        )
        
        # Save to database
        return await self.user_repository.create(user)

    async def authenticate(self, email: str, password: str) -> Optional[User]:
        """Authenticate a user"""
        user = await self.user_repository.find_by_email(email)
        if not user:
            return None
        
        if password != user.password:
            return None
        
        return user

    async def get_user_by_email(self, email: str) -> Optional[User]:
        """Get user by email"""
        return await self.user_repository.find_by_email(email)

    def create_token(self, user: User, homes: list[Home]) -> str:
        """Create access token for user"""
        access_token_expires = timedelta(minutes=ACCESS_TOKEN_EXPIRE_MINUTES)
        access_token = self._create_access_token(
            data={"sub": user.email, "homes": [str(h.home_id) for h in homes]},
            expires_delta=access_token_expires
        )
        return access_token
