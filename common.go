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

func getPath(keyClass, keyName string) string {
	var path string

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

	return path + "/" + keyName
}

func getPathPriv(keyClass, keyName string) string {
	return getPath(keyClass, keyName) + ".key"
}

func getPathPub(keyClass, keyName string) string {
	return getPath(keyClass, keyName) + ".pub"
}

func getPathCert(keyClass, keyName string) string {
	return getPath(keyClass, keyName) + ".crt"
}
