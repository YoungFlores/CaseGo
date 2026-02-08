package searchRepo

import (
	"context"
	"database/sql"

	"github.com/YoungFlores/Case_Go/Profile/internal/profile/models"
	"github.com/YoungFlores/Case_Go/Profile/internal/search/dto"
)

type SearchRepo interface {
	SearchProfile(ctx context.Context, req dto.SearchDTO, limit, offset uint64, sortBy, sortOrder string) ([]models.Profile, error)
	SearchByFio(ctx context.Context, req dto.SearchByFIODTO, limit, offset uint64) ([]models.Profile, error)
}

type PostgresSearchRepo struct {
	db *sql.DB
}

func NewPostgresSearchRepo(db *sql.DB) *PostgresSearchRepo {
	return &PostgresSearchRepo{db: db}
}
