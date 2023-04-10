package doubleclick

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewPricersPool(t *testing.T) {
	pp := NewPricersPool(encryptionKeyRaw, integrityKeyRaw, 1000000)
	pricer := pp.AcquirePricer()
	defer pp.ReleasePricer(pricer)
	encryptedPrice := "anCGGFJApcfB6ZGc6mindhpTrYXHY4ONo7lXpg"
	encryptedPriceBytes := []byte(encryptedPrice) // don't inline
	priceInMicros, err := pricer.DecryptRaw(encryptedPriceBytes)
	assert.Equal(t, uint64(1354000), priceInMicros)
	assert.Equal(t, nil, err)
}
