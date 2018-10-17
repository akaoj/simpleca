package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"time"
)


const VERSION string = "1.2.1"


func getHelp() string {
	return `Usage: simpleca <action>

Available actions:
	generate
	init
	rm
	sign
	version`
}


func main() {
	output, err := run()
	if err != nil {
		fmt.Println("Error: " + err.Error())
		os.Exit(2)
	}

	fmt.Println(output)
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

	var action string = os.Args[1]

	var state State
	var conf Conf
	var err error
	var msg string

	// Some actions might be fired without being inside a repo
	switch action {
	case "init":
		return "Folder initialized, please edit the configuration.json file to fit your organization", init_()
	case "help":
		var topic string = ""

		if len(os.Args) == 3 {
			topic = os.Args[2]
		}

		switch topic {
		case "":
			return getHelp(), nil
		case "rm":
			return getHelpRm(), nil
		case "generate":
			return getHelpGenerate(), nil
		case "init":
			return getHelpInit(), nil
		case "sign":
			return getHelpSign(), nil
		default:
			return "", errors.New("the action \"" + topic + "\" has no help available\n\n" + getHelp())
		}
	case "version":
		return "simpleca v" + VERSION, nil
	}

	if ! isRepo() {
		return "", errors.New(`The current folder does not appear to be a valid simpleca repository.
Please run "simpleca init" before running any other command.
`)
	}

	// These will fail if we run them in an unitialized repo but here it's safe now
	state, err = loadState()
	if err != nil {
		return "", err
	}

	conf, err = getConfig()
	if err != nil {
		return "", err
	}

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

		err = generate(&state, conf, class, keySize, keyType, keyName, clearText)
		if err != nil {
			return "", err
		}
	case "rm":
		if len(os.Args[2:]) < 1 {
			return "", errors.New("missing class\n\n" + getHelpRm())
		}

		var class string = os.Args[2]
		var keyName string

		commands := flag.NewFlagSet("rm", flag.ExitOnError)

		commands.StringVar(&keyName, "name", "", "")

		commands.Parse(os.Args[3:])

		switch class {
		case "root":
			return "", errors.New(`can't delete a root key, this is too dangerous: all intermediate and client keys will become orphans (no way to revoke them or sign new intermediates certificates).
If you want to get rid of this CA, remove the whole folder (or better: create a new one next to this one in case you need the old CA someday).`)
		case "intermediate":
			fmt.Print("Warning! You are about to delete an intermediate key and certificate, are you sure you want to do that (y/N)? ")

			var answer string = "n";
			fmt.Scanln(&answer)

			if answer != "y" && answer != "yes" && answer != "Y" {
				return "", errors.New("Aborting")
			}
		case "client":
			// Do nothing
		default:
			return "", errors.New("can't delete a " + class)
		}

		err = rm(&state, conf, class, keyName)
		if err != nil {
			return "", err
		}

		msg = class + " keys and certificates deleted"
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
	default:
		return "", errors.New("the action \"" + action + "\" does not exist\n\n" + getHelp())
	}

	state.LastModificationDate = time.Now()

	err = saveState(state)
	if err != nil {
		return "", err
	}

	return msg, nil
}
