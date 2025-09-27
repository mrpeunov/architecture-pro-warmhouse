from uuid import UUID

import asyncpg
from models.home import Home

class HomeRepository:
    def __init__(self, conn: asyncpg.Connection):
        self.conn = conn

    async def create(self, home: Home) -> Home:
        query = """
            INSERT INTO homes (address, created_at)
            VALUES ($1, $2)
            RETURNING home_id, address, created_at
        """
        row = await self.conn.fetchrow(
            query,
            home.address,
            home.created_at
        )
        return Home(
            home_id=row['home_id'],
            address=row['address'],
            created_at=row['created_at']
        )

    async def find_by_email(self, email: str) -> list[Home]:
        query = """
            SELECT h.home_id, h.address, h.created_at
            FROM homes h
            JOIN user_homes uh ON h.home_id = uh.home_id
            JOIN users u ON uh.email = u.email
            WHERE uh.email = $1
        """
        rows = await self.conn.fetch(query, email)
        return [
            Home(
                home_id=row['home_id'],
                address=row['address'],
                created_at=row['created_at']
            )
            for row in rows
        ]


    async def add_user_to_home(self, email: str, home_id: UUID) -> None:
        query = """
            INSERT INTO user_homes (email, home_id)
            VALUES ($1, $2)
            ON CONFLICT (email, home_id) DO NOTHING
        """
        await self.conn.execute(query, email, home_id)