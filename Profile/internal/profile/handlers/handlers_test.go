package profile_handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	repoMocksCat "github.com/YoungFlores/Case_Go/Profile/internal/profession_categories/repo/mocks"
	dto "github.com/YoungFlores/Case_Go/Profile/internal/profile/dto"
	handlers "github.com/YoungFlores/Case_Go/Profile/internal/profile/handlers"
	"github.com/YoungFlores/Case_Go/Profile/internal/profile/models"
	repoMocks "github.com/YoungFlores/Case_Go/Profile/internal/profile/repository/profile_repo/mocks"
	profileService "github.com/YoungFlores/Case_Go/Profile/internal/profile/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateProfileHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(repoMocks.ProfileRepoMock)
	catRepo := new(repoMocksCat.CategoryRepoMock)
	mockTx := new(repoMocks.TxMock)
	svc := profileService.NewProfileService(mockRepo, catRepo)
	handler := handlers.NewProfileHandler(svc)

	userID := int64(123)

	sex := 0
	reqBody := dto.CreateProfileRequest{
		Info: dto.ProfileInfoDTO{
			Name:     "John",
			Surname:  "Doe",
			Username: "johndoe",
			Sex:      &sex,
		},
	}
	bodyBytes, _ := json.Marshal(reqBody)

	expectedProfile := &models.Profile{
		UserID:    userID,
		Username:  "johndoe",
		Name:      "John",
		Surname:   "Doe",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Mock expectations for Service -> Repo interaction
	mockRepo.On("Begin", mock.Anything).Return(mockTx, nil)
	mockRepo.On("WithTx", mockTx).Return(mockRepo)

	mockRepo.On("CreateProfile", mock.Anything, mock.AnythingOfType("*models.Profile")).Return(expectedProfile, nil)
	mockRepo.On("AddSocial", mock.Anything, mock.Anything).Return([]models.UserSocialLink{}, nil)
	mockRepo.On("AddPurposes", mock.Anything, mock.Anything).Return([]models.UserPurpose{}, nil)

	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil)

	// Setup Router
	r := gin.Default()
	r.POST("/profile", handler.CreateProfileHandler)

	// Request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/profile", bytes.NewBuffer(bodyBytes))
	r.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.UserProfile
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "johndoe", response.UsrProfile.Username)

	mockRepo.AssertExpectations(t)
}

func TestGetUserProfileHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(repoMocks.ProfileRepoMock)
	catRepo := new(repoMocksCat.CategoryRepoMock)
	svc := profileService.NewProfileService(mockRepo, catRepo)
	handler := handlers.NewProfileHandler(svc)

	userID := int64(123) // Matches utils.go hardcoded value

	expectedProfile := &models.Profile{
		UserID:   userID,
		Username: "johndoe",
		IsActive: true,
	}

	mockRepo.On("GetUserProfile", mock.Anything, userID).Return(expectedProfile, nil)
	mockRepo.On("GetUserPurposes", mock.Anything, userID).Return([]models.UserPurpose{}, nil)
	mockRepo.On("GetUserSocials", mock.Anything, userID).Return([]models.UserSocialLink{}, nil)

	r := gin.Default()
	r.GET("/profile", handler.GetUserProfileHandler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/profile", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.UserProfile
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, userID, response.UsrProfile.UserID)

	mockRepo.AssertExpectations(t)
}
