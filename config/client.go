package config

import (
	"github.com/spf13/viper"
)

type LocalConfig struct {
	configPath string
	configMap  map[string]interface{}
}

func (c *LocalConfig) SetConfigPath(path string) *LocalConfig {
	c.configPath = path
	return c
}

func (c *LocalConfig) SetConfigMap(configMap map[string]interface{}) *LocalConfig {
	c.configMap = configMap
	return c
}

func (c *LocalConfig) Read() error {
	for name, conf := range c.configMap {
		v := viper.New()
		v.SetConfigName(name)
		v.AddConfigPath(c.configPath)

		err := v.ReadInConfig()
		if err != nil {
			return err
		}

		err = v.Unmarshal(conf)
		if err != nil {
			return err
		}
	}

	return nil
}
