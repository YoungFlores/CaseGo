package searchRepo

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/YoungFlores/Case_Go/Profile/internal/profile/models"
	"github.com/YoungFlores/Case_Go/Profile/internal/search/dto"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

func (r *PostgresSearchRepo) SearchProfile(ctx context.Context, req dto.SearchDTO, limit, offset uint64, sortBy, sortOrder string) ([]models.Profile, error) {
	query := psql.Select(
		"id", "user_id", "avatar", "is_active", "description",
		"username", "name", "surname", "patronymic", "city",
		"age", "sex", "profession", "case_count", "created_at", "updated_at",
	).From("profiles")

	conditions := sq.And{}

	conditions = append(conditions, sq.Eq{"is_active": true})

	if req.ProfessionID != nil {
		conditions = append(conditions, sq.Eq{"profession_id": req.ProfessionID})
	}

	if req.Profession != nil && req.ProfessionID == nil {
		conditions = append(conditions, sq.Eq{"profession": req.Profession})
	}

	if req.MinAge != nil && req.MaxAge != nil {
		conditions = append(conditions, sq.GtOrEq{"age": req.MinAge}, sq.LtOrEq{"age": req.MaxAge})
	}

	if req.City != nil {
		conditions = append(conditions, sq.Eq{"city": req.City})
	}

	if req.Sex != nil {
		conditions = append(conditions, sq.Eq{"sex": req.Sex})
	}

	query = query.Where(conditions)

	query = query.Limit(limit).Offset(offset)
	query = query.OrderBy(sortBy + " " + sortOrder)

	strSql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, strSql, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var profiles []models.Profile
	for rows.Next() {
		var p models.Profile

		err := rows.Scan(
			&p.ID,
			&p.UserID,
			&p.Avatar,
			&p.IsActive,
			&p.Description,
			&p.Username,
			&p.Name,
			&p.Surname,
			&p.Patronymic,
			&p.City,
			&p.Age,
			&p.Sex,
			&p.Profession,
			&p.CaseCount,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return profiles, nil

}

func (r *PostgresSearchRepo) SearchByFio(ctx context.Context, req dto.SearchByFIODTO, limit, offset uint64) ([]models.Profile, error) {
	query := psql.Select(
		"id", "user_id", "avatar", "is_active", "description",
		"username", "name", "surname", "patronymic", "city",
		"age", "sex", "profession", "case_count", "created_at", "updated_at",
	).From("profiles")

	conditions := sq.And{}

	if req.Name != nil {
		conditions = append(conditions, sq.ILike{"name": "%" + *req.Name + "%"})
	}

	if req.Surname != nil {
		conditions = append(conditions, sq.ILike{"surname": "%" + *req.Surname + "%"})
	}

	if req.Patronymic != nil {
		conditions = append(conditions, sq.ILike{"patronymic": "%" + *req.Patronymic + "%"})
	}

	if len(conditions) == 0 {
		return nil, fmt.Errorf("no conditions")
	}
	conditions = append(conditions, sq.Eq{"is_active": true})
	query = query.Where(conditions)
	query = query.Limit(limit).Offset(offset)
	query = query.OrderBy("name ASC")

	strSql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, strSql, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var profiles []models.Profile
	for rows.Next() {
		var p models.Profile

		err := rows.Scan(
			&p.ID,
			&p.UserID,
			&p.Avatar,
			&p.IsActive,
			&p.Description,
			&p.Username,
			&p.Name,
			&p.Surname,
			&p.Patronymic,
			&p.City,
			&p.Age,
			&p.Sex,
			&p.Profession,
			&p.CaseCount,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return profiles, nil

}
