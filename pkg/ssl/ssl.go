package ssl

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"github.com/go-acme/lego/v4/certcrypto"
	"math/big"
	"os"
	"time"
)

// Constants for all key types we support.
const (
	EC256   = KeyType("P256")
	EC384   = KeyType("P384")
	RSA2048 = KeyType("2048")
	RSA4096 = KeyType("4096")
	RSA8192 = KeyType("8192")
)

type SelfSSL struct {
	Domains          []string
	CommonName       string
	Country          string
	Organization     string
	OrganizationUint string
	Name             string
	KeyType          string
	Province         string
	City             string
	CertificatePath  string
	PrivateKeyPath   string
}

func GenerateSelfPem(selfSSL SelfSSL) error {

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	pkixName := pkix.Name{
		CommonName:         selfSSL.CommonName,
		Country:            []string{selfSSL.Country},
		Organization:       []string{selfSSL.Organization},
		OrganizationalUnit: []string{selfSSL.OrganizationUint},
	}
	if selfSSL.Province != "" {
		pkixName.Province = []string{selfSSL.Province}
	}
	if selfSSL.City != "" {
		pkixName.Locality = []string{selfSSL.City}
	}

	rootCA := &x509.Certificate{
		SerialNumber:          big.NewInt(time.Now().Unix()),
		Subject:               pkixName,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLen:            1,
		MaxPathLenZero:        false,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		DNSNames:              selfSSL.Domains,
	}

	// 生成证书
	derBytes, err := x509.CreateCertificate(rand.Reader, rootCA, rootCA, &privateKey.PublicKey, privateKey)
	if err != nil {
		return err
	}

	// 写入证书文件
	certOut, err := os.Create(selfSSL.CertificatePath)
	if err != nil {
		return err
	}
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()
	// fmt.Println("证书文件生成成功: cert.pem")

	// 生成私钥文件
	keyOut, err := os.OpenFile(selfSSL.PrivateKeyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	pem.Encode(keyOut, &pem.Block{Type: "EC PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})
	keyOut.Close()
	// fmt.Println("私钥文件生成成功: key.pem")
	return nil
}

func GeneratePrivateKey(keyType KeyType) (crypto.PrivateKey, error) {
	switch keyType {
	case EC256:
		return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	case EC384:
		return ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	case RSA2048:
		return rsa.GenerateKey(rand.Reader, 2048)
	case RSA4096:
		return rsa.GenerateKey(rand.Reader, 4096)
	case RSA8192:
		return rsa.GenerateKey(rand.Reader, 8192)
	}

	return nil, fmt.Errorf("invalid KeyType: %s", keyType)
}

func createPrivateKey(keyType string) (privateKey any, publicKey any, privateKeyBytes []byte, err error) {
	privateKey, err = certcrypto.GeneratePrivateKey(KeyType(keyType))
	if err != nil {
		return
	}
	var (
		caPrivateKeyPEM = new(bytes.Buffer)
	)
	if KeyType(keyType) == certcrypto.EC256 || KeyType(keyType) == certcrypto.EC384 {
		publicKey = &privateKey.(*ecdsa.PrivateKey).PublicKey
		publicKey = publicKey.(*ecdsa.PublicKey)
		block := &pem.Block{
			Type: "EC PRIVATE KEY",
		}
		privateBytes, sErr := x509.MarshalECPrivateKey(privateKey.(*ecdsa.PrivateKey))
		if sErr != nil {
			err = sErr
			return
		}
		block.Bytes = privateBytes
		_ = pem.Encode(caPrivateKeyPEM, block)
	} else {
		publicKey = &privateKey.(*rsa.PrivateKey).PublicKey
		_ = pem.Encode(caPrivateKeyPEM, &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey.(*rsa.PrivateKey)),
		})
	}
	privateKeyBytes = caPrivateKeyPEM.Bytes()
	return
}
