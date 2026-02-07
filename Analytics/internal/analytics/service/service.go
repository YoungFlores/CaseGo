package analyticsService

import analyticsRepo "github.com/YoungFlores/Case_Go/Analytics/internal/analytics/repo"

type AnalyticsService struct {
	repo analyticsRepo.AnalyticsRepo
}

func NewAnalyticsService(repo analyticsRepo.AnalyticsRepo) *AnalyticsService {
	return &AnalyticsService{
		repo: repo,
	}
}

