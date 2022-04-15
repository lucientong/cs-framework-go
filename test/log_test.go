package test

import (
	"cs/config"
	"cs/log"
	"testing"
)

func TestLog(t *testing.T) {
	config.AddConfigPath("../conf")
	config.Init()
	log.Init()

	log.Infof("start: %d", 123)
	log.Error("dddddd")
}
