package test

import "testing"
import "cs/config"

func TestConfig(t *testing.T) {
	config.AddConfigPath("../conf")
	config.SetConfigType("yaml")
	config.SetConfigName("config")

	config.Init()

	name := config.Use("info").GetString("name", "lucientong")
	schoolName := config.Use("info").Use("school").GetString("schoolName", "xiyou")
	departmentName := config.Use("info").Use("school").GetString("departmentName", "computer")
	age := config.MustUse("info").GetInt("age", 20)
	isMarried := config.Use("info").GetBool("isMarried", true)

	t.Logf("\nname: %s\nschoolName:%s\ndepartmentName:%s\nage:%d\nisMarried:%v", name, schoolName, departmentName, age, isMarried)
}
