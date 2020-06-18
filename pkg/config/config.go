package config

import (
	"fmt"
	"github.com/cxnam/prometheus-pusher/pkg/logger"
	"github.com/spf13/viper"
	"strings"
)

type Schema struct {
	Systemmetric struct {
		Prometheusurl   string `mapstructure:"prom_url"`
		Statuspageurl   string `mapstructure:"page_url"`
		Statuspagetoken string `mapstructure:"page_token"`
		Statuspageid    string `mapstructure:"page_id"`
	} `mapstructure:"systemmetric"`
}

var (
	log    = logger.GetLogger("Config")
	Config Schema
)

func init() {
	config := viper.New()
	config.SetConfigName("config")
	config.SetConfigType("yaml")
	config.AddConfigPath(".")        						// Look for config in current directory.
	config.AddConfigPath("pkg/config/")        // Look for config needed for tests.
	config.AddConfigPath("config/")  // Optionally look for config in the working directory.
	config.AddConfigPath("../config/")
	config.AddConfigPath("../../")

	config.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	config.AutomaticEnv()

	err := config.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	err = config.Unmarshal(&Config)
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	log.Infof("Config: %+v", Config)
}