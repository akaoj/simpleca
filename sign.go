package main

import (
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"math/big"
	"os"
	"time"
)


func sign(state *State, conf Conf, class, keyName, with string) error {
	var err error

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

	var privKeyPath string = getPathPriv(class, keyName)
	var pubKeyPath string = getPathPub(class, keyName)

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

	var privKey, pubKey interface{}

	var keyMemory *Element
	var ok bool

	keyMemory, ok = (*state).get(class, keyName)
	if !ok {
		return errors.New("key " + keyName + " is not known")
	}


	var keyType string = (*keyMemory).Type

	if keyType == "rsa" {
		privKey, err = x509.ParsePKCS1PrivateKey(privKeyPem.Bytes)
		if err != nil {
			return err
		}
	} else if keyType == "ecdsa" {
		privKey, err = x509.ParseECPrivateKey(privKeyPem.Bytes)
		if err != nil {
			return err
		}
	} else {
		return errors.New("key type " + keyType + " is not implemented")
	}

	// Get public key from private
	pubKey = getPubKey(privKey)

	certFile, err := os.OpenFile((*keyMemory).Path + ".crt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer certFile.Close()

	var cert []byte
	var certStruct *x509.Certificate

	var serial *big.Int

	serial, err = rand.Int(rand.Reader, (&big.Int{}).Exp(big.NewInt(2), big.NewInt(159), nil))
	if err != nil {
		return err
	}

	if with == "" {
		// Self-signed certificate
		certStruct = getCertForCA(serial, conf.CertificateDuration, conf.Organization, conf.Country, conf.Locality)

		cert, err = x509.CreateCertificate(rand.Reader, certStruct, certStruct, pubKey, privKey)
		if err != nil {
			return err
		}

	} else {
		var withElement *Element

		// Retrieve the element first from intermediate CAs, else from the root CA
		withElement, ok = (*state).get("intermediate", with)
		if !ok {
			withElement, ok = (*state).get("root", with)
			if !ok {
				return errors.New("can't find a CA named " + with)
			}
		}

		var withPrivKey interface{}
		var withCertificate *x509.Certificate

		// Load the keys
		withPrivKey, _, err = loadKeyPair(withElement.Type, (*withElement).Path)
		if err != nil {
			return err
		}

		// Load the certificate
		withCertificateBytes, err := ioutil.ReadFile((*withElement).Path + ".crt")
		if err != nil {
			return err
		}

		withCertificatePem, _ := pem.Decode(withCertificateBytes)

		withCertificate, err = x509.ParseCertificate(withCertificatePem.Bytes)
		if err != nil {
			return err
		}

		if class == "client" {
			certStruct = getCertForClient(serial, conf.CertificateDuration, keyName, conf.Organization, conf.Country, conf.Locality)
		} else {
			certStruct = getCertForCA(serial, conf.CertificateDuration, conf.Organization, conf.Country, conf.Locality)
		}

		cert, err = x509.CreateCertificate(rand.Reader, certStruct, withCertificate, pubKey, withPrivKey)
		if err != nil {
			return err
		}
	}

	pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: cert})

	// Save the serial number
	el, ok := (*state).get(class, keyName)
	if !ok {
		return errors.New("Can't save serial number for " + class + " key " + keyName)
	}

	(*el).SerialNumber = (*serial).String()

	return nil
}


func getCertForCA(serial *big.Int, duration int, organization, country, locality string) *x509.Certificate {
	return &x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			Organization:  []string{organization},
			Country:       []string{country},
			Locality:      []string{locality},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(0, duration, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}
}

func getCertForClient(serial *big.Int, duration int, commonName, organization, country, locality string) *x509.Certificate {
	return &x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			Organization:  []string{organization},
			Country:       []string{country},
			Locality:      []string{locality},
			CommonName:    commonName,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(0, duration, 0),
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature,
	}
}


func loadKeyPair(keyType, keyPath string) (interface{}, interface{}, error) {
	var err error

	var privKeyBytes []byte
	var privKeyPem *pem.Block
	var privKey, pubKey interface{}

	if _, err = os.Stat(keyPath + ".key"); os.IsNotExist(err) {
		return privKey, pubKey, errors.New("the private key " + keyPath + ".key does not exist")
	}
	if _, err = os.Stat(keyPath + ".pub"); os.IsNotExist(err) {
		return privKey, pubKey, errors.New("the public key " + keyPath + ".pub does not exist")
	}

	// Read file, decode it as PEM and load it
	privKeyBytes, err = ioutil.ReadFile(keyPath + ".key")
	if err != nil {
		return privKey, pubKey, err
	}
	privKeyPem, _ = pem.Decode(privKeyBytes)

	switch keyType {
	case "rsa":
		privKey, err = x509.ParsePKCS1PrivateKey(privKeyPem.Bytes)
		if err != nil {
			return privKey, pubKey, err
		}
	case "ecdsa":
		privKey, err = x509.ParseECPrivateKey(privKeyPem.Bytes)
		if err != nil {
			return privKey, pubKey, err
		}
	default:
		return privKey, pubKey, errors.New("key type " + keyType + " is not implemented")
	}

	pubKey = getPubKey(privKey)

	return privKey, pubKey, nil
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
