package client

import (
	"encoding/json"
	"io/ioutil"
)

// conf holds the paths to all critical files and settings
type conf struct {
	ApiKeyFilePath        string
	MACSecretFilePath     string
	PackageConfigFilePath string
	MiddlemanUploadURL    string
	StatsFile             string
}

// newConf reads and unmarshals the conf.json file
func newConf(path string) (conf, error) {
	c := conf{}
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return c, err
	}
	err = json.Unmarshal(b, &c)
	if err != nil {
		return c, err
	}
	return c, nil
}
