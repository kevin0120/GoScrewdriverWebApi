package tightening_device

type TighteningDeviceConfig struct {
	// 控制器型号
	Model string `yaml:"model" json:"model"`

	// 控制器协议类型
	Protocol string `yaml:"protocol" json:"protocol"`

	// 连接地址(如果在控制器上配了连接地址，则下属所有工具共用此地址进行通信)
	Endpoint string `yaml:"endpoint" json:"endpoint"`

	// 控制器序列号
	SN string `yaml:"sn" json:"sn"`

	// 控制器名字
	ControllerName string `yaml:"name" json:"name"`

	// 工具列表
	Tools []ToolConfig `yaml:"tools" json:"children"`
}

type ToolConfig struct {
	// 工具序列号
	SN string `yaml:"sn" json:"sn"`

	// 工具通道号
	Channel int `yaml:"channel" json:"channel"`

	// 连接地址
	Endpoint string `yaml:"endpoint" json:"endpoint"`
}

type SocketSelectorConfig struct {
	Enable   bool   `yaml:"enable" json:"enable"`
	Endpoint string `yaml:"endpoint" json:"endpoint"`
}

type Config struct {
	Enable         bool                     `yaml:"enable" json:"enable"`
	SocketSelector SocketSelectorConfig     `yaml:"socket_selector" json:"socket_selector"`
	Devices        []TighteningDeviceConfig `yaml:"devices" json:"devices"`
}

func NewConfig() Config {

	return Config{
		Enable: true,
		Devices: []TighteningDeviceConfig{

			{
				Model:          ModelDesoutterCvi3,
				Protocol:       "OpenProtocol",
				Endpoint:       "tcp://192.168.5.232:4545",
				SN:             "ControllerSn",
				ControllerName: "ControllerName",
				Tools: []ToolConfig{
					{
						SN:       "ToolSn",
						Endpoint: "tcp://192.168.0.26:4545",
						Channel:  1,
					},
				},
			},
		},
	}
}

func (c Config) Validate() error {

	return nil
}
