package command

import (
	"crypto/ed25519"
	"crypto/rsa"

	"github.com/dstout-devops/hephaestus/internal/keys"
)

type KeyGenerator interface {
	GenerateEd25519Key() (ed25519.PrivateKey, error)
	GenerateRSAKey(bits int) (*rsa.PrivateKey, error)
}

type DefaultKeyGenerator struct{}

func (d *DefaultKeyGenerator) GenerateEd25519Key() (ed25519.PrivateKey, error) {
	return keys.GenerateEd25519Key()
}

func (d *DefaultKeyGenerator) GenerateRSAKey(bits int) (*rsa.PrivateKey, error) {
	return keys.GenerateRSAKey(bits)
}
