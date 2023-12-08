package diagnostic

import (
	"github.com/kevin0120/GoScrewdriverWebApi/utils/toml"
	"time"
)

type Config struct {
	File   string        `yaml:"file"`
	Level  string        `yaml:"level"`
	MaxAge toml.Duration `yaml:"max_age"`
	Rotate toml.Duration `yaml:"rotate"`
}

func NewConfig() Config {
	return Config{
		File:   "STDERR",
		Level:  "DEBUG",
		MaxAge: toml.Duration(time.Duration(31 * 24 * time.Hour)),
		Rotate: toml.Duration(time.Duration(24 * time.Hour)),
	}
}
