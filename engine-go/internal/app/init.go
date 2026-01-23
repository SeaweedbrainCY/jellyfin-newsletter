package app

import (
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/config"
	"go.uber.org/zap"
)

func InitApplicationContext(config *config.Configuration, logger *zap.Logger) *ApplicationContext {
	return &ApplicationContext{
		Logger: logger,
		Config: config,
	}
}
