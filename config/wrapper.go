// Package config provides to users directly
package config

// Use 获取配置对象
func Use(key string) *Config {
	if vipers.Sub(key) == nil {
		return nil
	}
	return &Config{
		vipers.Sub(key),
	}
}

// MustUse 获取不到对应配置抛panic
func MustUse(key string) *Config {
	if vipers.Sub(key) == nil {
		panic("can not resolve key: " + key)
	}
	return Use(key)
}

// AllSettings 获取所有配置
func AllSettings() map[string]interface{} {
	return vipers.AllSettings()
}

// Set 设置配置
func Set(key string, value interface{}) {
	vipers.Set(key, value)
}
