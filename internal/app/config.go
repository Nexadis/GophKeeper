package app

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

type HTTPConfig struct {
	Address string
}
type DBConfig struct {
	URI string
}

type AppConfig struct {
	Debug bool
	HTTP  *HTTPConfig
	DB    *DBConfig
	Crt   string
}

func MustConfig() *AppConfig {
	loadDefaults()
	setEnv()
	loadConfig("config")
	c := AppConfig{}

	err := viper.Unmarshal(&c)
	if err != nil {
		log.Fatal(err)
	}

	if c.Debug {
		fmt.Println(viper.AllSettings())
	}
	return &c
}

func loadDefaults() {
	viper.SetDefault("debug", false)

	viper.SetDefault("http.address", ":8080")

	viper.SetDefault("db.uri", "postgresql://root:root@postgres:5432/keeper")

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
