package main

import (
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"flag"
	"io/ioutil"
	"math/big"
	"os"
	"time"
)


func sign(args []string) error {
	if len(args) < 1 {
		return errors.New("missing class\n\n" + getHelpSign())
	}

	var class string = args[0]
	var keyName string
	var with string

	var err error

	commands := flag.NewFlagSet("sign", flag.ExitOnError)

	commands.StringVar(&keyName, "name", "", "")
	commands.StringVar(&with, "with", "", "")

	commands.Parse(args[1:])

	switch class {
	case "root":
		if keyName == "" {
			keyName = "root"
		}
	case "intermediate":
		if keyName == "" {
			keyName = "intermediate"
		}
	case "client":
		if keyName == "" {
			keyName = "client"
		}
	default:
		return errors.New("can't sign a " + class)
	}

	var privKeyPath string = getPath(class, keyName, "priv")
	var pubKeyPath string = getPath(class, keyName, "pub")

	if _, err = os.Stat(privKeyPath); os.IsNotExist(err) {
		return errors.New("the private key " + keyName + " does not exist")
	}
	if _, err = os.Stat(pubKeyPath); os.IsNotExist(err) {
		return errors.New("the public key " + keyName + " does not exist")
	}

	// Read file, decode it as PEM and load it
	privKeyBinary, err := ioutil.ReadFile(privKeyPath)
	if err != nil {
		return err
	}
	privKeyPem, _ := pem.Decode(privKeyBinary)

	var privKeyData, pubKeyData interface{}

	// TODO: find a way to know dynamically
	var keyType string = "rsa"

	if keyType == "rsa" {
		privKeyData, err = x509.ParsePKCS1PrivateKey(privKeyPem.Bytes)
		if err != nil {
			return err
		}
	} else if keyType == "ecdsa" {
		privKeyData, err = x509.ParseECPrivateKey(privKeyPem.Bytes)
		if err != nil {
			return err
		}
	} else {
		return errors.New("key type " + keyType + " is not implemented")
	}

	// Get public key from private
	pubKeyData = getPubKey(privKeyData)

	certFile, err := os.OpenFile(getPath(class, keyName, "cert"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer certFile.Close()

	var cert []byte

	if with == "" {
		// Self-signed certificate
		cert, err = x509.CreateCertificate(rand.Reader, getCA(), getCA(), pubKeyData, privKeyData)
		if err != nil {
			return err
		}

	} else {
		// TODO: sign with the given CA
		// TODO: make sure the target and the optional CA exists
	}

	pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: cert})

	return nil
}


func getCA() *x509.Certificate {
	// TODO: allow for configuration
	return &x509.Certificate{
		SerialNumber: big.NewInt(1653),
		Subject: pkix.Name{
			Organization:  []string{"ORGANIZATION_NAME"},
			Country:       []string{"COUNTRY_CODE"},
			Province:      []string{"PROVINCE"},
			Locality:      []string{"CITY"},
			StreetAddress: []string{"ADDRESS"},
			PostalCode:    []string{"POSTAL_CODE"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}
}

func getHelpSign() string {
	return `Usage: simpleca sign [class] --with=<ca name>

Available classes:
	intermediate   sign an intermediate CA public key
	client         sign a client public key

--with string
	(optional) Sign the key with the given object (this should be the name of an intermediate CA, or "root" if you want
	to sign an intermediate CA)
`
}
