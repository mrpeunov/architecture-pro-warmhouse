import asyncpg
import os
from typing import Optional


def build_database_url() -> str:
    """Build database URL from individual environment variables"""
    db_user = os.getenv("DB_USER", "postgres")
    db_password = os.getenv("DB_PASSWORD", "postgres")
    db_host = os.getenv("DB_HOST", "localhost")
    db_port = os.getenv("DB_PORT", "5432")
    db_name = os.getenv("DB_NAME", "userdb")
    
    return f"postgresql://{db_user}:{db_password}@{db_host}:{db_port}/{db_name}"

DATABASE_URL = build_database_url()
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
