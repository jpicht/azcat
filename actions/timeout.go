package actions

import (
	"time"

	"github.com/spf13/pflag"
)

var (
	timeout  = pflag.String("timeout", "default", "Timeout")
	defaults = map[Mode]time.Duration{
		EMode.Read():   3600 * time.Second,
		EMode.Write():  3600 * time.Second,
		EMode.List():   60 * time.Second,
		EMode.Remove(): 60 * time.Second,
		EMode.Ping():   1 * time.Second,
	}
)

func getTimeout(mode Mode) time.Duration {
	if *timeout == "default" {
		if default_, ok := defaults[mode]; ok {
			return default_
		}
		log.WithField("mode", mode.String()).Fatal("No default timeout for mode")
	}

	duration, err := time.ParseDuration(*timeout)
	if err != nil {
		log.WithError(err).Fatal("Invalid timeout specified")
	}

	return duration
}
