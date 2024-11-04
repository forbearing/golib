package mqtt

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/forbearing/golib/config"
	"github.com/forbearing/golib/logger"
	"go.uber.org/zap"
)

var (
	client      mqtt.Client
	mu          sync.RWMutex
	initialized bool
	clientId    string
)

func Init() error {
	if !config.App.MqttConfig.Enable {
		return nil
	}

	mu.Lock()
	defer mu.Unlock()
	if initialized {
		return nil
	}

	opts, err := createClientOptions()
	if err != nil {
		return fmt.Errorf("create mqtt options failed: %w", err)
	}
	client = mqtt.NewClient(opts)
	if err := connect(client); err != nil {
		return fmt.Errorf("connect to mqtt broker failed: %w", err)
	}
	go monitorConnection()
	initialized = true
	return nil
}

func connect(client mqtt.Client) error {
	token := client.Connect()
	if !token.WaitTimeout(config.App.MqttConfig.ConnectTimeout) {
		return fmt.Errorf("connect timeout")
	}
	if err := token.Error(); err != nil {
		return err
	}
	zap.S().Infow("successfully connect to mqtt broker",
		"addr", config.App.MqttConfig.Addr,
		"client_id", clientId,
		"keepalive", config.App.Keepalive.String(),
		"connection_timeout", config.App.MqttConfig.ConnectTimeout.String(),
		"clean_session", config.App.MqttConfig.CleanSession,
		"auto_reconnect", config.App.MqttConfig.AutoReconnect,
	)
	return nil
}

func createClientOptions() (*mqtt.ClientOptions, error) {
	cfg := config.App.MqttConfig
	clientId = fmt.Sprintf("%s-%d",
		defaultIfEmpty(cfg.ClientPrefix, "mqtt-client"),
		rand.New(rand.NewSource(time.Now().UnixNano())).Int(),
	)

	opts := mqtt.NewClientOptions().
		AddBroker(cfg.Addr).
		SetAutoReconnect(cfg.AutoReconnect).
		SetClientID(clientId).
		SetProtocolVersion(5).
		SetKeepAlive(cfg.Keepalive).
		SetCleanSession(cfg.CleanSession).
		SetConnectTimeout(cfg.ConnectTimeout).
		SetCleanSession(cfg.CleanSession)
	if cfg.Username != "" {
		opts.SetUsername(cfg.Username)
		opts.SetPassword(cfg.Password)
	}
	if cfg.UseTLS {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: cfg.InsecureSkipVerify,
		}
		if len(cfg.CertFile) != 0 && len(cfg.KeyFile) != 0 {
			cert, err := loadCertificate(cfg.CertFile, cfg.KeyFile)
			if err != nil {
				return nil, fmt.Errorf("load certificate failed: %w", err)
			}
			tlsConfig.Certificates = []tls.Certificate{cert}
		}
		opts.SetTLSConfig(tlsConfig)
	}
	opts.SetConnectionLostHandler(func(c mqtt.Client, err error) {
		logger.Mqtt.Errorw("mqtt connection lost", "error", err, "client_id", clientId)
	})
	opts.SetOnConnectHandler(func(c mqtt.Client) {
		logger.Mqtt.Infow("mqtt client connected", "client_id", clientId)
	})

	return opts, nil
}

// loadCertificate loads TLS certificate
func loadCertificate(certFile, keyFile string) (tls.Certificate, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("load certificate failed: %w", err)
	}
	return cert, nil
}

func loadTLSConfig(caFile string) *tls.Config {
	// load tls config
	var tlsConfig tls.Config
	tlsConfig.InsecureSkipVerify = false
	if caFile != "" {
		certpool := x509.NewCertPool()
		ca, err := os.ReadFile(caFile)
		if err != nil {
			logger.Mqtt.Fatal(err.Error())
		}
		certpool.AppendCertsFromPEM(ca)
		tlsConfig.RootCAs = certpool
	}
	return &tlsConfig
}

// defaultIfEmpty returns default value if str is empty
func defaultIfEmpty(str, defaultStr string) string {
	if str == "" {
		return defaultStr
	}
	return str
}

