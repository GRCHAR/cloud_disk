package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// Config 为系统全局配置
type Config struct {
	Server struct {
		Port         int    `yaml:"port"`
		TempDirPath  string `yaml:"tempDirPath"`
		MergeDirPath string `yaml:"mergeDirPath"`
	}
}

var serverconfig Config

func init() {
	serverconfig = getConfig("./application.yaml")
}

func getConfig(filePath string) Config {
	config := Config{}
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("解析config.yaml读取错误: %v", err)
	}

	// fmt.Println(string(content))
	// fmt.Printf("init data: %v", config)
	if yaml.Unmarshal(content, &config) != nil {
		log.Fatalf("解析config.yaml出错: %v", err)
	}
	// fmt.Printf("File config: %v", config)
	return config
}

func GetConfig() *Config {
	return &serverconfig
}
