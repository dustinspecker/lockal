package dependency

import (
	"github.com/dustinspecker/lockal/internal/config"
)

type Dependency interface {
	Download(config.Config) error
}
