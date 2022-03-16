package internal

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

var (
	debug = pflag.BoolP("debug", "d", false, "Enable debug output")
)

func GetLog(module string) logrus.FieldLogger {
	if debug == nil {
		pflag.Parse()
	}

	if *debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	return logrus.StandardLogger().WithField("module", module)
}
