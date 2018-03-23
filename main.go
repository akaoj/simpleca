package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"time"
)


const VERSION string = "1.1"


func getHelp() string {
	return `Usage: simpleca <action>

Available actions:
	generate
	sign
	version
`
}


func main() {
	output, err := run()
	if err != nil {
		fmt.Print("Error: " + err.Error())
		os.Exit(2)
	}

	fmt.Print(output)
}


type stringArray []string

func (i *stringArray) String() string {
	return ""
}

func (s *stringArray) Set(value string) error {
	*s = append(*s, value)
	return nil
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
		var clearText bool = false

		commands := flag.NewFlagSet("generate", flag.ExitOnError)

		commands.StringVar(&keyType, "type", "", "")
		commands.IntVar(&keySize, "size", 0, "")
		commands.StringVar(&keyName, "name", "", "")
		commands.BoolVar(&clearText, "clear-text", false, "")

		commands.Parse(os.Args[3:])

		err := generate(&state, conf, class, keySize, keyType, keyName, clearText)
		if err != nil {
			return "", err
		}
	case "help":
		var topic string = ""

		if len(os.Args) == 3 {
			topic = os.Args[2]
		}

		switch topic {
		case "":
			return getHelp(), nil
		case "generate":
			return getHelpGenerate(), nil
		case "sign":
			return getHelpSign(), nil
		default:
			return "", errors.New("the action \"" + topic + "\" has no help available\n\n" + getHelp())
		}
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

		var altNames stringArray

		commands.StringVar(&keyName, "name", "", "")
		commands.StringVar(&with, "with", "", "")
		commands.Var(&altNames, "altname", "")

		commands.Parse(os.Args[3:])

		err := sign(&state, conf, class, keyName, with, altNames)
		if err != nil {
			return "", err
		}
	case "rm":
		return "", errors.New("the \"rm\" action is not yet implemented")
	case "version":
		return "simpleca v" + VERSION + "\n", nil
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
