package config

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Config groups all of the global configuration parameters into
// a single struct
type Config struct {
	AuthEndpoint string        `mapstructure:"authEndpoint"`
	AuthPath     string        `mapstructure:"authPath"`
	AuthTimeout  time.Duration `mapstructure:"authTimeout"`
	TokenPath    string        `mapstructure:"tokenPath"`
	Verbosity    string        `mapstructure:"verbosity"`
	RedirectURL  string        `mapstructure:"redirectURL"`
}

//GetLogger gets a logger according to the config
func (c Config) setupLogrus() {
	lvl, err := logrus.ParseLevel(c.Verbosity)
	if err != nil {
		logrus.SetLevel(logrus.InfoLevel)
	} else {
		logrus.SetLevel(lvl)
	}

	return
}

func setViperEnvBindings() {
	viper.BindEnv("authEndpoint", "AUTH_ENDPOINT")
	viper.BindEnv("authPath", "AUTH_PATH")
	viper.BindEnv("tokenPath", "TOKEN_PATH")
	viper.BindEnv("verbosity", "VERBOSITY")
	viper.BindEnv("redirectURL", "REDIRECT_URL")
	viper.BindEnv("authTimeout", "AUTH_TIMEOUT")
}

func setViperDefaults() {
	viper.SetDefault("authEndpoint", "auth.infra.whiteblock.io")
	viper.SetDefault("authPath", "/auth/realms/wb/protocol/openid-connect/auth")
	viper.SetDefault("tokenPath", "/auth/realms/wb/protocol/openid-connect/token")
	viper.SetDefault("redirectURL", "localhost:56666")
	viper.SetDefault("authTimeout", 120)
	viper.SetDefault("verbosity", "INFO")
}

func init() {
	setViperDefaults()
	setViperEnvBindings()

	//viper.AddConfigPath("/etc/whiteblock/")          // path to look for the config file in
	//viper.AddConfigPath("$HOME/.config/whiteblock/") // call multiple times to add many search paths
	//viper.SetConfigName("genesis")
	//viper.SetConfigType("yaml")

}

var (
	conf Config
	once = &sync.Once{}
)

// NewConfig creates a new config object from the global config
func NewConfig() Config {
	once.Do(func() {
		_ = viper.ReadInConfig()
		err := viper.Unmarshal(&conf)
		if err != nil {
			panic(err)
		}
		conf.setupLogrus()
	})
	return conf
}
