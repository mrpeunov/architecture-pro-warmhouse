from datetime import datetime
from typing import List
from uuid import UUID

from pydantic import BaseModel

from models.user import AuthToken


class Home(BaseModel):
    home_id: UUID
    address: str
    created_at: datetime


class HomeRequest(BaseModel):
    address: str


class HomeResponse(BaseModel):
    home_id: UUID
    address: str
    created_at: datetime
    auth_token: AuthToken