import asyncpg
from fastapi import HTTPException, status
from typing import List
from models.home import HomeRequest, HomeResponse
from models.user import AuthToken
from services import UserService
from services.home_service import HomeService


class HomeController:
    def __init__(self, conn: asyncpg.Connection):
        self.home_service = HomeService(conn)
        self.user_service = UserService(conn)

    async def create_home(self, home_data: HomeRequest, email: str) -> HomeResponse:
        """Create a new home"""
        try:
            home = await self.home_service.create_home(home_data)
            await self.home_service.add_user_to_home(home.home_id, email)

            homes = await self.home_service.get_homes(email)

            user = await self.user_service.get_user_by_email(email)
            access_token = self.user_service.create_token(user, homes)

            return HomeResponse(
                home_id=home.home_id,
                address=home.address,
                created_at=home.created_at,
                auth_token=AuthToken(access_token=access_token, token_type="Bearer"),
            )
        except Exception as e:
            raise HTTPException(
                status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
                detail="Internal server error"
            )

    async def get_user_homes(self, email: str) -> List[HomeResponse]:
        """Get homes for a specific user"""
        try:
            homes = await self.home_service.get_homes(email)
            user = await self.user_service.get_user_by_email(email)
            access_token = self.user_service.create_token(user, homes)
            return [
                HomeResponse(
                    home_id=home.home_id,
                    address=home.address,
                    created_at=home.created_at,
                    auth_token=AuthToken(access_token=access_token, token_type="Bearer"),
                )
                for home in homes
            ]
        except Exception:
            raise HTTPException(
                status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
                detail="Internal server error"
            )
