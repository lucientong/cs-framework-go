# C/S模式中间件

## 配置管理
基于[Viper](https://github.com/spf13/viper)实现，默认读取conf目录下的config.yaml文件。

- 如何设置配置文件？

```go
config.AddConfigPath("../conf")
config.SetConfigType("yaml")
config.SetConfigName("config")
```
使用 AddConfigPath 添加配置文件路径，框架会遍历查找。
使用 SetConfigType 设置配置文件类型，支持JSON, TOML, YAML, HCL, INI, envfile or Java properties
使用 SetConfigName 设置配置文件名称。

- example

config.yaml
```yaml
info:
#  name: "lucien"
  school:
    schoolName: "xupt"
    departmentName: "CS"
  age: 21
  isMarried: false

```
config_test.go
``` go
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

```

## ORM
基于[gorm](https://github.com/go-gorm/gorm)实现，默认读取conf目录下的config.yaml文件中db项配置。适配的数据库类型有：MySQL, PostgreSQL, SQlite, SQL Server，更改db下的driver既可更改驱动。
```yaml
db:
  default:
    driver: "mysql"  # MySQL
    # driver: "postgres" # PostgreSQL
    # driver: "sqlite"  # SQlite
    # driver: "sqlserver"  # SQL Server
```

- 如何使用？

1. 首先需要在配置文件中配置数据库信息（db项）

config.yaml
```yaml
info:
#  name: "lucien"
  school:
    schoolName: "xupt"
    departmentName: "CS"
  age: 21
  isMarried: false

db:
  default:
    driver: "mysql"
    host: "127.0.0.1"
    port: 3306
    database: "lucien_test_db"
    username: "root"
    password: "password"
```
1. 根据不同的表配置不同的结构体，增删改查与gorm相同


orm_test.go
```go
package test

import (
	"cs/config"
	"cs/db"
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

	data := DedaoUserData{}
	db.Default().First(&data)
	t.Log(data)
}

```


