package config

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type HTTPServerConfig struct {
	Up         bool
	Address    string
	JWTSecret  []byte
	TLS        bool
	CrtFile    string
	KeyFile    string
	ClientsDir string
	FrontDir   string
}

type HTTPClientConfig struct {
	Address string
	TLS     bool
	CrtFile string
	Retries int
}
type DBConfig struct {
	URI     string
	Timeout int64
}

type ServerConfig struct {
	Debug  bool
	HTTP   *HTTPServerConfig
	DB     *DBConfig
	Log    *LogConfig
	WarmUp time.Duration
}

type ClientConfig struct {
	Debug bool
	Log   *LogConfig
	HTTP  *HTTPClientConfig
}

type LogConfig struct {
	Level    string
	Outputs  []string
	Encoding string
}

func MustServerConfig() *ServerConfig {
	loadServerDefaults()
	loadConfig("server")
	setEnv()
	c := ServerConfig{}

	err := viper.Unmarshal(&c)
	if err != nil {
		log.Fatal(err)
	}

	return &c
}

func MustClientConfig() *ClientConfig {
	loadClientDefaults()
	loadConfig("client")
	setEnv()
	c := ClientConfig{}

	err := viper.Unmarshal(&c)
	if err != nil {
		log.Fatal(err)
	}
	return &c
}

func loadServerDefaults() {
	viper.SetDefault("debug", false)

	viper.SetDefault("http.up", true)
	viper.SetDefault("http.address", ":8443")
	viper.SetDefault("http.tls", true)
	viper.SetDefault("http.crtfile", "server.crt")
	viper.SetDefault("http.keyfile", "server.key")
	viper.SetDefault("http.frontdir", "frontend")
	viper.SetDefault("http.clientsdir", "clients")

	viper.SetDefault("db.uri", "postgresql://root:root@postgres:5432/keeper")
	viper.SetDefault("db.timeout", 10)

	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.outputs", []string{"stdout"})
	viper.SetDefault("log.encoding", "json")

	viper.SetDefault("warmup", 0)

	err := viper.WriteConfigAs("example_config.yaml")
	if err != nil {
		fmt.Println(err)
	}

}

func loadClientDefaults() {
	viper.SetDefault("debug", false)

	viper.SetDefault("http.address", "localhost:8443")
	viper.SetDefault("http.tls", true)
	viper.SetDefault("http.crtfile", "server.crt")
	viper.SetDefault("http.retries", 5)

	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.outputs", []string{"client.log"})
	viper.SetDefault("log.encoding", "console")

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
