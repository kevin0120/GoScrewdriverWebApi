package tightening_device

type ITighteningProtocol interface {

	// Name 协议名称
	Name() string

	// NewController 创建控制器
	NewController(cfg *TighteningDeviceConfig) (ITighteningController, error)
}

type ITighteningController interface {
	Model() string
	Start() error
}

type ITighteningTool interface {
	Mode() string
	Start() error
}
