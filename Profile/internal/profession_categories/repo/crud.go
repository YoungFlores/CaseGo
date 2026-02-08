package repo

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/YoungFlores/Case_Go/Profile/internal/profession_categories/models"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

func (r *PostgresCategoryRepo) CreateCategory(ctx context.Context, category *models.Category) (*models.Category, error) {
	query := psql.Insert("categories").Columns("name", "parent_id").Values(category.Name, category.ParentID)

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRowContext(ctx, sqlStr, args...).Scan(&category.ID, &category.Name, &category.ParentID)
	return category, err
}

func (r *PostgresCategoryRepo) GetCategories(ctx context.Context) ([]models.Category, error) {
	query := psql.Select("id", "name", "parent_id").From("categories")

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.ParentID); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func (r *PostgresCategoryRepo) GetCategoryByID(ctx context.Context, id int16) (*models.Category, error) {
	query := psql.Select("id", "name", "parent_id").From("categories").Where(sq.Eq{"id": id})

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var category models.Category
	err = r.db.QueryRowContext(ctx, sqlStr, args...).Scan(&category.ID, &category.Name, &category.ParentID)
	return &category, err
}

func (r *PostgresCategoryRepo) GetCategoriesByParent(ctx context.Context, parentID int16) ([]models.Category, error) {
	query := psql.Select("id", "name", "parent_id").From("categories").Where(sq.Eq{"parent_id": parentID})

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var categories []models.Category
	if err = r.db.QueryRowContext(ctx, sqlStr, args...).Scan(&categories); err != nil {
		return nil, err
	}
	return categories, nil
}
