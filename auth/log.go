package auth

import "github.com/sirupsen/logrus"

var (
	log = logrus.StandardLogger().WithField("module", "auth")
)
