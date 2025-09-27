import asyncpg
from typing import Optional
from datetime import datetime
from models.user import User
from database.connection import get_pool


class UserRepository:
    def __init__(self, conn: asyncpg.Connection):
        self.conn = conn

    async def create(self, user: User) -> User:
        query = """
            INSERT INTO users (email, password, name, created_at)
            VALUES ($1, $2, $3, $4)
            RETURNING email, password, name, created_at
        """
        row = await self.conn.fetchrow(
            query,
            user.email,
            user.password,
            user.name,
            user.created_at
        )
        return User(
            email=row['email'],
            password=row['password'],
            name=row['name'],
            created_at=row['created_at']
        )

    async def find_by_email(self, email: str) -> Optional[User]:
        query = """
            SELECT email, password, name, created_at
            FROM users
            WHERE email = $1
        """
        row = await self.conn.fetchrow(query, email)
        if row:
            return User(
                email=row['email'],
                password=row['password'],
                name=row['name'],
                created_at=row['created_at']
            )
        return None

