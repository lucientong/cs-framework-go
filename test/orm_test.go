package test

import (
	"cs/config"
	"cs/db"
	"cs/log"
	"testing"
)

type DedaoUserData struct {
	Id        int    `gorm:"id" json:"id"`
	Name      string `gorm:"name" json:"name"`
	Frequency int    `gorm:"frequency" json:"frequency"`
}

func (d *DedaoUserData) TableName() string {
	return "dedao_user_data"
}

func TestGorm(t *testing.T) {
	config.AddConfigPath("../conf")
	config.SetConfigType("yaml")
	config.SetConfigName("config")

	config.Init()
	log.Init()

	data := DedaoUserData{}
	db.Default().First(&data)
	t.Log(data)
}
