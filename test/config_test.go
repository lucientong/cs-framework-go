package test

import "testing"
import "cs/config"

func TestConfig(t *testing.T) {
	config.AddConfigPath("../conf")
	config.SetConfigType("yaml")
	config.SetConfigName("config")

	config.Init()

	name := config.Use("info").MustString("name", "lucientong")
	schoolName := config.Use("info").Use("school").MustString("schoolName", "xiyou")
	departmentName := config.Use("info").Use("school").MustString("departmentName", "computer")
	age := config.MustUse("info").MustInt("age", 20)
	isMarried := config.Use("info").MustBool("isMarried", true)

	t.Logf("\nname: %s\nschoolName:%s\ndepartmentName:%s\nage:%d\nisMarried:%v", name, schoolName, departmentName, age, isMarried)
}
