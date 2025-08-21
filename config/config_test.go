package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/forbearing/golib/config"
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
	config.Register[Wechat]()
	config.SetConfigFile(filename)
	if err := config.Init(); err != nil {
		t.Fatal(err)
	}
	// Register config after bootstrap
	config.Register[Nats]()

	wechat := config.Get[*Wechat]()
	assert.Equal(t, "wx123456789", wechat.AppID)
	assert.Equal(t, "myappsecret", wechat.AppSecret)
	assert.Equal(t, false, wechat.Enable)

	nats := config.Get[Nats]()
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
	config.Register[*Wechat]()
	config.SetConfigFile(filename)
	if err := config.Init(); err != nil {
		t.Fatal(err)
	}
	// Register config after bootstrap
	config.Register[*Nats]()
	wechat := config.Get[*Wechat]()

	assert.Equal(t, "wx123456789", wechat.AppID)
	assert.Equal(t, "myappsecret", wechat.AppSecret)
	assert.Equal(t, false, wechat.Enable)

	nats := config.Get[Nats]()
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

	os.Setenv("WECHAT_APP_SECRET", "my_app_secret")
	os.Setenv("NATS_USERNAME", "user_from_env")
	os.Setenv("NATS_PASSWORD", "pass_from_env")
	os.Setenv("NATS_TIMEOUT", "60s")

	// Register config before bootstrap
	config.Register[Wechat]()
	config.SetConfigFile(filename)
	if err := config.Init(); err != nil {
		t.Fatal(err)
	}
	// Register config after bootstrap
	config.Register[Nats]()

	wechat := config.Get[*Wechat]()

	assert.Equal(t, "wx123456789", wechat.AppID)
	assert.Equal(t, "my_app_secret", wechat.AppSecret)
	assert.Equal(t, false, wechat.Enable)

	nats := config.Get[Nats]()
	assert.Equal(t, "nats://127.0.0.1:4222", nats.URL)
	assert.Equal(t, "user_from_env", nats.Username)
	assert.Equal(t, "pass_from_env", nats.Password)
	assert.Equal(t, 60*time.Second, nats.Timeout)
	assert.Equal(t, true, nats.Enable)
}

func TestRegisterNonStructType(t *testing.T) {
	// These should be skipped silently without error or panic
	config.Register[string]()
	config.Register[int]()
	config.Register[*string]()
	config.Register[[]string]()
	config.Register[map[string]string]()

	// Should not panic or cause errors
	config.SetConfigFile(filename)
	if err := config.Init(); err != nil {
		t.Fatal(err)
	}
	// Getting non-registered configs should return zero values
	strVal := config.Get[string]()
	assert.Equal(t, "", strVal)

	intVal := config.Get[int]()
	assert.Equal(t, 0, intVal)
}

type Wechat struct {
	AppID     string `json:"app_id" mapstructure:"app_id" default:"myappid"`
	AppSecret string `json:"app_secret" mapstructure:"app_secret" default:"myappsecret"`
	Enable    bool   `json:"enable" mapstructure:"enable"`
}

type Nats struct {
	URL      string        `json:"url" mapstructure:"url" default:"nats://127.0.0.1:4222"`
	Username string        `json:"username" mapstructure:"username" default:"nats"`
	Password string        `json:"password" mapstructure:"password" default:"nats"`
	Timeout  time.Duration `json:"timeout" mapstructure:"timeout" default:"5s"`
	Enable   bool          `json:"enable" mapstructure:"enable"`
}

type TestConfig struct {
	Value string `json:"value" mapstructure:"value" default:"default_value"`
}