func monitorConnection() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if !client.IsConnected() {
			logger.Mqtt.Warn("mqtt client disconnected, attempting to reconnect...")
			if err := connect(client); err != nil {
				logger.Mqtt.Errorw("reconnect failed", "error", err)
				continue
			}
			logger.Mqtt.Info("successfully reconnected")
		}
	}
}

// Client returns the MQTT client instance
func Client() (mqtt.Client, error) {
	mu.RLock()
	defer mu.RUnlock()

	if !initialized {
		return nil, fmt.Errorf("mqtt client not initialized")
	}
	if client == nil {
		return nil, fmt.Errorf("mqtt client is nil")
	}

	return client, nil
}

// Health checks if the MQTT client is connected
func Health() error {
	c, err := Client()
	if err != nil {
		return err
	}

	if !c.IsConnected() {
		return fmt.Errorf("mqtt client is not connected")
	}

	return nil
}

// Close closes the MQTT client connection
func Close() error {
	mu.Lock()
	defer mu.Unlock()

	if client != nil {
		client.Disconnect(250) // 等待 250ms 完成断开
		client = nil
		initialized = false
	}
	return nil
}

func Publish(topic string, payload any, opts ...PublishOption) error {
	c, err := Client()
	if err != nil {
		return fmt.Errorf("get mqtt client failed: %w", err)
	}
	opt := DefaultPublishOption
	if len(opts) > 0 {
		opt = opts[0]
	}

	var data []byte

	switch v := payload.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	default:
		var err error
		if data, err = json.Marshal(v); err != nil {
			return fmt.Errorf("marshal payload failed: %w", err)
		}
	}

	token := c.Publish(topic, opt.QoS, opt.Retain, data)
	if !token.WaitTimeout(opt.Timeout) {
		return fmt.Errorf("publish timeout")
	}
	if err := token.Error(); err != nil {
		logger.Mqtt.Errorw("publish failed",
			"error", err,
			"topic", topic,
			"addr", config.App.MqttConfig.Addr,
		)
		return err
	}
	logger.Mqtt.Debugw("publish success",
		"topic", topic,
		"payload", string(data),
		"qos", opt.QoS,
	)
	return nil
}

type MessageHandler func(topic string, payload []byte) error

func Subscribe(topic string, handler MessageHandler, opts ...SubscribeOption) error {
	c, err := Client()
	if err != nil {
		return fmt.Errorf("get mqtt client failed: %w", err)
	}
	opt := DefaultSubscribeOption
	if len(opts) > 0 {
		opt = opts[0]
	}
	wrapper := func(client mqtt.Client, msg mqtt.Message) {
		logger.Mqtt.Debugw("received message",
			"topic", msg.Topic(),
			"payload", string(msg.Payload()),
		)
		if handler != nil {
			if err := handler(msg.Topic(), msg.Payload()); err != nil {
				logger.Mqtt.Errorw("handle message failed",
					"error", err,
					"topic", msg.Topic(),
				)
			}
		}
	}

	token := c.Subscribe(topic, opt.QoS, wrapper)
	if !token.WaitTimeout(opt.Timeout) {
		return fmt.Errorf("subscribe timeout")
	}
	if err := token.Error(); err != nil {
		logger.Mqtt.Errorw("subscribe failed",
			"error", err,
			"topic", topic,
			"addr", config.App.MqttConfig.Addr,
		)
		return err
	}

	logger.Mqtt.Infow("subscribe success",
		"topic", topic,
		"qos", opt.QoS,
	)
	return nil
}

func Unsubscribe(topics ...string) error {
	c, err := Client()
	if err != nil {
		return fmt.Errorf("get mqtt client failed: %w", err)
	}

	token := c.Unsubscribe(topics...)
	if !token.WaitTimeout(5 * time.Second) {
		return fmt.Errorf("unsubscribe timeout")
	}
	if err := token.Error(); err != nil {
		logger.Mqtt.Errorw("unsubscribe failed",
			"error", err,
			"topics", topics,
		)
		return err
	}

	logger.Mqtt.Infow("unsubscribe success", "topics", topics)
	return nil
}
