package server

import (
	"net"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	RecommendedHomeDir   = ".iam_study"
	RecommendedEnvPrefix = "IAM"
)

type CertKey struct {
	CertFile string
	KeyFile  string
}

type SecureServingInfo struct {
	BindAddress string
	BindPort    int
	CertKey     CertKey
}

func (s *SecureServingInfo) Address() string {
	return net.JoinHostPort(s.BindAddress, strconv.Itoa(s.BindPort))
}

type InsecureServingInfo struct {
	Address string
}

type JwtInfo struct {
	Realm      string
	Key        string
	Timeout    time.Duration
	MaxRefresh time.Duration
}

type Config struct {
	SecureServingInfo   *SecureServingInfo
	InsecureServingInfo *InsecureServingInfo
	Jwt                 *JwtInfo
	Mode                string
	Middlewares         []string
	Healthz             bool
	EnableProfiling     bool
	EnableMetrics       bool
}

type CompletedConfig struct {
	*Config
}

func NewConfig() *Config {
	return &Config{
		Healthz:         true,
		Mode:            gin.ReleaseMode,
		Middlewares:     []string{},
		EnableProfiling: true,
		EnableMetrics:   true,
		Jwt: &JwtInfo{
			Realm:      "iam jwt",
			Timeout:    1 * time.Hour,
			MaxRefresh: 1 * time.Hour,
		},
	}
}

func (config *Config) Complete() CompletedConfig {
<<<<<<< HEAD
	return CompletedConfig{c}
}

func (config *Config) New()
=======
	return CompletedConfig{config}
}

func (config CompletedConfig) New() (*GenericAPIServer, error) {
	server := &GenericAPIServer{
		SecureServingInfo:   config.SecureServingInfo,
		InsecureServingInfo: config.InsecureServingInfo,
		mode:                config.Mode,
		healthz:             config.Healthz,
		enableMetrics:       config.EnableMetrics,
		enableProfiling:     config.EnableProfiling,
		middlewares:         config.Middlewares,
		Engine:              gin.New(),
	}

	initGenericAPIServer(server)

	return server, nil
}
>>>>>>> 4f132fc (add apiserver run options)
