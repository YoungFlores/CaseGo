package mocks

import (
	"context"

	"github.com/YoungFlores/Case_Go/Profile/internal/profession_categories/models"
	"github.com/stretchr/testify/mock"
)

type CategoryRepoMock struct {
	mock.Mock
}

func (m *CategoryRepoMock) CreateCategory(ctx context.Context, category *models.Category) (*models.Category, error) {
	args := m.Called(ctx, category)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Category), args.Error(1)
}

func (m *CategoryRepoMock) GetCategoryByID(ctx context.Context, id int16) (*models.Category, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Category), args.Error(1)
}

func (m *CategoryRepoMock) GetCategories(ctx context.Context) ([]models.Category, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Category), args.Error(1)
}

func (m *CategoryRepoMock) GetCategoriesByParent(ctx context.Context, parentID int16) ([]models.Category, error) {
	args := m.Called(ctx, parentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Category), args.Error(1)
}

func (m *CategoryRepoMock) GetParentOfCategory(ctx context.Context, id int16) (*int16, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*int16), args.Error(1)
}
