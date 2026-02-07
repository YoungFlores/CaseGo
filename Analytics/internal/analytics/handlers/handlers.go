package analyticsHandlers

import analyticsService "github.com/YoungFlores/Case_Go/Analytics/internal/analytics/service"

type AnalyticsHandler struct {
	service *analyticsService.AnalyticsService
}

func NewAnalyticsHandler(service *analyticsService.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		service: service,
	}
}
