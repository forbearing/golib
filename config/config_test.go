package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/forbearing/golib/bootstrap"
	"github.com/forbearing/golib/config"
	. "github.com/forbearing/golib/util"
	"github.com/stretchr/testify/assert"
)

var configData = `
[wechat]
app_id = "wx123456789"

[nats]
username = "nuser"
password = "npass"
; timeout = "30s"
enable = true
`

var filename = "/tmp/config.ini"

func TestRegisterStruct(t *testing.T) {
	if err := os.WriteFile(filename, []byte(configData), 0o644); err != nil {
		t.Fatal(err)
	}

	// Register config before bootstrap
	config.Register[WechatConfig]("wechat")

	config.SetConfigFile(filename)
	RunOrDie(bootstrap.Bootstrap)

	// Register config after bootstrap
	config.Register[NatsConfig]("nats")

	wechat := config.Get[*WechatConfig]("wechat")

	assert.Equal(t, "wx123456789", wechat.AppID)
	assert.Equal(t, "myappsecret", wechat.AppSecret)
	assert.Equal(t, false, wechat.Enable)

	nats := config.Get[NatsConfig]("nats")
	assert.Equal(t, "nats://127.0.0.1:4222", nats.URL)
	assert.Equal(t, "nuser", nats.Username)
	assert.Equal(t, "npass", nats.Password)
	assert.Equal(t, 5*time.Second, nats.Timeout)
	assert.Equal(t, true, nats.Enable)
}

func TestRegisterStructPointer(t *testing.T) {
	if err := os.WriteFile(filename, []byte(configData), 0o644); err != nil {
		t.Fatal(err)
	}

	// Register config before bootstrap
	config.Register[*WechatConfig]("wechat")

	config.SetConfigFile(filename)
	RunOrDie(bootstrap.Bootstrap)

	// Register config after bootstrap
	config.Register[*NatsConfig]("nats")

	wechat := config.Get[*WechatConfig]("wechat")

	assert.Equal(t, "wx123456789", wechat.AppID)
	assert.Equal(t, "myappsecret", wechat.AppSecret)
	assert.Equal(t, false, wechat.Enable)

	nats := config.Get[NatsConfig]("nats")
	assert.Equal(t, "nats://127.0.0.1:4222", nats.URL)
	assert.Equal(t, "nuser", nats.Username)
	assert.Equal(t, "npass", nats.Password)
	assert.Equal(t, 5*time.Second, nats.Timeout)
	assert.Equal(t, true, nats.Enable)
}

func TestRegisterStructFromEnv(t *testing.T) {
	if err := os.WriteFile(filename, []byte(configData), 0o644); err != nil {
		t.Fatal(err)
	}

	os.Setenv("NATS_USERNAME", "user_from_env")
	os.Setenv("NATS_PASSWORD", "pass_from_env")
	os.Setenv("NATS_TIMEOUT", "60s")

	// Register config before bootstrap
	config.Register[WechatConfig]("wechat")

	config.SetConfigFile(filename)
	RunOrDie(bootstrap.Bootstrap)

	// Register config after bootstrap
	config.Register[NatsConfig]("nats")

	wechat := config.Get[*WechatConfig]("wechat")

	assert.Equal(t, "wx123456789", wechat.AppID)
	assert.Equal(t, "myappsecret", wechat.AppSecret)
	assert.Equal(t, false, wechat.Enable)

	nats := config.Get[NatsConfig]("nats")
	assert.Equal(t, "nats://127.0.0.1:4222", nats.URL)
	assert.Equal(t, "user_from_env", nats.Username)
	assert.Equal(t, "pass_from_env", nats.Password)
	assert.Equal(t, 60*time.Second, nats.Timeout)
	assert.Equal(t, true, nats.Enable)
}

type WechatConfig struct {
	AppID     string `json:"app_id" mapstructure:"app_id" default:"myappid"`
	AppSecret string `json:"app_secret" mapstructure:"app_secret" default:"myappsecret"`
	Enable    bool   `json:"enable" mapstructure:"enable"`
}

type NatsConfig struct {
	URL      string        `json:"url" mapstructure:"url" default:"nats://127.0.0.1:4222"`
	Username string        `json:"username" mapstructure:"username" default:"nats"`
	Password string        `json:"password" mapstructure:"password" default:"nats"`
	Timeout  time.Duration `json:"timeout" mapstructure:"timeout" default:"5s"`
	Enable   bool          `json:"enable" mapstructure:"enable"`
}
