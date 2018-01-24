package main

import (
	"errors"
	"fmt"
	"os"
	"time"
)


const VERSION string = "0.1"


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
	var err error

	state, err = loadState()
	if err != nil {
		return "", err
	}

	var action string = os.Args[1]

	switch action {
	case "generate":
		err := generate(&state, os.Args[2:])
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
		err := sign(&state, os.Args[2:])
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
