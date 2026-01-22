package profileService

import (
	"context"

	"time"

	"github.com/sewaustav/CaseGoProfile/internal/profile/dto"
	"github.com/sewaustav/CaseGoProfile/internal/profile/models"
)

func (s *ProfileService) CreateProfileService(
	ctx context.Context, 
	req dto.CreateProfileRequest, 
	userID int64) (
		*models.Profile, 
		[]models.UserSocialLink, 
		[]models.UserPurpose, 
		error,
		) {

	tx, err := s.repo.Begin(ctx)
    if err != nil {
        return nil, nil, nil, err
    }

	defer tx.Rollback()

	txRepo := s.repo.WithTx(tx)
	
	var sexPtr *models.UserSex
	if req.Info.Sex != nil {
		sexValue := models.UserSex(*req.Info.Sex)
		sexPtr = &sexValue
	}

	now := time.Now()
	
	profile := &models.Profile{
		UserID: userID,
		Avatar: req.Info.Avatar,
		IsActive: true,
		Description: req.Info.Description,
		Username: req.Info.Username,
		Name: req.Info.Name,
		Surname: req.Info.Surname,
		Email: req.Info.Email,
		CaseCount: 0,
		CreatedAt: now,
		UpdatedAt: now,
		// optional - nil possible
		Patronymic: req.Info.Patronymic,
		PhoneNumber: req.Info.PhoneNumber,
		Sex: sexPtr,
		Profession: req.Info.Profession,
	}

	var socialLinks []models.UserSocialLink

	for _, link := range req.SocialLinks {
		socialLinks = append(socialLinks, models.UserSocialLink{
			Type: link.Type, 
			URL: link.URL, 
			UserID: userID,
		})
	}

	var purposes []models.UserPurpose

	for _, purpose := range req.Purposes {
		purposes = append(purposes, models.UserPurpose{
			Purpose: purpose.Purpose,
			UserID: userID,
		})
	}

	createdProfile, err := txRepo.CreateProfile(ctx, profile)
	if err != nil {
		return nil, nil, nil, err 
	}
	addedLinks, err := txRepo.AddSocial(ctx, socialLinks)
	if err != nil {
		return nil, nil, nil, err
	}
	createdPurposes, err := txRepo.AddPurposes(ctx, purposes)
	if err != nil {
		return nil, nil, nil, err 
	}

	if err := tx.Commit(); err != nil {
        return nil, nil, nil, err
    }

	return createdProfile, addedLinks, createdPurposes, nil
	
}

// method for put
func(s *ProfileService) UpdateProfileService(ctx context.Context, userID int64, req dto.ProfileInfoDTO) (*models.Profile, error) {
	var sexPtr *models.UserSex
	if req.Sex != nil {
		sexValue := models.UserSex(*req.Sex)
		sexPtr = &sexValue
	}

	// todo - check is username unique

	now := time.Now()
	
	profile := &models.Profile{
		UserID: userID,
		Avatar: req.Avatar,
		IsActive: true,
		Description: req.Description,
		Username: req.Username,
		Name: req.Name,
		Surname: req.Surname,
		Email: req.Email,
		CaseCount: 0,
		UpdatedAt: now,
		// optional - nil possible
		Patronymic: req.Patronymic,
		PhoneNumber: req.PhoneNumber,
		Sex: sexPtr,
		Profession: req.Profession,
	}

	updatedProfile, err := s.repo.UpdateProfile(ctx, profile) 
	if err != nil {
		return nil, err 
	}

	return updatedProfile, nil
}


// partial - update. Note - in future special method for email/phone update
func (s *ProfileService) PatchProfile(
	ctx context.Context,
	userID int64,
	req dto.UpdateProfilePartialDTO,
) (*models.Profile, error) {

	// todo - check is username unique


	profile, err := s.repo.PathcProfile(ctx, userID, req)

	if err != nil {
		return nil, err 
	}

	return profile, nil 
}

func (s *ProfileService) UpdateSocialLink(ctx context.Context, req dto.SocialLinkDTO, userID int64, id int64) ([]models.UserSocialLink, error) {
	link := &models.UserSocialLink{
		ID: id,
		UserID: userID,
		Type: req.Type,
		URL: req.URL,
	}
	
	links, err := s.repo.EditSocial(ctx, link)

	if err != nil {
		return nil, err 
	}

	return links, nil

}

func (s *ProfileService) UpdatePurpose(ctx context.Context, req dto.UserPurposeDTO, userID, id int64) ([]models.UserPurpose, error) {
	purpose := &models.UserPurpose{
		ID: id,
		UserID: userID,
		Purpose: req.Purpose,
	}

	purposes, err := s.repo.EditPurpose(ctx, purpose)
	if err != nil {
		return nil, err 
	}

	return purposes, nil 
}