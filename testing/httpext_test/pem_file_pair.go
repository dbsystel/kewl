package httpext_test

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"os"
	"syscall"

	"github.com/pkg/errors"
)

var _ PemFilePair = &pemFilePairImpl{}

// pemFilePairImpl is the implementation for PemFilePair
type pemFilePairImpl struct {
	// PublicPath is the path to the public part
	PublicPath string
	// PrivatePath is the path to the private path
	PrivatePath string
}

func (k *pemFilePairImpl) Delete() error {
	if err := k.deleteIfExists(k.PublicPath); err != nil {
		return errors.Wrapf(err, "could not delete public key: %v", k.PublicPath)
	}
	if err := k.deleteIfExists(k.PrivatePath); err != nil {
		return errors.Wrapf(err, "could not delete private key: %v", k.PublicPath)
	}
	return nil
}

func (k *pemFilePairImpl) Generate(cert *x509.Certificate) error {
	if cert == nil {
		return errors.New("certificate may not be nil")
	}
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return errors.Wrapf(err, "could not generate a private key for: %v", cert)
	}
	publicBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &privateKey.PublicKey, privateKey)
	if err != nil {
		return errors.Wrapf(err, "could not create certificate from key pair for: %v", cert)
	}
	privateBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	if err := k.storePublicKey(publicBytes); err != nil {
		return errors.Wrapf(err, "could not store public key for: %v", cert)
	}
	if err := k.storePrivateKey(privateBytes); err != nil {
		return errors.Wrapf(err, "could not store private key for: %v", cert)
	}
	return nil
}

func (k *pemFilePairImpl) PublicKeyPath() string {
	return k.PublicPath
}

func (k *pemFilePairImpl) PrivateKeyPath() string {
	return k.PrivatePath
}

func (k *pemFilePairImpl) isAvailable() bool {
	return !(k.exists(k.PublicPath) || k.exists(k.PrivatePath))
}

func (k *pemFilePairImpl) storePublicKey(representation []byte) error {
	return k.store(k.PublicPath, representation, "CERTIFICATE")
}

func (k *pemFilePairImpl) storePrivateKey(representation []byte) error {
	return k.store(k.PrivatePath, representation, "RSA PRIVATE KEY")
}

func (k *pemFilePairImpl) store(path string, representation []byte, keyType string) error {
	buffer := new(bytes.Buffer)
	if err := pem.Encode(buffer, &pem.Block{Type: keyType, Bytes: representation}); err != nil {
		return errors.Wrapf(err, "could not encode to PEM type: %v", keyType)
	}
	if err := ioutil.WriteFile(path, buffer.Bytes(), 0600); err != nil {
		return errors.Wrapf(err, "could not store to file: %v", keyType)
	}
	return nil
}

func (k *pemFilePairImpl) exists(path string) bool {
	_, err := os.Stat(path)
	return !k.IsNotFileExistsErr(err)
}

func (k *pemFilePairImpl) deleteIfExists(path string) error {
	if !k.exists(path) {
		return nil
	}
	return os.Remove(path)
}

func (k *pemFilePairImpl) IsNotFileExistsErr(err error) bool {
	if err == nil {
		return false
	}
	if os.IsNotExist(err) {
		return true
	}
	if pathErr, ok := err.(*os.PathError); ok {
		return pathErr.Err.(syscall.Errno) == syscall.ENOENT
	}
	return false
}
