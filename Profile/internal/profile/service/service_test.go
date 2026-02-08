package profileService_test

import (
	"context"
	"testing"
	"time"

	dto "github.com/YoungFlores/Case_Go/Profile/internal/profile/dto"
	"github.com/YoungFlores/Case_Go/Profile/internal/profile/models"
	repo_mocks "github.com/YoungFlores/Case_Go/Profile/internal/profile/repository/profile_repo/mocks"
	service "github.com/YoungFlores/Case_Go/Profile/internal/profile/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateProfileService(t *testing.T) {
	mockRepo := new(repo_mocks.ProfileRepoMock)
	mockTx := new(repo_mocks.TxMock)
	svc := service.NewProfileService(mockRepo)

	ctx := context.Background()
	userID := int64(1)
	userInfo := models.UserIdentity{UserID: userID}

	sex := 0
	req := dto.CreateProfileRequest{
		Info: dto.ProfileInfoDTO{
			Name:     "John",
			Surname:  "Doe",
			Username: "johndoe",
			//Email:       "john@example.com",
			Sex:         &sex,
			Description: "Test user",
		},
		SocialLinks: []dto.SocialLinkDTO{
			{Type: "twitter", URL: "http://twitter.com/john"},
		},
		Purposes: []dto.UserPurposeDTO{
			{Purpose: "learning"},
		},
	}

	expectedProfile := &models.Profile{
		UserID:   userID,
		Username: "johndoe",
		//Email:    "john@example.com",
	}

	expectedLinks := []models.UserSocialLink{
		{Type: "twitter", URL: "http://twitter.com/john", UserID: userID},
	}

	expectedPurposes := []models.UserPurpose{
		{Purpose: "learning", UserID: userID},
	}

	// Mock expectations
	mockRepo.On("Begin", ctx).Return(mockTx, nil)
	mockRepo.On("WithTx", mockTx).Return(mockRepo) // Return same repo mock for simplicity

	// Since WithTx returns mockRepo, the following calls are on mockRepo
	mockRepo.On("CreateProfile", ctx, mock.AnythingOfType("*models.Profile")).Return(expectedProfile, nil)
	mockRepo.On("AddSocial", ctx, mock.AnythingOfType("[]models.UserSocialLink")).Return(expectedLinks, nil)
	mockRepo.On("AddPurposes", ctx, mock.AnythingOfType("[]models.UserPurpose")).Return(expectedPurposes, nil)

	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil) // Defer always calls Rollback

	// Execution
	result, err := svc.CreateProfileService(ctx, req, userInfo)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedProfile.Username, result.UsrProfile.Username)
	assert.Len(t, result.UsrSocials, 1)
	assert.Len(t, result.UsrPurposes, 1)

	mockRepo.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestGetUserProfileService(t *testing.T) {
	mockRepo := new(repo_mocks.ProfileRepoMock)
	svc := service.NewProfileService(mockRepo)

	ctx := context.Background()
	userID := int64(1)
	userInfo := models.UserIdentity{UserID: userID}

	expectedProfile := &models.Profile{
		UserID:   userID,
		Username: "johndoe",
		IsActive: true,
	}
	expectedLinks := []models.UserSocialLink{}
	expectedPurposes := []models.UserPurpose{}

	mockRepo.On("GetUserProfile", ctx, userID).Return(expectedProfile, nil)
	mockRepo.On("GetUserPurposes", ctx, userID).Return(expectedPurposes, nil)
	mockRepo.On("GetUserSocials", ctx, userID).Return(expectedLinks, nil)

	result, err := svc.GetUserProfileService(ctx, userInfo)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, userID, result.UsrProfile.UserID)

	mockRepo.AssertExpectations(t)
}

func TestUpdateProfileService(t *testing.T) {
	mockRepo := new(repo_mocks.ProfileRepoMock)
	svc := service.NewProfileService(mockRepo)

	ctx := context.Background()
	userID := int64(1)
	userInfo := models.UserIdentity{UserID: userID}

	req := dto.ProfileInfoDTO{
		Name:    "John",
		Surname: "Smith",
	}

	expectedProfile := &models.Profile{
		UserID:    userID,
		Name:      "John",
		Surname:   "Smith",
		IsActive:  true,
		UpdatedAt: time.Now(),
	}

	mockRepo.On("UpdateProfile", ctx, mock.AnythingOfType("*models.Profile")).Return(expectedProfile, nil)

	result, err := svc.UpdateProfileService(ctx, userInfo, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Smith", result.Surname)

	mockRepo.AssertExpectations(t)
}
