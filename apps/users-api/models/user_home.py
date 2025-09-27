from uuid import UUID

from pydantic import BaseModel


class UserHome(BaseModel):
    email: str
    home_id: UUID
