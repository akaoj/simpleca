package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
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

	if _, err := os.Stat(confPath); os.IsNotExist(err) {
		_, err := os.Create(confPath)
		if err != nil {
			return Conf{}, err
		}

		// No config file: create one
		conf = Conf{
			36,
			"SimpleCA",
			"France",
			"Paris",
		}

		b, err := json.MarshalIndent(conf, "", "    ")
		if err != nil {
			return Conf{}, err
		}

		err = ioutil.WriteFile(confPath, b, 0644)
		if err != nil {
			return Conf{}, err
		}

		return conf, nil
	}

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
