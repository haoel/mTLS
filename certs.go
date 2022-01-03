package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/ioutil"
	"log"
	"math/big"
	"time"
)

const path string = "./certs/"

func makeCA(subject *pkix.Name) (*x509.Certificate, *rsa.PrivateKey, error) {
	// creating a CA which will be used to sign all of our certificates using the x509 package from the Go Standard Library
	caCert := &x509.Certificate{
		SerialNumber:          big.NewInt(2019),
		Subject:               *subject,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10*365, 0, 0),
		IsCA:                  true, // <- indicating this certificate is a CA certificate.
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}
	// generate a private key for the CA
	caKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Printf("Generate the CA Private Key error: %v\n", err)
		return nil, nil, err
	}

	// create the CA certificate
	caBytes, err := x509.CreateCertificate(rand.Reader, caCert, caCert, &caKey.PublicKey, caKey)
	if err != nil {
		log.Printf("Create the CA Certificate error: %v\n", err)
		return nil, nil, err
	}

	// Create the CA PEM files
	caPEM := new(bytes.Buffer)
	pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})

	if err := ioutil.WriteFile(path + "ca.crt", caPEM.Bytes(), 0644); err != nil {
		log.Printf("Write the CA certificate file error: %v\n", err)
		return nil, nil, err
	}

	caPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(caPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caKey),
	})
	if err := ioutil.WriteFile(path + "ca.key", caPEM.Bytes(), 0644); err != nil {
		log.Printf("Write the CA certificate file error: %v\n", err)
		return nil, nil, err
	}
	return caCert, caKey, nil
}

func makeCert(caCert *x509.Certificate, caKey *rsa.PrivateKey, subject *pkix.Name, name string) error {
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject:      *subject,
		//IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		DNSNames:     []string{"localhost"},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	certKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Printf("Generate the Key error: %v\n", err)
		return err
	}
	certBytes, err := x509.CreateCertificate(rand.Reader, cert, caCert, &certKey.PublicKey, caKey)
	if err != nil {
		log.Printf("Generate the certificate error: %v\n", err)
		return err
	}

	certPEM := new(bytes.Buffer)
	pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})
	if err := ioutil.WriteFile(path + name + ".crt", certPEM.Bytes(), 0644); err != nil {
		log.Printf("Write the CA certificate file error: %v\n", err)
		return err
	}

	certKeyPEM := new(bytes.Buffer)
	pem.Encode(certKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certKey),
	})
	if err := ioutil.WriteFile(path + name + ".key", certKeyPEM.Bytes(), 0644); err != nil {
		log.Printf("Write the CA certificate file error: %v\n", err)
		return err
	}
	return nil
}

func main() {
	subject := pkix.Name{
		Country:            []string{"Earth"},
		Organization:       []string{"CA Company"},
		OrganizationalUnit: []string{"Engineering"},
		Locality:           []string{"Mountain"},
		Province:           []string{"Asia"},
		StreetAddress:      []string{"Bridge"},
		PostalCode:         []string{"123456"},
		SerialNumber:       "",
		CommonName:         "CA",
		Names:              []pkix.AttributeTypeAndValue{},
		ExtraNames:         []pkix.AttributeTypeAndValue{},
	}
	caCert, caKey, err := makeCA(&subject)
	if err != nil {
		log.Fatalf("make CA Certificate error!")
	}
	log.Println("Create the CA certificate successfully.")

	subject.CommonName = "Server"
	subject.Organization = []string{"Server Company"}
	if err := makeCert(caCert, caKey, &subject, "server"); err != nil {
		log.Fatal("make Server Certificate error!")
	}
	log.Println("Create and Sign the Server certificate successfully.")

	subject.CommonName = "Client A"
	subject.Organization = []string{"Client A Company"}
	if err := makeCert(caCert, caKey, &subject, "client.a"); err != nil {
		log.Fatal("make Client A Certificate error!")
	}
	log.Println("Create and Sign the Client A certificate successfully.")

	subject.CommonName = "Client B"
	subject.Organization = []string{"Client B Company"}
	if err := makeCert(caCert, caKey, &subject, "client.b"); err != nil {
		log.Fatal("make Client B Certificate error!")
	}
	log.Println("Create and Sign the Client B certificate successfully.")

}
