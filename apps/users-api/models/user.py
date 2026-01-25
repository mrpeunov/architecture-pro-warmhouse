from datetime import datetime
from typing import Optional
from pydantic import BaseModel, EmailStr


class User(BaseModel):
    email: str
    password: str
    name: str
    created_at: datetime


class UserRequest(BaseModel):
    email: EmailStr
    password: str
    name: str


class UserResponse(BaseModel):
    email: str
    name: str
    created_at: datetime


class LoginRequest(BaseModel):
    email: EmailStr
    password: str


class AuthToken(BaseModel):
    access_token: str
    token_type: str
