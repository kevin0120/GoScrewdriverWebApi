package hmi

type Config struct {
	Enable bool `yaml:"enable"`
}

func NewConfig() Config {
	return Config{
		Enable: true,
	}
}

func (c Config) Validate() error {
	return nil
}

type ControlData struct {
	D003 string `json:"D003"`
	D010 string `json:"D010"`
	D012 int    `json:"D012"`
	D018 int    `json:"D018"`
	D030 string `json:"D030"`
	D038 int    `json:"D038"`
	D042 string `json:"D042"`
	D043 string `json:"D043"`
	D050 string `json:"D050"`
	D150 string `json:"D150"`
}

type TighteningDeviceCmd struct {
	Sn          string      `json:"sn"`
	Mid         string      `json:"mid"`
	ControlData ControlData `json:"data"`
}
