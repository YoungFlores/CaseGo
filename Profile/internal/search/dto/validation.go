package dto

type SearchDTO struct {
	ProfessionID *int16  `json:"profession_id"`
	Profession   *string `json:"profession"`
	MinAge       *int    `json:"min_age"`
	MaxAge       *int    `json:"max_age"`
	City         *string `json:"city"`
	Sex          *int16  `json:"sex"`

	// add skills later
}

type SearchByFIODTO struct {
	Name       *string `json:"name"`
	Surname    *string `json:"surname"`
	Patronymic *string `json:"patronymic"`
}

type SearchHelpersDTO struct {
	Limit          *uint64
	Page           *uint64
	OrderBy        *string
	OrderDirection *string
}
