package main

import (
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
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

	var privKey, pubKey interface{}

	var keyInState *Element
	var ok bool

	keyInState, ok = (*state).get(class, keyName)
	if !ok {
		return errors.New("key " + keyName + " is not known")
	}

	privKey, pubKey, err = loadPrivKey((*keyInState).Type, (*keyInState).Path)

	var cert []byte
	var certStruct *x509.Certificate

	var serial *big.Int

	// If we are signing a client key, this will tell the user about the fullchain file
	var additionalMessage string

	serial, err = rand.Int(rand.Reader, (&big.Int{}).Exp(big.NewInt(2), big.NewInt(159), nil))
	if err != nil {
		return err
	}

	if with == "" {
		// Self-signed certificate
		if class == "client" {
			certStruct = getCertForClient(serial, conf.CertificateDuration, keyName, conf.Organization, conf.Country, conf.Locality)
		} else {
			certStruct = getCertForCA(serial, conf.CertificateDuration, conf.Organization, conf.Country, conf.Locality)
		}

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
		var withCertificatePem *pem.Block
		var withCertificateX509 *x509.Certificate

		// Load the keys
		withPrivKey, _, err = loadPrivKey((*withElement).Type, (*withElement).Path)
		if err != nil {
			return err
		}

		withCertificatePem, withCertificateX509, err = loadCertificate((*withElement).Path)
		if err != nil {
			return err
		}

		if class == "client" {
			certStruct = getCertForClient(serial, conf.CertificateDuration, keyName, conf.Organization, conf.Country, conf.Locality)
		} else {
			certStruct = getCertForCA(serial, conf.CertificateDuration, conf.Organization, conf.Country, conf.Locality)
		}

		cert, err = x509.CreateCertificate(rand.Reader, certStruct, withCertificateX509, pubKey, withPrivKey)
		if err != nil {
			return err
		}

		// If this is a client key, create the full chain too
		if class == "client" {
			var fullchainCertPath = getCertPath((*keyInState).Path + ".fullchain")

			fullchainCertFile, err := os.OpenFile(fullchainCertPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
			if err != nil {
				return err
			}
			defer fullchainCertFile.Close()

			pem.Encode(fullchainCertFile, &pem.Block{Type: "CERTIFICATE", Bytes: cert})
			pem.Encode(fullchainCertFile, withCertificatePem)

			additionalMessage = "A full chain certificate file is also available at " + fullchainCertPath
		}
	}

	var certPath string = getCertPath((*keyInState).Path)

	certFile, err := os.OpenFile(certPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer certFile.Close()

	pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: cert})

	// Save the serial number
	el, ok := (*state).get(class, keyName)
	if !ok {
		return errors.New("Can't save serial number for " + class + " key " + keyName)
	}

	(*el).SerialNumber = (*serial).String()

	fmt.Println(keyName + " key signed, certificate available in " + certPath)
	if additionalMessage != "" {
		fmt.Println(additionalMessage)
	}

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
