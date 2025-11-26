package config

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"runtime"
)

type DBConn struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"db_name"`
	Port     int    `yaml:"port"`
}

type RedisConn struct {
	Host     string
	Port     int
	Password string
}

type API struct {
	Port int `yaml:"port"`
}

type Config struct {
	Api   API       `yaml:"api"`
	DB    DBConn    `yaml:"db"`
	Redis RedisConn `yaml:"redis"`
}

// use this if config needs any other validation
func (c *Config) Validate() error {
	if c.DB.Password == "" {
		return fmt.Errorf("DB password is required")
	}

	return nil
}

func configPath() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	return filepath.Join(dir, "config.yaml")
}

// use 12 factor config
func NewConfig() (*Config, error) {
	data, err := os.ReadFile(configPath())
	if err != nil {
		return nil, fmt.Errorf("failed to load default config %+v", err)
	}

	var cfg Config

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall default config")
	}

	dbPassword := flag.String("password", "", "Password for db connection")

	if !flag.Parsed() {
		flag.Parse()
	}

	// check if env has it, we can use tag based env resolving(with reflection) but I skip for now
	if dbName := os.Getenv("DB_NAME"); dbName != "" {
		cfg.DB.DBName = dbName
	}

	if dbPass := os.Getenv("DB_PASS"); dbPass != "" || *dbPassword != "" {
		cfg.DB.Password = dbPass
	}

	if dbUser := os.Getenv("DB_USER"); dbUser != "" {
		cfg.DB.User = dbUser
	}

	return &cfg, nil
}
