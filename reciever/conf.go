package reciever

import (
	"encoding/json"
	"io/ioutil"
)

type conf struct {
	MacSecretFile        string
	PackageConfigFile    string
	ApiKeyFile           string
	MiddlemanDownloadURL string
}

func newConf(p string) (conf, error) {
	b, err := ioutil.ReadFile(p)
	if err != nil {
		return conf{}, err
	}

	c := conf{}
	err = json.Unmarshal(b, &c)
	if err != nil {
		return conf{}, err
	}
	return c, nil
}
