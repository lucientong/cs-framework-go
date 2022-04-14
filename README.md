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
2. 根据不同的表配置不同的结构体，增删改查与gorm相同


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

## 日志
基于[zap](https://github.com/uber-go/zap)实现，并基于[lumberjack](https://github.com/natefinch/lumberjack)实现日志切割，都需要配置在log条目下。

- 如何使用

1. 配置
1.1 日志在配置文件中可以配置 showConsole（是否在控制台显示），path（日志文件目录），level（日志文件级别，默认为INFO），appLogName（运行日志文件名，默认为app），errLogName（错误日志文件名，默认为err）。
1.2 日志切割可配置 maxSize（日志轮换前文件最大大小），maxAge（旧日志最长保留时间），maxBackups（保留旧日志文件最大数量，会因为maxAge配置而清除），compress（是否使用gzip压缩日志文件）
2. 使用

框架提供Debug(), Info()，Error()等zap原生格式日志，以及Debugf(), Infof()，Errorf()等printf格式日志，通过logger包直接调用。

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

db:
  default:
    driver: "mysql"
    host: "127.0.0.1"
    port: 3306
    database: "lucien_test_db"
    username: "root"
    password: "password"

log:
  showConsole: true
  path: "./log"
  level: "info" # 支持debug,info,warn,error。error会额外打到error.log
  maxSize: 64 # 日志文件最大的大小，单位MB
  maxBackups: 16 # 保留最大旧日志文件数，但过期（>MaxAge）仍可能删除
  maxAge: 30 # 保留日志的天数
  compress: false # 是否压缩备份日志

```

log_test.go

```go
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
```

控制台和文件输出如下：

```bash
[2022-04-14 16:55:31.923]	[INFO]	[test/log_test.go:14 test.TestLog]	start: [123]
[2022-04-14 16:55:31.925]	[ERROR]	[test/log_test.go:15 test.TestLog]	dddddd
```