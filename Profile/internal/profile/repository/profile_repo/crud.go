package profilerepo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	dto "github.com/YoungFlores/Case_Go/Profile/internal/profile/dto"
	"github.com/YoungFlores/Case_Go/Profile/internal/profile/models"
	repoerr "github.com/YoungFlores/Case_Go/Profile/internal/profile/repository/errors"
	"github.com/lib/pq"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

// --- Create Methods ---

func (r *PostgresProfileRepo) CreateProfile(ctx context.Context, profile *models.Profile) (*models.Profile, error) {
	now := time.Now()
	query := psql.Insert("profiles").
		Columns("user_id", "avatar", "is_active", "description", "username", "name", "surname", "patronymic", "email", "phone_number", "sex", "profession", "case_count", "created_at", "updated_at").
		Values(profile.UserID, profile.Avatar, profile.IsActive, profile.Description, profile.Username, profile.Name, profile.Surname, profile.Patronymic, profile.Email, profile.PhoneNumber, profile.Sex, profile.Profession, profile.CaseCount, now, now).
		Suffix("RETURNING id, created_at, updated_at")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRowContext(ctx, sql, args...).Scan(&profile.ID, &profile.CreatedAt, &profile.UpdatedAt)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return nil, &repoerr.RepoError{
					Field: extractField(pgErr.Constraint),
					Err:   repoerr.ErrConflict,
				}
			}
		}
		return nil, fmt.Errorf("scan profile: %w", err)
	}

	return profile, nil
}

func (r *PostgresProfileRepo) AddSocial(ctx context.Context, links []models.UserSocialLink) ([]models.UserSocialLink, error) {
	// TODO: return err
	if len(links) == 0 {
		return links, nil
	}

	query := psql.Insert("user_social_links").Columns("user_id", "type", "url")
	for _, link := range links {
		query = query.Values(link.UserID, link.Type, link.URL)
	}
	query = query.Suffix("RETURNING id")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	for i := range ids {
		links[i].ID = ids[i]
	}

	return links, nil
}

func (r *PostgresProfileRepo) AddPurposes(ctx context.Context, purposes []models.UserPurpose) ([]models.UserPurpose, error) {
	if len(purposes) == 0 {
		return purposes, nil
	}

	query := psql.Insert("user_purposes").Columns("user_id", "purpose")
	for _, p := range purposes {
		query = query.Values(p.UserID, p.Purpose)
	}
	query = query.Suffix("RETURNING id")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for i := 0; rows.Next(); i++ {
		if err := rows.Scan(&purposes[i].ID); err != nil {
			return nil, err
		}
	}

	return purposes, nil
}

// --- Get Methods ---

func (r *PostgresProfileRepo) GetProfileByID(ctx context.Context, id int64) (*models.Profile, error) {
	query := psql.Select("*").From("profiles").Where(sq.Eq{"id": id})
	return r.fetchProfile(ctx, query)
}

func (r *PostgresProfileRepo) GetUserProfile(ctx context.Context, userID int64) (*models.Profile, error) {
	query := psql.Select("*").From("profiles").Where(sq.Eq{"user_id": userID})
	return r.fetchProfile(ctx, query)
}

func (r *PostgresProfileRepo) GetUserByProfileID(ctx context.Context, id, userID int64) (int64, error) {
	sql, args, err := psql.Select("user_id").From("profiles").Where(sq.Eq{"id": id, "user_id": userID}).ToSql()
	if err != nil {
		return 0, err
	}

	var resID int64
	err = r.db.QueryRowContext(ctx, sql, args...).Scan(&resID)
	return resID, err
}

func (r *PostgresProfileRepo) GetUserSocials(ctx context.Context, userID int64) ([]models.UserSocialLink, error) {
	sql, args, err := psql.Select("id", "user_id", "type", "url").From("user_social_links").Where(sq.Eq{"user_id": userID}).ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []models.UserSocialLink
	for rows.Next() {
		var l models.UserSocialLink
		if err := rows.Scan(&l.ID, &l.UserID, &l.Type, &l.URL); err != nil {
			return nil, err
		}
		links = append(links, l)
	}
	return links, nil
}

