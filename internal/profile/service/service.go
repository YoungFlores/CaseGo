package profileService

import profilerepo "github.com/sewaustav/CaseGoProfile/internal/profile/repository/profile_repo"


type ProfileService struct {
	repo profilerepo.ProfileRepo
}

func NewProfileService(repo profilerepo.ProfileRepo) *ProfileService {
	return &ProfileService{
		repo: repo,
	}
}