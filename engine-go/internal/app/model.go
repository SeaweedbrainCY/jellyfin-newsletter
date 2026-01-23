package app

import (
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/config"
	"go.uber.org/zap"
)

type ApplicationContext struct {
	Config *config.Configuration
	Logger *zap.Logger
}
