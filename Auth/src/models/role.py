import enum
from datetime import datetime

from sqlalchemy import ForeignKey, UniqueConstraint, func
from sqlalchemy.orm import Mapped, mapped_column
from .base import MyBaseModel
from sqlalchemy.dialects.postgresql import ENUM as PG_ENUM


class RoleEnum(str, enum.Enum):
	USER = "user"
	MODERATOR = "moderator"
	ADMIN = "admin"
	SUPER_ADMIN = "super_admin"


class UserRole(MyBaseModel):
	"""Связь пользователя с ролями"""
	__tablename__ = "user_roles"

	id: Mapped[int] = mapped_column(primary_key=True, autoincrement=True)
	user_id: Mapped[int] = mapped_column(
		ForeignKey("users.id", ondelete="CASCADE"),
		nullable=False
	)

	role: Mapped[RoleEnum] = mapped_column(
		PG_ENUM(RoleEnum, name="user_role_enum"),
		nullable=False,
		default=RoleEnum.USER,
		server_default="USER"
	)

	granted_at: Mapped[datetime] = mapped_column(server_default=func.now())
	granted_by: Mapped[int | None] = mapped_column(
		ForeignKey("users.id", ondelete="SET NULL"),
		nullable=True
	)

	__table_args__ = (
		UniqueConstraint("user_id", "role", name="uq_user_role_unique"),
	)
