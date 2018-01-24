package main

import (
	"crypto/rsa"
	"crypto/ecdsa"
)

const rootPath = "root"
const intermediatesPath = "intermediates"
const clientsPath = "clients"


func getPubKey(priv interface{}) interface{} {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	default:
		return nil
	}
}

func getPath(keyClass, keyName, keyType string) string {
	var extension string
	var path string

	switch keyType {
	case "priv":
		extension = "key"
	case "pub":
		extension = "pub"
	case "cert":
		extension = "crt"
	default:
		panic("the key type " + keyType + " is not valid")
	}

	switch keyClass {
	case "root":
		path = "root"
	case "intermediate":
		path = "intermediates"
	case "client":
		path = "clients"
	default:
		panic("the key class " + keyClass + " is not valid")
	}

	return path + "/" + keyName + "." + extension
}
