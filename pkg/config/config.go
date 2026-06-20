package config

import (
	"bytes"
	"embed"
	"strings"

	viper2 "github.com/spf13/viper"
)

//go:embed config.yaml
var configFile embed.FS

var v *viper2.Viper

var defaultJWTSecret = "change-me-in-production-use-a-long-random-string"

func init() {
	data, err := configFile.ReadFile("config.yaml")
	if err != nil {
		panic(err.Error())
	}
	v = viper2.New()
	v.SetConfigType("yaml")
	if err := v.ReadConfig(bytes.NewReader(data)); err != nil {
		panic(err.Error())
	}

	s, _ := Get("config.jwt.secret").(string)
	if s == "" || s == defaultJWTSecret {
		panic("FATAL: jwt.secret must be changed from the default value in production. Set it in config.yaml or via JWT_SECRET env var.")
	}
}

func Get(fileKey string) interface{} {
	index := strings.Index(fileKey, ".")
	if index == -1 {
		return v.Get(fileKey)
	}
	return v.Get(fileKey[index+1:])
}
