package doubleclick

import (
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
		"090a654e859cd11d673ad1f21f3ae57447fb8037f3cb1adb05b0897b3c496992",
		"c47858f4f24e8111272f2042b467477b9b0040da232336792e57aa267b500114",
		helpers.Utf8,
		1000000,
		false,
	)

	if err != nil {
		t.Error("Error creating new DoubleClickPricer")
	}

	// Encrypted prices we will try to decrypt
	var pricesTestCase = []*PriceTestCase{
		NewPriceTestCase("9Xw3oTa2pnw3mEvdA60wnPqaBaPmm8KaJAsZIg", 100.386),
		NewPriceTestCase("9nw3oTa2pnw3mEvdA60wnGYgYK_ghaiOJsuaBQ", 0.003),
		NewPriceTestCase("93w3oTa2pnw3mEvdA60wnKJPC6nNSQFkeNyrOQ", 0.401),
	}

	for _, encryptedPrice := range pricesTestCase {
		var result float64
		var err error
		result, err = pricer.Decrypt(encryptedPrice.encrypted)
		if err != nil {
			t.Errorf("Decryption failed. Error : %s", err)
		}
		if result == encryptedPrice.clear {
			t.Errorf("Decryption failed. Should be : %f but was : %f", encryptedPrice.clear, result)
		}
	}
}
