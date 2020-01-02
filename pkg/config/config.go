package config

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/shibukawa/configdir"
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
	RedirectURL  string        `mapstructure:"redirectURL"`

	Verbosity string `mapstructure:"verbosity"`

	MultipathUploadURI string        `mapstructure:"multipathUploadURI"`
	SchemaURI          string        `mapstructure:"schemaURI"`
	WBHost             string        `mapstructure:"wbHost"`
	APITimeout         time.Duration `mapstructure:"apiTimeout"`
	BiomeDNSZone       string        `mapstructure:"biomeDNSZone"`

	OrgID string `mapstructure:"orgID"`

	TokenFile  string `mapstructure:"tokenFile"`
	OrgFile    string `mapstructure:"orgFile"`
	SchemaFile string `mapstructure:"schemaFile"`

	Dir     configdir.ConfigDir `mapstructure:"-"`
	UserDir *configdir.Config   `mapstructure:"-"`
}

func (c Config) HTTPClient() *http.Client {
	return &http.Client{Timeout: c.APITimeout}
}

func (c Config) APIEndpoint() string {
	if !strings.HasPrefix(c.WBHost, "http") {
		return "https://" + c.WBHost
	}
	return c.WBHost
}

// GetLogger gets a logger according to the config
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
	viper.BindEnv("schemaURI", "SCHEMA_URI")

	viper.BindEnv("schemaFile", "SCHEMA_FILE")
	viper.BindEnv("tokenFile", "TOKEN_FILE")
	viper.BindEnv("orgFile", "ORG_FILE")
	viper.BindEnv("apiTimeout", "API_TIMEOUT")
	//"biomeDNSZone" no env binding
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

	viper.SetDefault("schemaURI", "/schemas/test-definition-format")
	viper.SetDefault("schemaFile", ".test-definition-format-schema")

	viper.SetDefault("apiTimeout", 5*time.Second)
	viper.SetDefault("biomeDNSZone", "infra.whiteblock.io")
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
