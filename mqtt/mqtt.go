package mqtt

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/forbearing/golib/config"
	"github.com/forbearing/golib/logger"
)

var Client mqtt.Client

func Init() error {
	timeout := 10 * time.Second
	logger.Mqtt.Infow("connecting mqtt broker", "addr", config.App.MqttConfig.Addr)
	connectAddress := config.App.MqttConfig.Addr
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	clientID := fmt.Sprintf("agv-backend-%d", r.Int())

	opts := mqtt.NewClientOptions()
	opts.AutoReconnect = true
	opts.ProtocolVersion = 5
	opts.ConnectTimeout = timeout
	opts.CleanSession = true
	opts.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	opts.AddBroker(connectAddress)
	// opts.SetUsername(config.App.MqttConfig.Username)
	// opts.SetPassword(config.App.MqttConfig.Password)
	opts.SetClientID(clientID)
	opts.SetKeepAlive(time.Second * 60)

	// Optional: 设置CA证书
	// opts.SetTLSConfig(loadTLSConfig("caFilePath"))

	Client = mqtt.NewClient(opts)
	if token := Client.Connect(); token.WaitTimeout(timeout) && token.Error() != nil {
		return token.Error()
	}
	logger.Mqtt.Infow("successfully connected mqtt broker", "addr", config.App.MqttConfig.Addr)

	go func() {
		// 检查 Client 是否处理连接状态，如果已经断开了，则尝试重连
		time.Sleep(5 * time.Second)
		for {
			if !Client.IsConnected() {
				logger.Mqtt.Warn("mqtt client disconnected, try to reconnect...")
				if token := Client.Connect(); token.WaitTimeout(timeout) && token.Error() != nil {
					logger.Mqtt.Error("failed to reconnect: ", token.Error())
				} else {
					logger.Mqtt.Info("successfully reconnected")
				}
			}
			time.Sleep(3 * time.Second)
		}
	}()

	return nil
}

func Publish(topic string, payload any) {
	payload, _ = json.Marshal(payload)
	qos := 0
	if token := Client.Publish(topic, byte(qos), false, payload); token.Wait() && token.Error() != nil {
		logger.Mqtt.Errorw(token.Error().Error(), "topic", topic, "addr", config.App.MqttConfig.Addr)
	} else {
		logger.Mqtt.Debugw("publish success", "topic", topic, "payload", payload)
	}
}

func Subscribe(topic string, handlers ...mqtt.MessageHandler) {
	_handler := func(client mqtt.Client, msg mqtt.Message) {
		logger.Mqtt.Infow("received", "topic", topic, "payload", msg.Payload())
	}
	if len(handlers) > 0 {
		_handler = handlers[0]
	}
	qos := 0
	Client.Subscribe(topic, byte(qos), _handler)
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
