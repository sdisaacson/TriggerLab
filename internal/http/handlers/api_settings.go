package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/sdisaacson/TriggerLab/internal/api"
	"github.com/sdisaacson/TriggerLab/internal/config"
)

type apiSettings struct {
	cfg config.Config
}

// ApiSettings returns application settings.
func (s *apiSettings) ApiSettings(c echo.Context) error {
	return c.JSON(http.StatusOK, api.AppSettings{
		Limits: api.AppSettingsLimits{
			MaxRequests:        int(s.cfg.MaxRequests),
			MaxWebhookBodySize: int(s.cfg.MaxRequestBodySize),
			SessionLifetimeSec: int(s.cfg.SessionTTL.Seconds()),
		},
	})
}