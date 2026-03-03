package app

import (
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/config"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"go.uber.org/zap"
)

type ApplicationContext struct {
	Config    *config.Configuration
	Logger    *zap.Logger
	Localizer *i18n.Localizer
}
