package validate

import (
	"encoding/json"
	"io/ioutil"
	"path"

	"github.com/whiteblock/genesis-cli/pkg/config"
	"github.com/whiteblock/genesis-cli/pkg/message"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	jschema "github.com/xeipuuv/gojsonschema"
)

var (
	conf             = config.NewConfig()
	ErrUnknownFormat = errors.New(message.UnknownFormat)
	ErrNoSchema      = errors.New(message.NoSchema)
)

func schemaFromCache() ([]byte, error) {
	if !conf.UserDir.Exists(conf.SchemaFile) {
		return nil, nil
	}
	return conf.UserDir.ReadFile(conf.SchemaFile)
}

func storeSchema(data []byte) error {
	return conf.UserDir.WriteFile(conf.SchemaFile, data)
}

func getSchema() ([]byte, error) {
	resp, err := conf.HTTPClient().Get(path.Join(conf.APIEndpoint(), conf.SchemaURI))
	if err != nil {
		return schemaFromCache()
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		log.WithField("error", storeSchema(data)).Debug("failed to store schema locally")
	}
	return data, err
}

func forceIntoJSON(in []byte) ([]byte, error) {
	if json.Valid(in) {
		return in, nil
	}

	out, err := yaml.YAMLToJSON(in)
	if err == nil {
		return out, nil
	}

	return nil, ErrUnknownFormat
}

func Schema(file []byte) (*jschema.Result, error) {
	data, err := forceIntoJSON(file)
	if err != nil {
		return nil, err
	}

	schema, err := getSchema()
	if err != nil {
		return nil, errors.Wrap(err, ErrNoSchema.Error())
	}
	if schema == nil {
		return nil, ErrNoSchema
	}

	schemaLoader := jschema.NewBytesLoader(schema)
	fileLoader := jschema.NewBytesLoader(data)
	return jschema.Validate(schemaLoader, fileLoader)
}
