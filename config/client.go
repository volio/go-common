package config

import (
	"github.com/spf13/viper"
)

type LocalConfigClient interface {
	SetConfigPath(path string) LocalConfigClient
	SetConfigMap(configMap map[string]interface{}) LocalConfigClient
	Read() error
}

type localConfigClient struct {
	configPath string
	configMap  map[string]interface{}
}

func (c *localConfigClient) SetConfigPath(path string) LocalConfigClient {
	c.configPath = path
	return c
}

func (c *localConfigClient) SetConfigMap(configMap map[string]interface{}) LocalConfigClient {
	c.configMap = configMap
	return c
}

func (c *localConfigClient) Read() error {
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
