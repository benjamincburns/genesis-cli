package config

import (
	"sync"
	"time"

	"github.com/shibukawa/configdir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Config groups all of the global configuration parameters into
// a single struct
type Config struct {
	AuthEndpoint       string        `mapstructure:"authEndpoint"`
	AuthPath           string        `mapstructure:"authPath"`
	AuthTimeout        time.Duration `mapstructure:"authTimeout"`
	TokenPath          string        `mapstructure:"tokenPath"`
	Verbosity          string        `mapstructure:"verbosity"`
	RedirectURL        string        `mapstructure:"redirectURL"`
	MultipathUploadURI string        `mapstructure:"multipathUploadURI"`
	WBHost             string        `mapstructure:"wbHost"`
	OrgID              string        `mapstructure:"orgID"`
	TokenFile          string        `mapstructure:"tokenFile"`
	OrgFile            string        `mapstructure:"orgFile"`

	Dir     configdir.ConfigDir `mapstructure:"-"`
	UserDir *configdir.Config   `mapstructure:"-"`
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
	viper.BindEnv("MultipathUploadURI", "MULTIPART_UPLOAD_URI")
	viper.BindEnv("wbHost", "WB_HOST")
	viper.BindEnv("orgID", "ORG_ID")
	viper.BindEnv("tokenFile", "TOKEN_FILE")
	viper.BindEnv("orgFile", "ORG_FILE")
}

func setViperDefaults() {
	viper.SetDefault("authEndpoint", "auth.infra.whiteblock.io")
	viper.SetDefault("authPath", "/auth/realms/wb/protocol/openid-connect/auth")
	viper.SetDefault("tokenPath", "/auth/realms/wb/protocol/openid-connect/token")
	viper.SetDefault("redirectURL", "localhost:56666")
	viper.SetDefault("authTimeout", 120*time.Second)
	viper.SetDefault("verbosity", "INFO")
	viper.SetDefault("multipathUploadURI", "/api/v1/testexecution/organizations/%s/files")
	viper.SetDefault("wbHost", "www.infra.whiteblock.io")
	viper.SetDefault("orgID", "")
	viper.SetDefault("tokenFile", ".auth_token")
	viper.SetDefault("orgFile", ".org_name")
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
		conf.Dir = configdir.New("whiteblock", "genesis-cli")
		conf.UserDir = conf.Dir.QueryFolders(configdir.Global)[0]
		conf.UserDir.MkdirAll()
	})
	return conf
}