func (r *PostgresProfileRepo) GetUserPurposes(ctx context.Context, userID int64) ([]models.UserPurpose, error) {
	sql, args, err := psql.Select("id", "user_id", "purpose").From("user_purposes").Where(sq.Eq{"user_id": userID}).ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var purposes []models.UserPurpose
	for rows.Next() {
		var p models.UserPurpose
		if err := rows.Scan(&p.ID, &p.UserID, &p.Purpose); err != nil {
			return nil, err
		}
		purposes = append(purposes, p)
	}
	return purposes, nil
}

func (r *PostgresProfileRepo) GetAllUsers(ctx context.Context, limit int) ([]models.Profile, error) {
	sql, args, err := psql.Select("*").From("profiles").Where(sq.Eq{"is_active": true}).Limit(uint64(limit)).ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []models.Profile
	for rows.Next() {
		var p models.Profile
		err := rows.Scan(&p.ID, &p.UserID, &p.Avatar, &p.IsActive, &p.Description, &p.Username, &p.Name, &p.Surname, &p.Patronymic, &p.Email, &p.PhoneNumber, &p.Sex, &p.Profession, &p.CaseCount, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, p)
	}
	return profiles, nil
}

// --- Update Methods ---

func (r *PostgresProfileRepo) UpdateProfile(ctx context.Context, profile *models.Profile) (*models.Profile, error) {
	sql, args, err := psql.Update("profiles").
		Set("avatar", profile.Avatar).
		Set("description", profile.Description).
		Set("username", profile.Username).
		Set("name", profile.Name).
		Set("surname", profile.Surname).
		Set("patronymic", profile.Patronymic).
		Set("email", profile.Email).
		Set("phone_number", profile.PhoneNumber).
		Set("sex", profile.Sex).
		Set("profession", profile.Profession).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"user_id": profile.UserID}).
		Suffix("RETURNING updated_at").
		ToSql()

	if err != nil {
		return nil, err
	}

	err = r.db.QueryRowContext(ctx, sql, args...).Scan(&profile.UpdatedAt)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" {
				return nil, &repoerr.RepoError{
					Field: extractField(pgErr.Constraint),
					Err:   repoerr.ErrConflict,
				}
			}
		}
		return nil, fmt.Errorf("scan profile: %w", err)
	}
	return profile, err
}

func (r *PostgresProfileRepo) PathcProfile(ctx context.Context, userID int64, updates dto.UpdateProfilePartialDTO) (*models.Profile, error) {
	query := psql.Update("profiles").Where(sq.Eq{"user_id": userID}).Set("updated_at", time.Now())

	if updates.Avatar != nil {
		query = query.Set("avatar", *updates.Avatar)
	}
	if updates.Username != nil {
		query = query.Set("username", *updates.Username)
	}
	if updates.Name != nil {
		query = query.Set("name", *updates.Name)
	}
	if updates.Surname != nil {
		query = query.Set("surname", *updates.Surname)
	}
	if updates.Patronymic != nil {
		query = query.Set("patronymic", updates.Patronymic)
	}
	if updates.Email != nil {
		query = query.Set("email", *updates.Email)
	}
	if updates.PhoneNumber != nil {
		query = query.Set("phone_number", updates.PhoneNumber)
	}
	if updates.Sex != nil {
		query = query.Set("sex", updates.Sex)
	}
	if updates.Description != nil {
		query = query.Set("description", *updates.Description)
	}
	if updates.Profession != nil {
		query = query.Set("profession", updates.Profession)
	}

	sql, args, err := query.Suffix("RETURNING id, user_id, avatar, is_active, description, username, name, surname, patronymic, email, phone_number, sex, profession, case_count, created_at, updated_at").ToSql()
	if err != nil {
		return nil, err
	}

	var p models.Profile
	err = r.db.QueryRowContext(ctx, sql, args...).Scan(&p.ID, &p.UserID, &p.Avatar, &p.IsActive, &p.Description, &p.Username, &p.Name, &p.Surname, &p.Patronymic, &p.Email, &p.PhoneNumber, &p.Sex, &p.Profession, &p.CaseCount, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" {
				return nil, &repoerr.RepoError{
					Field: extractField(pgErr.Constraint),
					Err:   repoerr.ErrConflict,
				}
			}
		}
		return nil, fmt.Errorf("scan profile: %w", err)
	}

	return &p, err
}

