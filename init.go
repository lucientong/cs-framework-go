package cs

import (
	"cs/config"
	"cs/log"
)

func Init() {
	config.Init()
	log.Init()
}

func Run() {
	log.Info("====== App start ======")
	log.Infof("%+v", config.AllSettings())
}
