package httpext

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/pkg/errors"
)

// TLSConfig denotes the TLS configuration for Server
type TLSConfig struct {
	// MinVersion is the minimum version string (e.g. 1.2)
	MinVersion string
	// CaCertsFile is the path to the CA certificates file to be used
	CaCertsFile string
	// PrivateKeyFile is the key file for the private key
	PrivateKeyFile string
	// PublicKeyFile is the key file for the public key
	PublicKeyFile string
}

// Config denotes the configuration parameters for Server
type Config struct {
	// ServerNamespace is a name for the server to be used as namespace for the prometheus Metrics
	ServerNamespace string
	// Addr is the address to listen to
	Addr string
	// TLSConfig denotes the TLSConfig to be used
	TLSConfig TLSConfig
	// ReadTimeout denotes the timeout for reading the complete request
	ReadTimeout time.Duration
	// ReadHeaderTimeout denotes the timeout for reading the headers
	ReadHeaderTimeout time.Duration
	// WriteTimeout denotes the timeout for writing the complete response
	WriteTimeout time.Duration
	// IdleTimeout denotes the timeout for an idling request
	IdleTimeout time.Duration
	// MaxHeaderBytes denotes the maximum acceptable size of the request header in bytes
	MaxHeaderBytes int
}

// NewDefaultConfig returns a new Config with the default options set
func NewDefaultConfig() Config {
	return Config{
		Addr: ":8443",
		TLSConfig: TLSConfig{
			MinVersion:     "1.2",
			CaCertsFile:    "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt",
			PrivateKeyFile: "/etc/ssl/private/tls.key",
			PublicKeyFile:  "/etc/ssl/private/tls.crt",
		},
		ReadTimeout:    10 * time.Second, // nolint:gomnd
		WriteTimeout:   10 * time.Second, // nolint:gomnd
		MaxHeaderBytes: 10 << 10,         // nolint:gomnd
	}
}

// AddCommandLineFlags adds all flags to flag.CommandLine for parsing the provided Config, using it's values as default
func AddCommandLineFlags(result *Config) {
	AddFlagsToFlagSet(result, flag.CommandLine)
}

// AddFlagsToFlagSet works like AddCommandLineFlags but let's you specify the target flag.FlagSet
func AddFlagsToFlagSet(result *Config, flagSet *flag.FlagSet) {
	flagSet.StringVar(
		&(result.Addr), "listen.addr",
		result.Addr, "the listen address for the web-serverImpl",
	)
	flagSet.StringVar(
		&(result.TLSConfig.MinVersion), "tls.min",
		result.TLSConfig.MinVersion, "the minimum tls version to be used",
	)
	flagSet.StringVar(
		&(result.TLSConfig.CaCertsFile), "tls.ca.certs",
		result.TLSConfig.CaCertsFile, "the path to the file containing the ca certificates",
	)
	flagSet.StringVar(
		&(result.TLSConfig.PrivateKeyFile), "tls.private.key",
		result.TLSConfig.PrivateKeyFile, "the path to the file containing the private key for the server",
	)
	flagSet.StringVar(
		&(result.TLSConfig.PublicKeyFile), "tls.public.key",
		result.TLSConfig.PublicKeyFile, "the path to the file containing the public key for the server",
	)
	flagSet.DurationVar(
		&(result.ReadTimeout), "timeout.read",
		result.ReadTimeout, "the timeout for reading the whole request",
	)
	flagSet.DurationVar(
		&(result.ReadHeaderTimeout), "timeout.read.header",
		result.ReadHeaderTimeout, "the timeout for reading the header",
	)
	flagSet.DurationVar(
		&(result.WriteTimeout), "timeout.write",
		result.WriteTimeout, "the timeout for writing the response",
	)
	flagSet.DurationVar(
		&(result.IdleTimeout), "timeout.idle",
		result.IdleTimeout, "the timeout for idle requests",
	)
	flagSet.IntVar(
		&(result.MaxHeaderBytes), "max.header.bytes",
		result.MaxHeaderBytes, "the maximum header bytes",
	)
}

// CreateTLSConfig creates a tls.Config from our TLSConfig
func CreateTLSConfig(config TLSConfig) (*tls.Config, error) {
	if config.MinVersion == "" && config.CaCertsFile == "" &&
		config.PublicKeyFile == "" && config.PrivateKeyFile == "" {
		return nil, nil
	}

	minVersion, err := resolveTLSVersion(config.MinVersion)
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	if len(config.CaCertsFile) > 0 {
		caCert, err := ioutil.ReadFile(config.CaCertsFile)
		if err != nil {
			return nil, errors.Wrapf(err, "could not read ca certs from file: %v", config.CaCertsFile)
		}
		if !caCertPool.AppendCertsFromPEM(caCert) {
			return nil, fmt.Errorf("could not add certificates from PEM file: %v", config.CaCertsFile)
		}
	}

	return &tls.Config{RootCAs: caCertPool, MinVersion: minVersion}, nil //nolint:gosec
}

// resolveTLSVersion resolves a version string to the corresponding constant in the tls package
func resolveTLSVersion(version string) (uint16, error) {
	switch version {
	case "1.2":
		return tls.VersionTLS12, nil
	case "1.3":
		return tls.VersionTLS13, nil
	}
	return 0, fmt.Errorf("invalid TLS version: %v", version)
}
