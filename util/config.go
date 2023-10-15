package util

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

type loaderFn func(conf *ConfigData, env string, ext string) error

type TlsConfigData struct {
	Enabled string `yaml:"enabled" json:"enabled" env:"GODRINK_SERVER_TLS_ENABLED"`
	Ca      string `yaml:"ca" json:"ca" env:"GODRINK_SERVER_TLS_CA"`
	Cert    string `yaml:"cert" json:"cert" env:"GODRINK_SERVER_TLS_CERT"`
	Key     string `yaml:"key" json:"key" env:"GODRINK_SERVER_TLS_KEY"`
}

type ServerConfigData struct {
	Host string         `yaml:"host" json:"host" env:"GODRINK_SERVER_HOST"`
	Port string         `yaml:"port" json:"port" env:"GODRINK_SERVER_PORT"`
	Tls  *TlsConfigData `yaml:"tls" json:"tls"`
}

type DatabaseTlsConfigData struct {
	Enabled    bool   `yaml:"enabled" json:"enabled" env:"GODRINK_DATABASE_TLS_ENABLED"`
	SkipVerify bool   `yaml:"skip-verify" json:"skipVerify" env:"GODRINK_DATABASE_TLS_SKIP_VERIFY"`
	Ca         string `yaml:"ca" json:"ca" env:"GODRINK_DATABASE_TLS_CA"`
	Cert       string `yaml:"cert" json:"cert" env:"GODRINK_DATABASE_TLS_CERT"`
	Key        string `yaml:"key" json:"key" env:"GODRINK_DATABASE_TLS_KEY"`
}

type DatabaseConfigData struct {
	Dialect string                 `yaml:"dialect" json:"dialect" env:"GODRINK_DATABASE_DIALECT"`
	Host    string                 `yaml:"host" json:"host" env:"GODRINK_DATABASE_HOST"`
	Port    string                 `yaml:"port" json:"port" env:"GODRINK_DATABASE_PORT"`
	User    string                 `yaml:"user" json:"user" env:"GODRINK_DATABASE_USER"`
	Pass    string                 `yaml:"pass" json:"pass" env:"GODRINK_DATABASE_PASS"`
	Name    string                 `yaml:"name" json:"name" env:"GODRINK_DATABASE_NAME"`
	Tls     *DatabaseTlsConfigData `yaml:"tls" json:"tls"`
}

type ConfigData struct {
	Environment string
	Server      *ServerConfigData   `yaml:"server" json:"server"`
	Database    *DatabaseConfigData `yaml:"database" json:"database"`
}

const configPrefix = "config"

var loaders = &map[string]loaderFn{
	"yml":  loadYaml,
	"yaml": loadYaml,
	"json": loadJson,
	"env":  loadEnvFile,
}

var config = &ConfigData{
	Environment: "development",
	Server: &ServerConfigData{
		Host: "localhost",
		Port: "8080",
		Tls: &TlsConfigData{
			Enabled: "false",
		},
	},
	Database: &DatabaseConfigData{
		Dialect: "postgres",
		Host:    "localhost",
		Port:    "5432",
		User:    "postgres",
		Pass:    "postgres",
		Name:    "postgres",
		Tls: &DatabaseTlsConfigData{
			Enabled:    true,
			SkipVerify: true,
		},
	},
}

func init() {
	environment := os.Getenv("GODRINK_ENV")
	for ext, loader := range *loaders {
		loader(config, "", ext)
		if environment != "" {
			loader(config, environment, ext)
		}
	}
}

func GetConfig() ConfigData {
	return *config
}

func loadYaml(conf *ConfigData, env string, ext string) error {
	if env != "" {
		env = fmt.Sprintf(".%s", env)
	}

	filename := fmt.Sprintf("%s%s.%s", configPrefix, env, ext)
	file, err := loadFile(JoinFromRoot(filename))
	if err != nil {
		defer file.Close()
		return fmt.Errorf("error loading YAML file %s: %w", filename, err)
	}
	if file != nil {
		return yaml.NewDecoder(file).Decode(conf)
	}
	return nil
}

func loadJson(conf *ConfigData, env string, ext string) error {
	if env != "" {
		env = fmt.Sprintf(".%s", env)
	}

	filename := fmt.Sprintf("%s%s.%s", configPrefix, env, ext)
	file, err := loadFile(JoinFromRoot(filename))
	if err != nil {
		return fmt.Errorf("error loading JSON file %s: %w", filename, err)
	}
	if file != nil {
		defer file.Close()
		return json.NewDecoder(file).Decode(conf)
	}
	return nil
}

func loadEnvFile(conf *ConfigData, env string, ext string) error {
	if env != "" {
		env = fmt.Sprintf(".%s", env)
	}

	filename := fmt.Sprintf(".%s%s", ext, env)
	if err := expectFileNotExist(godotenv.Load(JoinFromRoot(filename))); err != nil {
		return fmt.Errorf("error loading env file %s: %w", filename, err)
	}
	return envconfig.Process("GODRINK", conf)
}

func loadFile(path string) (*os.File, error) {
	file, err := os.Open(path)
	if err = expectFileNotExist(err); err != nil {
		return nil, fmt.Errorf("file open error: %w", err)
	}
	return file, nil
}

func expectFileNotExist(err error) error {
	if os.IsNotExist(err) {
		return nil
	}
	return err
}
