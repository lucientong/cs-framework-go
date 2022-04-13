package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

var vipers *Config

var (
	configType = "yaml"
	configName = "config"
	configPath = []string{"./conf"}
)

// Config 配置
type Config struct {
	*viper.Viper
}

// SetConfigType 设置配置文件类型
// support JSON, TOML, YAML, HCL, INI, envfile or Java properties formats
func SetConfigType(ct string) {
	configType = ct
}

// SetConfigName 设置配置文件名称
func SetConfigName(cn string) {
	configName = cn
}

// AddConfigPath 设置配置文件路径
func AddConfigPath(cp string) {
	configPath = append(configPath, cp)
}

// Use 获取配置
func (c *Config) Use(key string) *Config {
	if c.Sub(key) == nil {
		return nil
	}
	return &Config{
		c.Sub(key),
	}
}

// GetInt 获取int值
func (c *Config) GetInt(key string, defaultValue int) int {
	val := c.Get(key)
	if val == nil {
		return defaultValue
	}
	value, err := cast.ToIntE(val)
	if err != nil {
		return defaultValue
	}
	return value
}

// GetInt8 获取int8值
func (c *Config) GetInt8(key string, defaultValue int8) int8 {
	val := c.Get(key)
	if val == nil {
		return defaultValue
	}
	value, err := cast.ToInt8E(val)
	if err != nil {
		return defaultValue
	}
	return value
}

// GetInt16 获取int16值
func (c *Config) GetInt16(key string, defaultValue int16) int16 {
	val := c.Get(key)
	if val == nil {
		return defaultValue
	}
	value, err := cast.ToInt16E(val)
	if err != nil {
		return defaultValue
	}
	return value
}

// GetInt64 获取int64值
func (c *Config) GetInt64(key string, defaultValue int64) int64 {
	val := c.Get(key)
	if val == nil {
		return defaultValue
	}
	value, err := cast.ToInt64E(val)
	if err != nil {
		return defaultValue
	}
	return value
}

// GetUint 获取uint值
func (c *Config) GetUint(key string, defaultValue uint) uint {
	val := c.Get(key)
	if val == nil {
		return defaultValue
	}
	value, err := cast.ToUintE(val)
	if err != nil {
		return defaultValue
	}
	return value
}

// GetUint8 获取uint8值
func (c *Config) GetUint8(key string, defaultValue uint8) uint8 {
	val := c.Get(key)
	if val == nil {
		return defaultValue
	}
	value, err := cast.ToUint8E(val)
	if err != nil {
		return defaultValue
	}
	return value
}

// GetUint16 获取uint16值
func (c *Config) GetUint16(key string, defaultValue uint16) uint16 {
	val := c.Get(key)
	if val == nil {
		return defaultValue
	}
	value, err := cast.ToUint16E(val)
	if err != nil {
		return defaultValue
	}
	return value
}

// GetUint32 获取uint32值
func (c *Config) GetUint32(key string, defaultValue uint32) uint32 {
	val := c.Get(key)
	if val == nil {
		return defaultValue
	}
	value, err := cast.ToUint32E(val)
	if err != nil {
		return defaultValue
	}
	return value
}

// GetUint64 获取uint64值
func (c *Config) GetUint64(key string, defaultValue uint64) uint64 {
	val := c.Get(key)
	if val == nil {
		return defaultValue
	}
	value, err := cast.ToUint64E(val)
	if err != nil {
		return defaultValue
	}
	return value
}

// GetFloat32 获取float32值
func (c *Config) GetFloat32(key string, defaultValue float32) float32 {
	val := c.Get(key)
	if val == nil {
		return defaultValue
	}
	value, err := cast.ToFloat32E(val)
	if err != nil {
		return defaultValue
	}
	return value
}

// GetFloat64 获取float64值
func (c *Config) GetFloat64(key string, defaultValue float64) float64 {
	val := c.Get(key)
	if val == nil {
		return defaultValue
	}
	value, err := cast.ToFloat64E(val)
	if err != nil {
		return defaultValue
	}
	return value
}

// GetString 获取string值
func (c *Config) GetString(key string, defaultValue string) string {
	val := c.Get(key)
	if val == nil {
		return defaultValue
	}
	value, err := cast.ToStringE(val)
	if err != nil {
		return defaultValue
	}
	return value
}

// GetBool 获取bool值
func (c *Config) GetBool(key string, defaultValue bool) bool {
	val := c.Get(key)
	if val == nil {
		return defaultValue
	}
	value, err := cast.ToBoolE(val)
	if err != nil {
		return defaultValue
	}
	return value
}

// Init 初始化
func Init() {
	vipers = &Config{
		viper.New(),
	}
	vipers.SetConfigType(configType)
	vipers.SetConfigName(configName)
	for _, cp := range configPath {
		vipers.AddConfigPath(cp)
	}
	if err := vipers.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}
	vipers.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
	vipers.WatchConfig()
}
