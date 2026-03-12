from typing import Annotated

from fastapi import APIRouter, Depends

from ...dependencies import get_user_by_token
from models.user import User
from schemas.user import UserResponse

router = APIRouter()


@router.get("/me", response_model=UserResponse)
async def read_users_me_endpoint(
    current_user: Annotated[User, Depends(get_user_by_token)]
):
    return current_user
