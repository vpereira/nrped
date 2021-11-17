package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"io/ioutil"
	"log"
	"math/big"
	"time"
)

// creating self-signed certs
// TODO, add the IP address in the certificate
func createCertificate(country string, organization string,
	organizationalUnit string, pemFile string, keyFile string) bool {

	ca := &x509.Certificate{
		SerialNumber: big.NewInt(7829),
		Subject: pkix.Name{
			Country:            []string{country},
			Organization:       []string{organization},
			OrganizationalUnit: []string{organizationalUnit},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0), // 10 years
		SubjectKeyId:          []byte{1, 2, 3, 4, 5},
		BasicConstraintsValid: true,
		IsCA:                  true,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth,
			x509.ExtKeyUsageServerAuth},
		KeyUsage: x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}

	priv, _ := rsa.GenerateKey(rand.Reader, 2048)
	pub := &priv.PublicKey
	ca_b, err := x509.CreateCertificate(rand.Reader, ca, ca, pub, priv)

	if err != nil {
		log.Println("create ca failed", err)
		return false
	}

	log.Println("write to", pemFile)
	ioutil.WriteFile(pemFile, ca_b, 0777)

	priv_b := x509.MarshalPKCS1PrivateKey(priv)
	log.Println("write to", keyFile)
	ioutil.WriteFile(keyFile, priv_b, 0777)

	return true
}

func main() {
	ret := createCertificate("Germany", "NRPED", "Server", "server.pem", "server.key")

	if ret != true {
		log.Println("create ca failed")
		return
	}

	ret = createCertificate("Germany", "NRPED", "Client", "client.pem", "client.key")

	if ret != true {
		log.Println("create ca failed")
		return
	}

}
