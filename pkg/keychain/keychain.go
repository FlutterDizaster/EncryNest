package keychain

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"log/slog"
	"math/big"
	"os"
	"path/filepath"
	"time"
)

const (
	CertTTL   = time.Hour * 24 * 365
	KeyLenght = 2048
)

// LoadTLSCertificateSettings - settings for loading TLS certificate.
// If directory is empty, it will be ignored.
// CertFile and KeyFile must be set.
// If GeneratingSettings is not nil, new certificate will be generated if it does not exist.
type TLSCertificateSettings struct {
	Directory string
	CertFile  string
	KeyFile   string
}

func LoadTLSCertificate(settings TLSCertificateSettings) (*tls.Certificate, error) {
	var keyFilePath string
	var certFilePath string

	if settings.Directory != "" {
		keyFilePath = filepath.Join(settings.Directory, settings.KeyFile)
		certFilePath = filepath.Join(settings.Directory, settings.CertFile)
	} else {
		keyFilePath = settings.KeyFile
		certFilePath = settings.CertFile
	}

	privateKey, err := LoadPrivateKeyPEM(keyFilePath)
	if err != nil {
		return nil, err
	}

	certificate, err := LoadX509CertificatePEM(certFilePath)
	if err != nil {
		return nil, err
	}

	cert := &tls.Certificate{
		Certificate: [][]byte{certificate},
		PrivateKey:  privateKey,
	}

	return cert, nil
}

func GenerateTLSCertificate(
	settings TLSCertificateSettings,
	generationSettings GenerateCertSettings,
) (*tls.Certificate, error) {
	// Generating certificates
	certificate, privateKey, err := GenerateX509Certificate(generationSettings)
	if err != nil {
		return nil, err
	}

	// Saving certificates
	var keyFilePath string
	var certFilePath string

	if settings.Directory != "" {
		keyFilePath = filepath.Join(settings.Directory, settings.KeyFile)
		certFilePath = filepath.Join(settings.Directory, settings.CertFile)
	} else {
		keyFilePath = settings.KeyFile
		certFilePath = settings.CertFile
	}

	err = SaveX509CertificatePEM(certificate, certFilePath)
	if err != nil {
		return nil, err
	}

	err = SavePrivateKeyPEM(privateKey, keyFilePath)
	if err != nil {
		return nil, err
	}

	cert := &tls.Certificate{
		Certificate: [][]byte{certificate},
		PrivateKey:  privateKey,
	}

	return cert, nil
}

type GenerateCertSettings struct {
	Subject  pkix.Name
	Issuer   pkix.Name
	DNSNames []string
}

func GenerateX509Certificate(settings GenerateCertSettings) ([]byte, *rsa.PrivateKey, error) {
	// Generatink private key
	privKey, err := rsa.GenerateKey(rand.Reader, KeyLenght)
	if err != nil {
		slog.Error("Error while generating private key", slog.Any("err", err))
		return nil, nil, err
	}

	// Generating certificate data
	notBefore := time.Now()
	notAfter := notBefore.Add(CertTTL)

	// Generating random serial number
	serialNumber, err := rand.Int(rand.Reader, big.NewInt(1<<63-1))
	if err != nil {
		slog.Error("Error while generating serial number", slog.Any("err", err))
		return nil, nil, err
	}

	// Constructing template
	certTemplate := x509.Certificate{
		SerialNumber: serialNumber,
		Subject:      settings.Subject,
		Issuer:       settings.Issuer,
		NotBefore:    notBefore,
		NotAfter:     notAfter,
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		DNSNames:     settings.DNSNames,
	}

	cert, err := x509.CreateCertificate(
		rand.Reader,
		&certTemplate,
		&certTemplate,
		&privKey.PublicKey,
		privKey,
	)
	if err != nil {
		slog.Error("Error while creating certificate", slog.Any("err", err))
		return nil, nil, err
	}

	return cert, privKey, nil
}

func SavePrivateKeyPEM(privKey *rsa.PrivateKey, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	if err = pem.Encode(
		file,
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privKey),
		},
	); err != nil {
		return err
	}

	return nil
}

func SaveX509CertificatePEM(cert []byte, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	if err = pem.Encode(file, &pem.Block{Type: "CERTIFICATE", Bytes: cert}); err != nil {
		return err
	}

	return nil
}

func LoadPrivateKeyPEM(path string) (*rsa.PrivateKey, error) {
	keyPEM, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyPEM)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func LoadX509CertificatePEM(path string) ([]byte, error) {
	certPEM, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(certPEM)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the certificate")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}

	return cert.Raw, nil
}
