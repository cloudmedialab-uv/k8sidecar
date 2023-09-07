package config

import (
	"log"
	"os"
	"sync"
)

type Config struct {
	data map[string]string
}

var instance *Config

var once sync.Once

func NewIntance() *Config {
	once.Do(func() {
		instance = &Config{
			data: make(map[string]string),
		}
	})
	return instance
}

func (c *Config) Load(key string, defaultValue string) {
	value, exist := os.LookupEnv(key)

	if !exist || value == "" {
		if defaultValue == "" {
			log.Fatalf("La variable de entorno %s no esta definida", key)
		}
		c.data[key] = defaultValue
	} else {
		c.data[key] = value
	}
}

func (c *Config) Get(key string) string {
	return c.data[key]
}
