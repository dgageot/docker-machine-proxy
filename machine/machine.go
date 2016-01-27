package machine

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
)

type DockerMachine struct {
	Url      string
	CertPath string
}

func (dm *DockerMachine) ReadTLSConfig() (*tls.Config, error) {
	caCert, err := ioutil.ReadFile(dm.CertPath + "/ca.pem")
	if err != nil {
		return nil, err
	}

	serverCert, err := ioutil.ReadFile(dm.CertPath + "/server.pem")
	if err != nil {
		return nil, err
	}

	serverKey, err := ioutil.ReadFile(dm.CertPath + "/server-key.pem")
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()

	ok := certPool.AppendCertsFromPEM(caCert)
	if !ok {
		return nil, errors.New("There was an error reading certificate")
	}

	keypair, err := tls.X509KeyPair(serverCert, serverKey)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		RootCAs:      certPool,
		Certificates: []tls.Certificate{keypair},
	}, nil
}
