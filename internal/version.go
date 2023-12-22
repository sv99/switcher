package internal

import (
	"github.com/rs/zerolog/log"
)

// filled from ldflags
// version defined in the Makefile
var (
	Version   = ""
	GoVersion = ""
	BuildTime = ""
)

func PrintVersion() {
	log.Debug().Msgf("Version: %s", Version)
	log.Debug().Msgf("%s", GoVersion)
	log.Debug().Msgf("Build Time: %s", BuildTime)
}
