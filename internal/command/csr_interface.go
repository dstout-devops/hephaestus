// CSRGenerator defines an interface for generating CSRs.
package command

import (
	"crypto"

	"github.com/dstout-devops/hephaestus/internal/config"
)

type CSRGenerator interface {
	GenerateCSR(privKey crypto.PrivateKey, csrConfig config.CSRConfig) ([]byte, error)
}

type DefaultCSRGenerator struct{}

func (d *DefaultCSRGenerator) GenerateCSR(privKey crypto.PrivateKey, csrConfig config.CSRConfig) ([]byte, error) {
	return []byte{}, nil
}
