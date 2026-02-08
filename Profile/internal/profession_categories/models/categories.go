package models

type Category struct {
	ID       int16  `json:"id" db:"id"`
	ParentID *int16 `json:"parent_id" db:"parent_id"`
	Name     string `json:"name" db:"name"`
}
type CategoryDTO struct {
	ParentID *int16 `json:"parent_id" validate:"omitempty,gte=0"`
	Name     string `json:"name"`
}
