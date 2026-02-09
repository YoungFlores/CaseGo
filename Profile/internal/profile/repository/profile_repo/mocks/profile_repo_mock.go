package mocks

import (
	"context"
	"database/sql"

	dto "github.com/YoungFlores/Case_Go/Profile/internal/profile/dto"
	"github.com/YoungFlores/Case_Go/Profile/internal/profile/models"
	profilerepo "github.com/YoungFlores/Case_Go/Profile/internal/profile/repository/profile_repo"
	"github.com/stretchr/testify/mock"
)

type ProfileRepoMock struct {
	mock.Mock
}

func (m *ProfileRepoMock) Begin(ctx context.Context) (profilerepo.Tx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(profilerepo.Tx), args.Error(1)
}

func (m *ProfileRepoMock) WithTx(tx profilerepo.Tx) profilerepo.ProfileRepo {
	args := m.Called(tx)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(profilerepo.ProfileRepo)
}

func (m *ProfileRepoMock) CreateProfile(ctx context.Context, user *models.Profile) (*models.Profile, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Profile), args.Error(1)
}

func (m *ProfileRepoMock) AddSocial(ctx context.Context, links []models.UserSocialLink) ([]models.UserSocialLink, error) {
	args := m.Called(ctx, links)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.UserSocialLink), args.Error(1)
}

func (m *ProfileRepoMock) AddPurposes(ctx context.Context, purposes []models.UserPurpose) ([]models.UserPurpose, error) {
	args := m.Called(ctx, purposes)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.UserPurpose), args.Error(1)
}

func (m *ProfileRepoMock) GetProfileByID(ctx context.Context, id int64) (*models.Profile, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Profile), args.Error(1)
}

func (m *ProfileRepoMock) GetUserByProfileID(ctx context.Context, id, userID int64) (int64, error) {
	args := m.Called(ctx, id, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *ProfileRepoMock) GetUserProfile(ctx context.Context, userID int64) (*models.Profile, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Profile), args.Error(1)
}

func (m *ProfileRepoMock) GetUserSocials(ctx context.Context, userID int64) ([]models.UserSocialLink, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.UserSocialLink), args.Error(1)
}

func (m *ProfileRepoMock) GetUserPurposes(ctx context.Context, userID int64) ([]models.UserPurpose, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.UserPurpose), args.Error(1)
}

func (m *ProfileRepoMock) GetAllUsers(ctx context.Context, limit int) ([]models.Profile, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Profile), args.Error(1)
}

func (m *ProfileRepoMock) UpdateProfile(ctx context.Context, user *models.Profile) (*models.Profile, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Profile), args.Error(1)
}

func (m *ProfileRepoMock) PatchProfile(ctx context.Context, userID int64, updates dto.UpdateProfilePartialDTO) (*models.Profile, error) {
	args := m.Called(ctx, userID, updates)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Profile), args.Error(1)
}

func (m *ProfileRepoMock) UpdateLinks(ctx context.Context, links []models.UserSocialLink) ([]models.UserSocialLink, error) {
	args := m.Called(ctx, links)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.UserSocialLink), args.Error(1)
}

func (m *ProfileRepoMock) EditSocial(ctx context.Context, link *models.UserSocialLink) ([]models.UserSocialLink, error) {
	args := m.Called(ctx, link)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.UserSocialLink), args.Error(1)
}

func (m *ProfileRepoMock) UpdatePurposes(ctx context.Context, purposes []models.UserPurpose) ([]models.UserPurpose, error) {
	args := m.Called(ctx, purposes)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.UserPurpose), args.Error(1)
}

func (m *ProfileRepoMock) EditPurpose(ctx context.Context, purpose *models.UserPurpose) ([]models.UserPurpose, error) {
	args := m.Called(ctx, purpose)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.UserPurpose), args.Error(1)
}

func (m *ProfileRepoMock) DeletePurpose(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *ProfileRepoMock) DeleteSocial(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *ProfileRepoMock) DeleteProfile(ctx context.Context, userID int64) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *ProfileRepoMock) DeleteProfileWithoutRecovery(ctx context.Context, userID int64) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// TxMock implements profilerepo.Tx
type TxMock struct {
	mock.Mock
}

func (m *TxMock) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	a := m.Called(ctx, query, args)
	return a.Get(0).(sql.Result), a.Error(1)
}

func (m *TxMock) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	args := m.Called(ctx, query)
	return args.Get(0).(*sql.Stmt), args.Error(1)
}

func (m *TxMock) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	a := m.Called(ctx, query, args)
	return a.Get(0).(*sql.Rows), a.Error(1)
}

func (m *TxMock) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	a := m.Called(ctx, query, args)
	return a.Get(0).(*sql.Row) // this might be tricky to mock properly as *sql.Row is a struct
}

func (m *TxMock) Commit() error {
	args := m.Called()
	return args.Error(0)
}

func (m *TxMock) Rollback() error {
	args := m.Called()
	return args.Error(0)
}
