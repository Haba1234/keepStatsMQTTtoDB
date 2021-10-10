package app

import (
	"keepStatsMQTTtoDB/internal/config"
)

// ClientMQTTConf структура конфигурации.
type ClientMQTTConf struct {
	ClientID string // ClientID - уникальное имя клиента для брокеров.

}

type ServerMQTTConf struct {
	Schema   string           // Schema - тип подключения.
	Host     string           // Host - адрес MQTT сервера.
	Port     string           // Port - порт MQTT сервера.
	User     string           // User - логин для подключения к MQTT серверу.
	Password string           // Password - пароль для подключения к MQTT серверу.
	Topics   map[string]Topic // Topics - карта топиков и QOS
}

// StorageConf структура конфигурации.
type StorageConf struct {
	URL    string // URL - IP:port базы данных.
	Bucket string // Bucket - параметры подключения.
	Org    string // Org - параметры подключения.
	Token  string // Token - параметры подключения.
}

// Point структура точки для передачи в БД.
type Point struct {
	Measurement string
	Tag         map[string]string
	Field       string
	Value       interface{}
}

type Topic struct {
	Measurement string
	Qos         byte
}

// ConvertConfigClientMQTT преобразует структуры.
func ConvertConfigClientMQTT(cfg config.MQTTConf) ClientMQTTConf {
	return ClientMQTTConf{
		ClientID: cfg.ClientID,
	}
}

// ConvertConfigServerMQTT преобразует структуры.
func ConvertConfigServerMQTT(cfg config.Server) ServerMQTTConf {
	return ServerMQTTConf{
		Schema:   cfg.Schema,
		Host:     cfg.Host,
		Port:     cfg.Port,
		User:     cfg.User,
		Password: cfg.Password,
		Topics:   sliceToStrMap(cfg.Qos, cfg.Topics),
	}
}

// ConvertConfigStorage преобразует структуры.
func ConvertConfigStorage(cfg config.StorageConf) StorageConf {
	return StorageConf{
		URL:    cfg.URL,
		Bucket: cfg.Bucket,
		Org:    cfg.Org,
		Token:  cfg.Token,
	}
}

// TODO добавить запись QOS из конфига.

func sliceToStrMap(qos byte, topics map[string]config.Topic) map[string]Topic {
	t := map[string]Topic{}

	for _, val := range topics {
		t[val.Name] = Topic{
			Measurement: val.Measurement,
			Qos:         qos,
		}
	}
	return t
}
