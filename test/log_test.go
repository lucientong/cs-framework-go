package test

import (
	"cs/config"
	"cs/logger"
	"testing"
)

func TestLog(t *testing.T) {
	config.AddConfigPath("../conf")
	config.Init()
	logger.Init()

	logger.Infof("start: %d", 123)
	logger.Error("dddddd")
}
