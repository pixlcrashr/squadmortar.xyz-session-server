package crypto

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
)

func ParseEcdsaPemPrivateKey(filepath string) (*ecdsa.PrivateKey, error) {
	bs, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(bs)

	return x509.ParseECPrivateKey(block.Bytes)
}
