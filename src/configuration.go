package main

import (
	"encoding/json"
	"io/ioutil"
)


const confPath = "configuration.json"


type Conf struct {
	CertificateDuration int
	Organization string
	Country string
	Locality string
}


func getConfig() (Conf, error) {
	var conf Conf

	confFile, err := ioutil.ReadFile(confPath)
	if err != nil {
		return Conf{}, err
	}

	err = json.Unmarshal(confFile, &conf)
	if err != nil {
		return Conf{}, err
	}

	return conf, nil
}
