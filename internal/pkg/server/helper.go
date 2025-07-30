package server

import (
	"path/filepath"
	"strings"

	"github.com/pachirode/iam_study/pkg/log"
	"github.com/pachirode/iam_study/pkg/utils/homedir"
	"github.com/spf13/viper"
)

func LoadConfig(cfg string, defaultName string) {
	if cfg != "" {
		viper.SetConfigFile(cfg)
	} else {
		viper.AddConfigPath(".")
		viper.AddConfigPath(filepath.Join(homedir.HomeDir(), RecommendedHomeDir))
		viper.SetConfigName(defaultName)
	}

	viper.SetConfigFile("yaml")
	viper.AutomaticEnv()
	viper.SetEnvPrefix(RecommendedEnvPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	if err := viper.ReadInConfig(); err != nil {
		log.Warnf("WARNING: viper failed to discover and load configuration file: %s", err.Error())
	}
}
