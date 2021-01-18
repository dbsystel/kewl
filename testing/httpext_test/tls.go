package httpext_test

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/dbsystel/kewl/pkg/panicutils"
	"k8s.io/apimachinery/pkg/util/uuid"
)

var CACertificate = &x509.Certificate{
	SerialNumber: big.NewInt(1),
	Subject: pkix.Name{
		Organization:  []string{"DB Systel GmbH"},
		Country:       []string{"DE"},
		Province:      []string{"Hesse"},
		Locality:      []string{"Frankfurt am Main"},
		StreetAddress: []string{"Juergen-Ponto-Platz 1"},
		PostalCode:    []string{"60329"},
	},
	NotBefore:             time.Now().AddDate(0, 0, -1),
	NotAfter:              time.Now().AddDate(1, 0, 0),
	IsCA:                  true,
	ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
	KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	BasicConstraintsValid: true,
}

const maxTmpFileFind = 10

// PemFilePair denotes a pair of PEM files for a certificate
type PemFilePair interface {
	// PublicKeyPath returns the path to the public part
	PublicKeyPath() string
	// PrivateKeyPath returns the path to the private part
	PrivateKeyPath() string
	// Generate generates the files using the provided x509.Certificate
	Generate(cert *x509.Certificate) error
	// Delete deletes the files
	Delete() error
	// isAvailable returns true if both files do not exist (internal)
	isAvailable() bool
}

// NewPemFilePair creates a new PemFilePair for the provide base path (adding .pem and .key as extensions)
func NewPemFilePair(basePath string) PemFilePair {
	return &pemFilePairImpl{PublicPath: basePath + ".pem", PrivatePath: basePath + ".key"}
}

// NewTmpPemFilePair creates a new PemFilePair generating a name for it and using the
func NewTmpPemFilePair() (PemFilePair, error) {
	tmpDir := os.TempDir()
	if !strings.HasSuffix(tmpDir, "/") {
		tmpDir += "/"
	}
	var result PemFilePair
	for i := 0; i < maxTmpFileFind; i++ {
		if result = NewPemFilePair(fmt.Sprintf("%vgenerated-%v", tmpDir, uuid.NewUUID())); result.isAvailable() {
			return result, nil
		}
	}
	return nil, fmt.Errorf("unable to create temporary pem file pair after %v tries", maxTmpFileFind)
}

// MustNewTmpPemFilePair equals NewTmpPemFilePair, but panics in case of an error
func MustNewTmpPemFilePair() PemFilePair {
	result, err := NewTmpPemFilePair()
	panicutils.PanicIfError(err)
	return result
}
