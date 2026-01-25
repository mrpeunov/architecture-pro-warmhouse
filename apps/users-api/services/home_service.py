from datetime import datetime
from typing import List
import uuid

import asyncpg

from models.home import Home, HomeRequest
from repositories.home_repository import HomeRepository
from repositories.user_repository import UserRepository


class HomeService:
    def __init__(self, conn: asyncpg.Connection):
        self.home_repository = HomeRepository(conn)
        self.user_repository = UserRepository(conn)

    async def create_home(self, home_data: HomeRequest) -> Home:
        """Create a new home"""
        home = Home(
            home_id=uuid.uuid4(),
            address=home_data.address,
            created_at=datetime.now()
        )
        
        return await self.home_repository.create(home)

    async def get_homes(self, email: str) -> list[Home]:
        """Get homes by user ID"""
        return await self.home_repository.find_by_email(email)

    async def add_user_to_home(self, home_id: uuid.UUID, email: str) -> bool:
        """Add a user to a home"""
        try:
            # Check if user exists
            user = await self.user_repository.find_by_email(email)
            if not user:
                return False
            
            # Add user to home
            await self.home_repository.add_user_to_home(email, home_id)
            return True
        except Exception:
            return False
