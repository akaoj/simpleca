package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"time"
)


const VERSION string = "0.2"


func main() {
	output, err := run()
	if err != nil {
		fmt.Println("Error: " + err.Error())
		os.Exit(2)
	}

	fmt.Print(output)
}

func run() (string, error) {
	if len(os.Args) < 2 {
		return "", errors.New("no action given\n\n" + getHelp())
	}

	var state State
	var conf Conf
	var err error

	state, err = loadState()
	if err != nil {
		return "", err
	}

	conf, err = getConfig()
	if err != nil {
		return "", err
	}

	var action string = os.Args[1]

	switch action {
	case "generate":
		if len(os.Args[2:]) < 1 {
			return "", errors.New("missing class\n\n" + getHelpGenerate())
		}

		var class string = os.Args[2]
		var keySize int
		var keyType string
		var keyName string

		commands := flag.NewFlagSet("generate", flag.ExitOnError)

		commands.StringVar(&keyType, "type", "", "")
		commands.IntVar(&keySize, "size", 0, "")
		commands.StringVar(&keyName, "name", "", "")

		commands.Parse(os.Args[3:])

		err := generate(&state, conf, class, keySize, keyType, keyName)
		if err != nil {
			return "", err
		}
	case "help":
		return getHelp(), nil
	case "info":
		return "", errors.New("the \"info\" action is not yet implemented")
	case "show":
		return "", errors.New("the \"list\" action is not yet implemented")
	case "sign":
		if len(os.Args[2:]) < 1 {
			return "", errors.New("missing class\n\n" + getHelpSign())
		}

		var class string = os.Args[2]
		var keyName string
		var with string

		commands := flag.NewFlagSet("sign", flag.ExitOnError)

		commands.StringVar(&keyName, "name", "", "")
		commands.StringVar(&with, "with", "", "")

		commands.Parse(os.Args[3:])

		err := sign(&state, conf, class, keyName, with)
		if err != nil {
			return "", err
		}
	case "version":
		return "simpleca v." + VERSION + "\n", nil
	default:
		return "", errors.New("the action \"" + action + "\" does not exist\n\n" + getHelp())
	}

	state.LastModificationDate = time.Now()

	err = saveState(state)
	if err != nil {
		return "", err
	}

	return "", nil
}
