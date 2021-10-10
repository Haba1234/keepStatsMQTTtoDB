package clientmqtt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"keepStatsMQTTtoDB/internal/app"
)

// ClientMQTT структура клиента MQTT.
type ClientMQTT struct {
	ctx        context.Context
	cfgClient  app.ClientMQTTConf
	client     mqtt.Client
	opts       *mqtt.ClientOptions
	pointsCh   chan<- app.Point
	serverName string
	server     app.ServerMQTTConf
}

// NewClient конструктор.
func NewClient(cfgClient app.ClientMQTTConf, serverName string, cgfServ app.ServerMQTTConf) *ClientMQTT {
	return &ClientMQTT{
		cfgClient:  cfgClient,
		serverName: serverName,
		server:     cgfServ,
	}
}

func (c *ClientMQTT) Start(ctx context.Context, pointsCh chan<- app.Point) error {
	mqtt.ERROR = log.New(os.Stdout, "[ERROR] ", 0)
	mqtt.CRITICAL = log.New(os.Stdout, "[CRIT] ", 0)
	mqtt.WARN = log.New(os.Stdout, "[WARN]  ", 0)

	c.ctx = ctx // TODO удалить?
	c.pointsCh = pointsCh

	if len(c.server.Topics) == 0 {
		return errors.New("no topics are specified in the configuration file")
	}

	// TODO добавить проверки по брокеру на ошибки конфигурации.
	c.opts = mqtt.NewClientOptions().
		AddBroker(fmt.Sprintf("%s://%s:%s", c.server.Schema, c.server.Host, c.server.Port)).
		SetUsername(c.server.User).
		SetPassword(c.server.Password).
		SetDefaultPublishHandler(c.messageHandler).
		SetOnConnectHandler(c.connectHandler).
		SetConnectionLostHandler(c.connectLostHandler).
		SetClientID(c.cfgClient.ClientID).
		SetOrderMatters(false).
		SetCleanSession(true).
		SetAutoReconnect(true).
		SetConnectRetry(true).
		SetConnectRetryInterval(5 * time.Second). // TODO добавить в конфиг.
		SetMaxReconnectInterval(5 * time.Second).
		SetKeepAlive(30 * time.Second)

	c.client = mqtt.NewClient(c.opts)

	token := c.client.Connect()
	select {
	case <-token.Done():
		if token.Error() != nil {
			return token.Error()
		}
		break
	case <-c.ctx.Done():
		return errors.New("context canceled")
	}

	if err := c.sub(c.client, c.server.Topics); err != nil {
		return err
	}

	log.Println("Status: ", c.client.IsConnected())
	return nil
}

func (c *ClientMQTT) Stop() error {
	if c.client != nil && c.client.IsConnected() {
		c.client.Disconnect(500)
	}
	return nil
}

func (c *ClientMQTT) sub(client mqtt.Client, topics map[string]app.Topic) error {
	filters := map[string]byte{}
	for name, val := range topics {
		filters[name] = val.Qos
	}

	token := client.SubscribeMultiple(filters, nil)
	select {
	case <-c.ctx.Done():
		return nil
	case <-token.Done():
		if token.Error() != nil {
			return token.Error()
		}
		log.Printf("Subscribed to topics %v", filters)
	}
	return nil
}

func (c *ClientMQTT) connectHandler(_ mqtt.Client) {
	log.Println("client connected to server:", c.serverName)
}

func (c *ClientMQTT) connectLostHandler(_ mqtt.Client, err error) {
	log.Printf("server: %s. Connect lost: %v", c.serverName, err)
}

func (c *ClientMQTT) messageHandler(_ mqtt.Client, msg mqtt.Message) {
	log.Printf("received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	go c.sendPointToStorage(msg)
}

func (c *ClientMQTT) sendPointToStorage(msg mqtt.Message) {
	message := msg
	var m interface{}
	if err := json.Unmarshal(message.Payload(), &m); err != nil {
		log.Printf("server: %s. Message could not be parsed (%s): %s", c.serverName, message.Payload(), err)
	}

	val, ok := c.server.Topics[message.Topic()]
	if !ok {
		log.Println("accepted topic was not found in the database. Recording in DB was canceled")
		return
	}

	tag := map[string]string{
		"server": c.serverName,
		"topic":  message.Topic(),
	}
	p := app.Point{
		Measurement: val.Measurement,
		Tag:         tag,
		Field:       "Value",
		Value:       m,
	}
	c.pointsCh <- p
}

// TODO TLS
/* https://webdevelop.pro/blog/using-golang-and-mqtt
func (c *MQTTConnector) configureMqttConnection() {
	connOpts := MQTT.NewClientOptions().
		AddBroker(c.config.URL).
		SetClientID(c.clientID).
		SetCleanSession(true).
		SetConnectionLostHandler(c.onConnectionLost).
		SetOnConnectHandler(c.onConnected).
		SetAutoReconnect(false) // we take care of re-connect ourselves

	// Username/password authentication
	if c.config.Username != "" && c.config.Password != "" {
		connOpts.SetUsername(c.config.Username)
		connOpts.SetPassword(c.config.Password)
	}

	// SSL/TLS
	if strings.HasPrefix(c.config.URL, "ssl") {
		tlsConfig := &tls.Config{}
		// Custom CA to auth broker with a self-signed certificate
		if c.config.CaFile != "" {
			caFile, err := ioutil.ReadFile(c.config.CaFile)
			if err != nil {
				logger.Printf("MQTTConnector.configureMqttConnection() ERROR: failed to read CA file %s:%s\n", c.config.CaFile, err.Error())
			} else {
				tlsConfig.RootCAs = x509.NewCertPool()
				ok := tlsConfig.RootCAs.AppendCertsFromPEM(caFile)
				if !ok {
					logger.Printf("MQTTConnector.configureMqttConnection() ERROR: failed to parse CA certificate %s\n", c.config.CaFile)
				}
			}
		}
		// Certificate-based client authentication
		if c.config.CertFile != "" && c.config.KeyFile != "" {
			cert, err := tls.LoadX509KeyPair(c.config.CertFile, c.config.KeyFile)
			if err != nil {
				logger.Printf("MQTTConnector.configureMqttConnection() ERROR: failed to load client TLS credentials: %s\n",
					err.Error())
			} else {
				tlsConfig.Certificates = []tls.Certificate{cert}
			}
		}

		connOpts.SetTLSConfig(tlsConfig)
	}

	c.client = MQTT.NewClient(connOpts)
} */
