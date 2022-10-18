package config

import (
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

//go:generate moq -rm -out config_mock.go . Config
type Config interface {
	AuthToken() string
	SetAuthToken(string) error

	Leagues() []string
	SetLeagues([]string) error
}

func NewConfig() (Config, error) {
	home, err := homedir.Dir()
	if err != nil {
		return nil, err
	}

	// Search config in home directory with name ".cobra" (without extension).
	viper.AddConfigPath(home)
	viper.SetConfigName(".football-data")
	viper.SetConfigType("json")

	if err := viper.ReadInConfig(); err != nil {
		err := viper.WriteConfig()
		if err != nil {
			return nil, err
		}

		if err := viper.ReadInConfig(); err != nil {
			return nil, err
		}
	}

	return &viperCfg{}, nil
}

func NewTestConfig() *ConfigMock {
	cfg := testCfg{}
	mock := &ConfigMock{}
	mock.AuthTokenFunc = func() string {
		return cfg.AuthToken()
	}
	mock.SetAuthTokenFunc = func(token string) error {
		return cfg.SetAuthToken(token)
	}

	return mock
}

type viperCfg struct {
}

func (c viperCfg) AuthToken() string {
	return viper.GetString("footballdata.apikey")
}

func (c viperCfg) SetAuthToken(token string) error {
	viper.Set("footballdata.apikey", token)
	return viper.WriteConfig()
}

func (c viperCfg) Leagues() []string {
	return viper.GetStringSlice("footballdata.competitions")
}

func (c viperCfg) SetLeagues(leagues []string) error {
	viper.Set("footballdata.competitions", leagues)
	return viper.WriteConfig()
}

type testCfg struct {
	authToken string
}

func (c *testCfg) AuthToken() string {
	return c.authToken
}

func (c *testCfg) SetAuthToken(token string) error {
	c.authToken = token
	return nil
}
