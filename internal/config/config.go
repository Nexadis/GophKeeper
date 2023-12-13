package config

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// HTTPServerConfig - конфиг HTTP сервера
type HTTPServerConfig struct {
	Up         bool   // Up - запускать ли HTTP-сервер
	Address    string // Address - адрес на котором сервер ждет подключений
	JWTSecret  []byte // JWTSecret - секрет для JWT
	TLS        bool   // TLS - использовать ли TLS
	CrtFile    string // CrtFile - путь до сертификата сервера
	KeyFile    string // KeyFile - путь до ключа сервера
	ClientsDir string // ClientsDir - директория с скомпилированными клиентами под разные архитектуры
	FrontDir   string // FrontDir - директория с Frontend'ом
}

// HTTPClientConfig - конфиг HTTP клиента
type HTTPClientConfig struct {
	Address string // Address - адрес подключения к серверу
	TLS     bool   // TLS - использовать ли TLS
	CrtFile string // CrtFile - путь до сертификата сервера для доверенного подключения
	Retries int    // Retries - количество повторных попыток для подключения
}

// DBConfig - конфиг для подключения к БД
type DBConfig struct {
	URI     string
	Timeout int64
}

// ServerConfig - общий конфиг сервера, независимо от транспорта
type ServerConfig struct {
	Debug  bool
	HTTP   *HTTPServerConfig
	DB     *DBConfig
	Log    *LogConfig
	WarmUp time.Duration
}

// ClientConfig - общий конфиг клиента, независимо от транспорта
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

// MustServerConfig - создаёт конфиг сервера и определяет его из значений по умолчанию, файла с конфигом и переменных окружения
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

// MustClientConfig - создаёт конфиг клиента и определяет его из значений по умолчанию, файла с конфигом и переменных окружения
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
