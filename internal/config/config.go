package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

type HTTPConfig struct {
	Up      bool
	Address string
	TLS     bool
	CrtFile string
	KeyFile string
}
type DBConfig struct {
	URI string
}

type AppConfig struct {
	Debug bool
	HTTP  *HTTPConfig
	DB    *DBConfig
	Log   *LogConfig
}

type LogConfig struct {
	Level    string
	Outputs  []string
	Encoding string
}

func MustConfig() *AppConfig {
	loadDefaults()
	loadConfig("config")
	setEnv()
	c := AppConfig{}

	err := viper.Unmarshal(&c)
	if err != nil {
		log.Fatal(err)
	}

	return &c
}

func loadDefaults() {
	viper.SetDefault("debug", false)

	viper.SetDefault("http.up", true)
	viper.SetDefault("http.address", ":8080")
	viper.SetDefault("http.tls", false)
	viper.SetDefault("http.crtfile", "server.crt")
	viper.SetDefault("http.keyfile", "server.key")

	viper.SetDefault("db.uri", "postgresql://root:root@postgres:5432/keeper")

	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.outputs", []string{"stdout"})
	viper.SetDefault("log.encoding", "json")

	viper.SetDefault("crt", "server.crt")

	err := viper.WriteConfigAs("example_config.yaml")
	if err != nil {
		fmt.Println(err)
	}

}

func setEnv() {
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix("keeper")
	viper.AutomaticEnv()
}

func loadConfig(c string) {
	viper.SetConfigName(c)
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Error reading config file, %s\n", err)
	}
}
