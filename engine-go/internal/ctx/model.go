package ctx

import (
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/config"
	"go.uber.org/zap"
)

type Context struct {
	Config *config.Configuration
	Logger *zap.Logger
}
