package openprotocol

import (
	"github.com/kevin0120/GoScrewdriverWebApi/utils/toml"
	"github.com/kevin0120/GoScrewdriverWebApi/utils/typeDef"
	"time"
)

const (
	OpenProtocolDefaultGetTollInfoPeriod = toml.Duration(time.Hour * 12)
)

type Config struct {
	SkipJobs          []int         `yaml:"skip_job"`
	DataIndex         int           `yaml:"data_index"`
	VinIndex          []int         `yaml:"vin_index"`
	GetToolInfoPeriod toml.Duration `yaml:"tool_info_period"`
	DefaultMode       string        `yaml:"default_mode"`
}

func NewConfig() Config {

	return Config{
		SkipJobs:          []int{250},
		DataIndex:         1,
		VinIndex:          []int{0, 1},
		GetToolInfoPeriod: OpenProtocolDefaultGetTollInfoPeriod,
		DefaultMode:       typeDef.MODE_JOB,
	}
}

func (c Config) Validate() error {
	return nil
}
