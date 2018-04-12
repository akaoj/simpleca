package main

import (
	"errors"
	"os"
)


func getHelpRm() string {
	return `Usage: simpleca rm <class> [--name=<name>]

Delete a key pair and all associated certificates.

--name string
	(optional) The key name.`
}


func rm(state *State, conf Conf, class string, name string) error {
	var err error

	var path string

	switch class {
	case "root":
		return errors.New("can't delete a root key")
	case "intermediate":
		path = IntermediatesPath
		if name == "" {
			name = "intermediate"
		}
	case "client":
		path = ClientsPath
		if name == "" {
			name = "client"
		}
	default:
		return errors.New("can't delete a " + class)
	}

	var fullPath string = path + "/" + name

	var privKeyPath string = getPrivKeyPath(fullPath)
	var pubKeyPath string = getPubKeyPath(fullPath)
	var certPath string = getCertPath(fullPath)
	var fullCertPath string = getFullCertPath(fullPath)

	for _, file := range []string{privKeyPath, pubKeyPath, certPath, fullCertPath} {
		if _, err = os.Stat(file); err == nil {
			err = os.Remove(file)
			if err != nil {
				return err
			}
		}
	}

	if class == "client" {
		delete((*state).Clients, name)
	} else if class == "intermediate" {
		delete((*state).Intermediates, name)
	}

	return nil
}
