package main

import (
	"bufio"
	"crypto/rsa"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"syscall"
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


func getPrivKeyPath(path string) string {
	return path + ".key"
}
func getPubKeyPath(path string) string {
	return path + ".pub"
}
func getCertPath(path string) string {
	return path + ".crt"
}


// Load private key file and return both private and public keys
func loadPrivKey(keyType, path string) (interface{}, interface{}, error) {
	var err error

	var privKeyBytes []byte
	var privKeyPem *pem.Block
	var privKeyDecryptedBytes []byte
	var privKey, pubKey interface{}

	var privKeyPath string = getPrivKeyPath(path)

	if _, err = os.Stat(privKeyPath); os.IsNotExist(err) {
		return privKey, pubKey, errors.New("the private key " + privKeyPath + " does not exist")
	}

	// Read file, decode it as PEM and load it
	privKeyBytes, err = ioutil.ReadFile(privKeyPath)
	if err != nil {
		return privKey, pubKey, err
	}
	privKeyPem, _ = pem.Decode(privKeyBytes)

	if x509.IsEncryptedPEMBlock(privKeyPem) {
		var password string

		password, err = getpass("The file " + privKeyPath + " is encrypted, please enter the password to unlock it: ")
		if err != nil {
			return privKey, pubKey, err
		}

		privKeyDecryptedBytes, err = x509.DecryptPEMBlock(privKeyPem, []byte(password))
		if err != nil {
			return privKey, pubKey, err
		}
	} else {
		privKeyDecryptedBytes = privKeyPem.Bytes
	}

	switch keyType {
	case "rsa":
		privKey, err = x509.ParsePKCS1PrivateKey(privKeyDecryptedBytes)
		if err != nil {
			return privKey, pubKey, err
		}
	case "ecdsa":
		privKey, err = x509.ParseECPrivateKey(privKeyDecryptedBytes)
		if err != nil {
		return privKey, pubKey, err
		}
	default:
		return privKey, pubKey, errors.New("key type " + keyType + " is not implemented")
	}

	pubKey = getPubKey(privKey)

	return privKey, pubKey, nil
}


// Load certificate
func loadCertificate(path string) (certificatePem *pem.Block, certificateX509 *x509.Certificate, err error) {
	var rawCertificateBytes []byte

	var certPath string = getCertPath(path)

	if _, err = os.Stat(certPath); os.IsNotExist(err) {
		return certificatePem, certificateX509, errors.New("certificate " + certPath + " does not exist")
	}

	// Load the certificate
	rawCertificateBytes, err = ioutil.ReadFile(certPath)
	if err != nil {
		return
	}

	certificatePem, _ = pem.Decode(rawCertificateBytes)

	certificateX509, err = x509.ParseCertificate(certificatePem.Bytes)
	if err != nil {
			return
	}

	return
}


// Thank you go for not providing a getpass() equivalent in the stdlib
func getpass(prompt string) (string, error) {
	var err error
	var ws syscall.WaitStatus

	// Print the prompt if needed
	if prompt != "" {
		fmt.Print(prompt)
	}

	attrs := syscall.ProcAttr{
		Dir: "",
		Env: []string{},
		Files: []uintptr{os.Stdin.Fd()},
		Sys: nil,
	}

	pid, err := syscall.ForkExec("/bin/stty", []string{"stty", "-echo"}, &attrs)
	if err != nil {
		return "", err
	}

	_, err = syscall.Wait4(pid, &ws, 0, nil)
	if err != nil {
		return "", err
	}

	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	pid, err = syscall.ForkExec("/bin/stty", []string{"stty", "echo"}, &attrs)
	if err != nil {
		return "", err
	}

	_, err = syscall.Wait4(pid, &ws, 0, nil)
	if err != nil {
		return "", err
	}

	// Print the carriage return that has been swallowed
	fmt.Println()

	return strings.TrimSpace(text), nil
}
