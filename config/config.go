package config

import (
	"encoding/json"
	"fmt"
	"github.com/kevin0120/GoScrewdriverWebApi/services/diagnostic"
	"github.com/kevin0120/GoScrewdriverWebApi/services/http/httpd"
	"github.com/kevin0120/GoScrewdriverWebApi/services/opclient/openprotocol"
	"github.com/kevin0120/GoScrewdriverWebApi/services/opclient/tightening_device"
	"github.com/kevin0120/GoScrewdriverWebApi/utils"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Database struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UdpClient struct {
	RemoteHost string `json:"remote_host"`
	RemotePort int    `json:"remote_port"`
	LocalPort  int    `json:"local_port"`
}
type Config struct {
	Hostname string `yaml:"hostname" json:"hostname"`
	DataDir  string `yaml:"data_dir" json:"data_dir"`
	SN       string `yaml:"serial_no" json:"serial_no"`
	Path     string `yaml:"-" json:"-"`
	DocPath  string `yaml:"doc_path" json:"doc_path"`

	OpPort           int                      `json:"op_port"`
	UdpClient        *UdpClient               `json:"udp_client"`
	Database         *Database                `json:"database"`
	TighteningDevice tightening_device.Config `json:"tightening_device"`
	OpenProtocol     openprotocol.Config      `json:"openprotocol"`
	Logging          diagnostic.Config        `json:"logging"`
	HTTP             httpd.Config             `json:"httpd"`
}

var MyConfig *Config

func init() {
	// 生成默认配置
	MyConfig = getDefaultConfig()
	// 读取配置文件并覆盖默认配置
	exePath, err := os.Getwd()
	if err != nil {
		fmt.Println("无法获取可执行文件的路径:", err)
	}
	err = readConfigFile(exePath+"/config/config.json", MyConfig)
	if err != nil {
		fmt.Println("Failed to read config file:", err)
	}
}

func readConfigFile(filename string, config *Config) error {
	file, err := os.Open(filename)
	if err != nil {
		// 配置文件不存在，直接返回
		if os.IsNotExist(err) {
			err := os.MkdirAll(filepath.Dir(filename), 0755)
			if err != nil {
				fmt.Println("无法创建目录:", err)
				return err
			}
			out, err := json.MarshalIndent(config, "", "  ")
			err = ioutil.WriteFile(filename, out, fs.ModePerm)
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, config)
	if err != nil {
		return err
	}
	return nil
}

func getDefaultConfig() *Config {
	// 生成默认配置的逻辑
	return &Config{
		Hostname: "localhost",
		SN:       utils.GenerateID(),
		DocPath:  filepath.Join("./", "doc"),
		DataDir:  filepath.Join("./", ".rush"),

		OpPort: 4545,
		Database: &Database{Host: "192.168.10.122",
			Port: 8082, Username: "ROOT",
			Password: "!23!QQA"},
		UdpClient: &UdpClient{
			RemoteHost: "211.254.254.250",
			RemotePort: 8080,
			LocalPort:  50004},
		TighteningDevice: tightening_device.NewConfig(),
		OpenProtocol:     openprotocol.NewConfig(),
		Logging:          diagnostic.NewConfig(),
		HTTP:             httpd.NewConfig(),
	}
}

func GetConfig() *Config {
	return MyConfig
}
