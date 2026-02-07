from datetime import datetime
from typing import Optional

from sqlalchemy import (
	String,
	Boolean,
	ForeignKey,
	Index,
	func,
)
from sqlalchemy.orm import Mapped, mapped_column
from sqlalchemy.dialects.postgresql import TIMESTAMP

from .base import MyBaseModel


class RefreshToken(MyBaseModel):
	"""
	Refresh-токены.
	"""
	__tablename__ = "refresh_tokens"

	id: Mapped[int] = mapped_column(
		primary_key=True,
		autoincrement=True
	)

	user_id: Mapped[int] = mapped_column(
		ForeignKey("users.id", ondelete="CASCADE"),
		nullable=False,
		index=True
	)

	token_hash: Mapped[str] = mapped_column(
		String(128),
		nullable=False,
		unique=True,
		index=True
	)

	is_revoked: Mapped[bool] = mapped_column(
		Boolean,
		default=False,
		nullable=False,
		server_default="false"
	)

	expires_at: Mapped[datetime] = mapped_column(
		TIMESTAMP(timezone=True),
		nullable=False,
		index=True
	)

	created_at: Mapped[datetime] = mapped_column(
		TIMESTAMP(timezone=True),
		server_default=func.now(),
		nullable=False
	)

	revoked_at: Mapped[Optional[datetime]] = mapped_column(
		TIMESTAMP(timezone=True),
		nullable=True
	)

	user_agent: Mapped[Optional[str]] = mapped_column(
		String(255),
		nullable=True
	)

	ip_address: Mapped[Optional[str]] = mapped_column(
		String(45),
		nullable=True
	)

	__table_args__ = (
		Index("idx_refresh_user_active", "user_id", "is_revoked"),
	)
