import asyncpg
import os
from typing import Optional

DATABASE_URL = os.getenv("DATABASE_URL", "postgresql://postgres:postgres@localhost:5432/userdb")
_pool: Optional[asyncpg.Pool] = None


async def get_pool() -> asyncpg.Pool:
    """Get a database connection from the pool"""
    global _pool
    if _pool is None:
        await init_database()
    return _pool


async def init_database():
    """Initialize the database connection pool"""
    global _pool
    if _pool is None:
        _pool = await asyncpg.create_pool(
            DATABASE_URL,
            min_size=1,
            max_size=10,
            command_timeout=60
        )


async def close_database():
    """Close the database connection pool"""
    global _pool
    if _pool:
        await _pool.close()
        _pool = None
