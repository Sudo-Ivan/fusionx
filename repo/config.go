package repo

import (
	"strconv"
	"time"

	"github.com/0x2e/fusion/model"
	"gorm.io/gorm"
)

func NewConfig(db *gorm.DB) *Config {
	return &Config{
		db: db,
	}
}

type Config struct {
	db *gorm.DB
}

func (c *Config) Get(key string) (string, error) {
	var config model.Config
	err := c.db.Where("key = ?", key).First(&config).Error
	if err != nil {
		return "", err
	}
	return config.Value, nil
}

func (c *Config) Set(key, value string) error {
	config := model.Config{
		Key:   key,
		Value: value,
	}
	return c.db.Save(&config).Error
}

func (c *Config) GetInt(key string, defaultValue int) (int, error) {
	value, err := c.Get(key)
	if err != nil {
		if err == ErrNotFound {
			return defaultValue, nil
		}
		return 0, err
	}
	
	result, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue, nil
	}
	return result, nil
}

func (c *Config) SetInt(key string, value int) error {
	return c.Set(key, strconv.Itoa(value))
}

func (c *Config) GetDuration(key string, defaultValue time.Duration) (time.Duration, error) {
	value, err := c.Get(key)
	if err != nil {
		if err == ErrNotFound {
			return defaultValue, nil
		}
		return 0, err
	}
	
	result, err := time.ParseDuration(value)
	if err != nil {
		return defaultValue, nil
	}
	return result, nil
}

func (c *Config) SetDuration(key string, value time.Duration) error {
	return c.Set(key, value.String())
}
