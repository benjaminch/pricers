package doubleclick

import (
	"fmt"
	"testing"

	"../helpers"
)

func BuildNewDoubleClickPricer(encryptionKey, integrityKey string, keyDecodingMode helpers.KeyDecodingMode, scaleFactor float64, isDebugMode bool) (*DoubleClickPricer, error) {
	return NewDoubleClickPricer(encryptionKey, integrityKey, keyDecodingMode, scaleFactor, isDebugMode)
}

type PriceTestCase struct {
	encrypted string
	clear     float64
}

func NewPriceTestCase(encrypted string, clear float64) *PriceTestCase {
	return &PriceTestCase{encrypted: encrypted, clear: clear}
}

func TestDecrypt(t *testing.T) {

	// Create a pricer with:
	// - HEX keys
	// - Price scale factor as micro
	// - No debug mode
	var pricer *DoubleClickPricer
	var err error
	pricer, err = BuildNewDoubleClickPricer(
		"skU7Ax_NL5pPAFyKdkfZjZz2-VhIN8bjj1rVFOaJ_5o=",
		"arO23ykdNqUQ5LEoQ0FVmPkBd7xB5CO89PDZlSjpFxo=",
		helpers.Utf8,
		1000000,
		false,
	)

	if err != nil {
		t.Error("Error creating new DoubleClickPricer")
	}

	// Encrypted prices we will try to decrypt
	var pricesTestCase = []*PriceTestCase{
		NewPriceTestCase("WEp8wQAAAABnFd5EkB2k1wJeFcAj-Z_JVOeGzA==", 100.0),
		NewPriceTestCase("WEp8sQAAAACwF6CtLJrXSRFBM8UiTTIyngN-og==", 1900.0),
		NewPriceTestCase("WEp8nQAAAAADG-y45xxIC1tMWuTjzmDW6HtroQ==", 2700.0),
	}

	for _, encryptedPrice := range pricesTestCase {
		var result = pricer.Decrypt(encryptedPrice.encrypted)
		fmt.Println(result)
		if result == encryptedPrice.clear {
			t.Errorf("Decryption failed. Should be : %f but was : %f", encryptedPrice.clear, result)
		}
	}
}
