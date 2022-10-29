package config

import (
	"github.com/caarlos0/env/v6"
)

type Config struct {
	Env        string `env:"GO_ENV" envDefault:"production"`
	Port       int    `env:"PORT" envDefault:"80"`
	DBHost     string `env:"DB_HOST" envDefault:"db"`
	DBPort     int    `env:"DB_PORT" envDefault:"3306"`
	DBUser     string `env:"DB_USER" envDefault:"admin"`
	DBPassword string `env:"DB_PASSWORD" envDefault:"password"`
	DBName     string `env:"DB_NAME" envDefault:"point_app"`
	RedisHost  string `env:"REDIS_HOST" envDefault:"db"`
	RedisPort  int    `env:"REDIS_PORT" envDefault:"36379"`
}

// 環境変数の構造体を返却
//
// @return *Config 環境変数の構造体
//
// @return error エラー
func New() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
