package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"
)

const (
	ECDSA = "ecdsa"
	RSA = "rsa"
)

const statePath = "state.json"

type Element struct {
	Path string
	Type string
	Size int
	CreatedOn time.Time
	ValidUntil time.Time
}

type State struct {
	Root map[string]Element
	Intermediates map[string]Element
	Clients map[string]Element
	LastModificationDate time.Time
}


func (s *State) set(class, name string, value Element) {
	switch class {
	case "root":
		if s.Root == nil {
			s.Root = make(map[string]Element)
		}
		s.Root[name] = value
	case "intermediate":
		if s.Intermediates == nil {
			s.Intermediates = make(map[string]Element)
		}
		s.Intermediates[name] = value
	case "client":
		if s.Clients == nil {
			s.Clients = make(map[string]Element)
		}
		s.Clients[name] = value
	}
}

func (s *State) get(class, name string) Element {
	switch class {
	case "root":
		return s.Root[name]
	case "intermediate":
		return s.Intermediates[name]
	case "client":
		return s.Clients[name]
	}

	return Element{}
}


func loadState() (State, error) {
	var s State

	if _, err := os.Stat(statePath); os.IsNotExist(err) {
		f, err := os.Create(statePath)
		if err != nil {
			return State{}, err
		}

		_, err = f.WriteString("{}\n")
		if err != nil {
			return State{}, err
		}
	}

	stateFile, err := ioutil.ReadFile(statePath)
	if err != nil {
		return State{}, err
	}

	err = json.Unmarshal(stateFile, &s)
	if err != nil {
		return State{}, err
	}

	return s, nil
}

func saveState(s State) error {
	b, err := json.Marshal(s)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(statePath, b, 0644)

	return err
}
