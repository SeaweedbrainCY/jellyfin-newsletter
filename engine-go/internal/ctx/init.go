package ctx

import (
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/config"
	"go.uber.org/zap"
)

func InitContext(config *config.Configuration, logger *zap.Logger) *Context {
	return &Context{
		Logger: logger,
		Config: config,
	}
}
