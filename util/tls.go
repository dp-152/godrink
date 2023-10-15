package util

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

type TlsParams struct {
	CaPath     string
	CertPath   string
	KeyPath    string
	SkipVerify bool
}

func MakeTlsConfig(config TlsParams) (*tls.Config, error) {
	certPool, err := GetCertPool(config.CaPath)
	if err != nil {
		return nil, fmt.Errorf("error getting cert pool: %w", err)
	}

	var clientCert *tls.Certificate
	if config.CertPath != "" && config.KeyPath != "" {
		cert, err := tls.LoadX509KeyPair(config.CertPath, config.KeyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load key pair for cert %s: %w", config.CertPath, err)
		}
		clientCert = &cert
	}

	return &tls.Config{
		RootCAs:            certPool,
		Certificates:       []tls.Certificate{*clientCert},
		InsecureSkipVerify: config.SkipVerify,
	}, nil
}

func GetCertPool(certPaths ...string) (*x509.CertPool, error) {
	certPool, _ := x509.SystemCertPool()
	if certPool == nil {
		certPool = x509.NewCertPool()
	}

	for _, certPath := range certPaths {
		pemBytes, err := os.ReadFile(certPath)
		if err != nil {
			return nil, fmt.Errorf("error loading pem blocks from ca file: %w", err)
		}
		certPool.AppendCertsFromPEM(pemBytes)
	}

	return certPool, nil
}

func LoadPemFile(pemPath string) ([]*pem.Block, error) {
	pemBytes, err := os.ReadFile(pemPath)
	if err != nil {
		return nil, fmt.Errorf("cannot read PEM file at %s: %w", pemPath, err)
	}

	pemBlocks := make([]*pem.Block, 0, 10)

	for block, rest := pem.Decode(pemBytes); block != nil; block, rest = pem.Decode(rest) {
		pemBlocks = append(pemBlocks, block)
	}

	return pemBlocks, nil
}
