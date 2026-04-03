package app

import (
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/clock"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/config"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/i18n"
	"go.uber.org/zap"
)

func InitApplicationContext(
	config *config.Configuration,
	logger *zap.Logger,
	localizer *i18n.Localizer,
	clock clock.Interface,
) *ApplicationContext {
	return &ApplicationContext{
		Logger:    logger,
		Config:    config,
		Localizer: localizer,
		Clock:     clock,
	}
}
