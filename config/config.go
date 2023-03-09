package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

type config struct {
	Debug   bool   `yaml:"debug"`
	LogPath string `yaml:"logPath"`

	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"server"`

	Forwards []string `yaml:"forwards"`

	Settings struct {
		VerifyPeriod int `yaml:"verifyPeriod"`
		LimitPeriod  int `yaml:"limitPeriod"`
		LimitTimes   int `yaml:"limitTimes"`
	} `yaml:"settings"`

	Machine struct {
		Name   string `yaml:"name"`
		Secret string `yaml:"secret"`
	} `yaml:"machine"`
}

var Config = new(config)

func init() {
	//viper.SetConfigFile(filepath.Dir(os.Args[0]) + "/config.yml")
	viper.SetConfigFile("./config.yml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(fmt.Errorf("Config read failed: %v \n", err))
	}
	if err := viper.Unmarshal(Config); err != nil {
		log.Fatal(fmt.Errorf("Config unmarshal failed: %v \n", err))
	}
	viper.WatchConfig()
}