func (r *PostgresProfileRepo) UpdateLinks(ctx context.Context, links []models.UserSocialLink) ([]models.UserSocialLink, error) {
	for i, link := range links {
		sql, args, err := psql.Update("user_social_links").
			Set("type", link.Type).
			Set("url", link.URL).
			Where(sq.Eq{"id": link.ID, "user_id": link.UserID}).
			ToSql()
		if err != nil {
			return nil, err
		}

		res, err := r.db.ExecContext(ctx, sql, args...)
		if err != nil {
			return nil, err
		}

		rows, err := res.RowsAffected()
		if err != nil {
			return nil, err
		}

		if rows == 0 {
			return nil, &repoerr.RepoError{
				Err: repoerr.ErrFrobidden,
			}
		}

		links[i] = link
	}
	return links, nil
}

func (r *PostgresProfileRepo) EditSocial(ctx context.Context, link *models.UserSocialLink) ([]models.UserSocialLink, error) {
	_, err := r.UpdateLinks(ctx, []models.UserSocialLink{*link})
	if err != nil {
		return nil, err
	}
	return r.GetUserSocials(ctx, link.UserID)
}

func (r *PostgresProfileRepo) UpdatePurposes(ctx context.Context, purposes []models.UserPurpose) ([]models.UserPurpose, error) {
	for i, p := range purposes {
		sql, args, err := psql.Update("user_purposes").
			Set("purpose", p.Purpose).
			Where(sq.Eq{"id": p.ID, "user_id": p.UserID}).
			ToSql()
		if err != nil {
			return nil, err
		}

		res, err := r.db.ExecContext(ctx, sql, args...)
		if err != nil {
			return nil, fmt.Errorf("exec update purpose: %w", err)
		}

		rows, err := res.RowsAffected()
		if err != nil {
			return nil, err
		}
		if rows == 0 {
			return nil, &repoerr.RepoError{
				Field: "id",
				Err:   repoerr.ErrFrobidden,
			}
		}

		purposes[i] = p
	}
	return purposes, nil
}

func (r *PostgresProfileRepo) EditPurpose(ctx context.Context, purpose *models.UserPurpose) ([]models.UserPurpose, error) {
	_, err := r.UpdatePurposes(ctx, []models.UserPurpose{*purpose})
	if err != nil {
		return nil, err
	}
	return r.GetUserPurposes(ctx, purpose.UserID)
}

// --- Delete Methods ---

func (r *PostgresProfileRepo) DeletePupose(ctx context.Context, id int64) error {
	sql, args, err := psql.Delete("user_purposes").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, sql, args...)
	return err
}

func (r *PostgresProfileRepo) DeleteSocial(ctx context.Context, id int64) error {
	sql, args, err := psql.Delete("user_social_links").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, sql, args...)
	return err
}

func (r *PostgresProfileRepo) DeleteProfile(ctx context.Context, userID int64) error {
	query := psql.Update("profiles").
		Set("is_active", false).
		Where(sq.Eq{"user_id": userID}).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	res, err := r.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	log.Printf("Rows affected: %d", rows)

	return nil
}

func (r *PostgresProfileRepo) DeleteProfileWithoutRecovery(ctx context.Context, userID int64) error {
	// Хард удаление
	sql, args, err := psql.Delete("profiles").Where(sq.Eq{"user_id": userID}).ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, sql, args...)
	return err
}
