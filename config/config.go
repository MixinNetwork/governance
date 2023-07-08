package config

import (
	_ "embed"
	"encoding/json"

	"github.com/pelletier/go-toml"
)

const (
	BuildVersion = "BUILD_VERSION"
)

//go:embed config.toml
var configtoml []byte

//go:embed apps.json
var appsjson []byte

type Configuration struct {
	Database struct {
		Path string `toml:"path"`
	} `toml:"database"`
	Mixin struct {
		ClientID   string `toml:"client-id"`
		SessionID  string `toml:"session-id"`
		PrivateKey string `toml:"private-key"`
		PinToken   string `toml:"pin-token"`
		Pin        string `toml:"pin"`
	} `toml:"mixin"`
	Governance struct {
		FeeAssetID string `toml:"fee-asset-id"`
		Fee        string `toml:"fee"`
	} `toml:"governance"`
	Environment string `toml:"environment"`
	Port        string `toml:"port"`
}

var AppConfig *Configuration

func InitConfiguration(env string) {
	var cfs map[string]*Configuration

	err := toml.Unmarshal(configtoml, &cfs)
	if err != nil {
		panic(err)
	}
	AppConfig = cfs[env]
	if AppConfig == nil {
		panic("Invalid environment")
	}
}

type App struct {
	AppID      string `json:"app_id"`
	SessionID  string `json:"session_id"`
	PrivateKey string `json:"private_key"`
	PinToken   string `json:"pin_token"`
	Pin        string `json:"pin"`
}

func FetchApps() ([]*App, error) {
	var apps []*App
	err := json.Unmarshal(appsjson, &apps)
	if err != nil {
		return nil, err
	}
	return apps, nil
}
