package main

import (
	"io/ioutil"
	"encoding/json"
	"os"
)


var folders [3]string = [...]string{"root", "intermediates", "clients"}


func getHelpInit() string {
	return `Usage: simpleca init

Init the current folder as a simpleca repository.
It will create root/, intermediates/ and clients/ folders as well as an empty state file and a generic configuration file.

You can run it multiple times, it won't overwrite or delete your content.`
}


// Do some cheap checks to know if we're in a simpleca folder
func isRepo() bool {
	var err error

	if _, err = os.Stat(statePath); os.IsNotExist(err) {
		return false
	}
	if _, err = os.Stat(confPath); os.IsNotExist(err) {
		return false
	}

	for _, f := range folders {
		if _, err = os.Stat(f); os.IsNotExist(err) {
			return false
		}
	}

	return true
}


// Can't call it init() because of go
func init_() error {
	var err error

	// Create and chmod all subfolders
	for _, f := range folders {
		err = os.Mkdir(f, 0700)
		if err != nil && ! os.IsExist(err) {
			return err
		}

		err = os.Chmod(f, 0700)
		if err != nil {
			return err
		}
	}

	// Init empty State and configuration if needed
	if _, err = os.Stat(statePath); os.IsNotExist(err) {
		f, err := os.Create(statePath)
		if err != nil {
			return err
		}

		_, err = f.WriteString("{}\n")
		if err != nil {
			return err
		}
	}

	if _, err = os.Stat(confPath); os.IsNotExist(err) {
		_, err := os.Create(confPath)
		if err != nil {
			return err
		}

		// No config file: create one
		var conf Conf = Conf{
			36,
			"SimpleCA",
			"France",
			"Paris",
		}

		b, err := json.MarshalIndent(conf, "", "    ")
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(confPath, b, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}
