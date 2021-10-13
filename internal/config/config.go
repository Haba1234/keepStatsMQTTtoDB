package config

import (
	"github.com/BurntSushi/toml"
)

// Config структура конфигурации.
type Config struct {
	Logger  LogConf           // Logger - конфигурация Loggera.
	MQTT    MQTTConf          // MQTT - конфигурация MQTT клиента.
	Servers map[string]Server // Servers - конфигурация серверов (брокеров) MQTT.
	Storage StorageConf       // Storage - конфигурация для подключения к БД.
}

// LogConf структура конфигурации.
type LogConf struct {
	Level string `toml:"log-level"` // Level - уровень логирования.
}

// MQTTConf структура конфигурации.
type MQTTConf struct {
	ClientID string `toml:"clientID"` // ClientID - имя клиента.
}

type Server struct {
	Schema   string           `toml:"schema"`   // Schema - тип подключения.
	Host     string           `toml:"server"`   // Host - адрес MQTT сервера.
	Port     string           `toml:"port"`     // Port - порт MQTT сервера.
	User     string           `toml:"user"`     // User - логин для подключения к MQTT серверу.
	Password string           `toml:"password"` // Password - пароль для подключения к MQTT серверу.
	Qos      byte             `toml:"qos"`      // Qos - качество обслуживания.
	Topics   map[string]Topic `toml:"topics"`   // Topics - слайс топиков для подписки.
}

type Topic struct {
	Measurement string `toml:"measurement"`
	Name        string `toml:"name"`
}

// StorageConf структура конфигурации.
type StorageConf struct {
	URL    string `toml:"db-uri"` // URL - IP:port базы данных.
	Bucket string `toml:"bucket"` // Bucket - параметры подключения.
	Org    string `toml:"org"`    // Org - параметры подключения.
	Token  string `toml:"token"`  // Token - параметры подключения.
}

// NewConfig конструктор.
func NewConfig(path string) (*Config, error) {
	// default values
	cfg := Config{
		Logger:  LogConf{},
		MQTT:    MQTTConf{},
		Servers: map[string]Server{},
		Storage: StorageConf{},
	}
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return &cfg, err
	}
	return &cfg, nil
}
