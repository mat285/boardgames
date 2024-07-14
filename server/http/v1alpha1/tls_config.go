package v1alpha1

import (
	"crypto/tls"
)

type TLS struct {
	CertFile string `json:"certFile" yaml:"certFile"`
	KeyFile  string `json:"keyFile" yaml:"keyFile"`
}

func (c TLS) TLSConfig() *tls.Config {
	if len(c.CertFile) == 0 || len(c.KeyFile) == 0 {
		return nil
	}
	cert, err := tls.LoadX509KeyPair(c.CertFile, c.KeyFile)
	if err != nil {
		return nil
	}

	return &tls.Config{
		Certificates: []tls.Certificate{
			cert,
		},
	}
}
