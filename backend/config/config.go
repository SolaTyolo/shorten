package config

import (
	"log"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	Name string
}

func Init(cfg string) error {
	c := Config{
		Name: cfg,
	}

	if err := c.initConfig(); err != nil {
		return err
	}

	// watching config
	c.watchConfig()
	return nil
}

func (c *Config) initConfig() error {
	if c.Name != "" {
		// get config file from path
		viper.SetConfigFile(c.Name)
	} else {
		// default config
		viper.AddConfigPath("conf")
		viper.SetConfigName("config")
	}

	viper.SetConfigType("yaml")
	// get environment , prefix
	viper.AutomaticEnv()
	viper.SetEnvPrefix("GOSERVER")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	// parse config
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}

func (c *Config) watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("Config file chaned: %s", e.Name)
	})
}
