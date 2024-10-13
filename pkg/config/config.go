package config

import "github.com/spf13/viper"

type Config struct {
	Server struct {
		Address      string
		ReadTimeout  int
		WriteTimeout int
	}

	Database struct {
		Host     string
		Port     int
		User     string
		Password string
		Name     string
	}

	JWT struct {
		SecretKey     string
		TokenLifetime int
	}

	Log struct {
		Level string
	}
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
