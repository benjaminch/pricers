package doubleclick

import (
	"testing"

	"github.com/benjaminch/openrtb-pricers/helpers"
	testshelpers "github.com/benjaminch/openrtb-pricers/tests_helpers"
)

func BuildNewDoubleClickPricer(encryptionKey string, integrityKey string, keyDecodingMode helpers.KeyDecodingMode, scaleFactor float64, isDebugMode bool) (*DoubleClickPricer, error) {
	return NewDoubleClickPricer(encryptionKey, integrityKey, keyDecodingMode, scaleFactor, isDebugMode)
}

type PriceTestCase struct {
	encrypted string
	clear     float64
}

func NewPriceTestCase(encrypted string, clear float64) PriceTestCase {
	return PriceTestCase{encrypted: encrypted, clear: clear}
}

func TestDecryptWithHexaKeys(t *testing.T) {

	// Create a pricer with:
	// - HEX keys
	// - Price scale factor as micro
	// - No debug mode
	var pricer *DoubleClickPricer
	var err error
	pricer, err = BuildNewDoubleClickPricer(
		"652f83ada0545157a1b7fb0c0e09f59e7337332fe7abd4eb10449b8ee6c39135",
		"bd0a3dfb82ad95c5e63e159a62f73c6aca98ba2495322194759d512d77eb2bb5",
		helpers.Hexa,
		1000000,
		false,
	)

	if err != nil {
		t.Error("Error creating new DoubleClickPricer")
	}

	// Encrypted prices we will try to decrypt
	var pricesTestCase = []PriceTestCase{
		NewPriceTestCase("anCGGFJApcfB6ZGc6mindhpTrYXHY4ONo7lXpg", 1.354),
		NewPriceTestCase("ce131TRp7waIZI2qOiRr2DMm2sSIeGh_wIAwVQ", 3.24),
		NewPriceTestCase("K6tfPnPvN_5E2xS3GssrFYeouJJRkBQqxR_FxQ", 1),
		NewPriceTestCase("lEzCWnwgB21Dy2_H43PKZeZaNDstZZElZRFTDQ", 0.89),
		NewPriceTestCase("L91lB6giyIXh2o4CeUf0F7sCXozKWRXAUeMUfg", 100),
		NewPriceTestCase("8WY0BgWbds1eEVNFkrXVIr1GU08iueKrP0wXfw", 0.01),
	}

	for _, encryptedPrice := range pricesTestCase {
		var result float64
		var err error
		result, err = pricer.Decrypt(encryptedPrice.encrypted, false)
		if err != nil {
			t.Errorf("Decryption failed. Error : %s", err)
		}
		if !testshelpers.FloatEquals(result, encryptedPrice.clear) {
			t.Errorf("Decryption failed. Should be : %f but was : %f", encryptedPrice.clear, result)
		}
	}
}

func TestDecryptWithUtf8Keys(t *testing.T) {
	// TODO : To be implemented
}

func TestDecryptWithDebug(t *testing.T) {
	// TODO : To be implemented
}

func TestEncryptWithHexaKeys(t *testing.T) {
	// TODO : To be implemented
}

func TestEncryptWithUtf8Keys(t *testing.T) {
	// TODO : To be implemented
}

func TestEncryptWithScaleFactor(t *testing.T) {
	// TODO : To be implemented
}

func TestEncryptWithDebug(t *testing.T) {
	// TODO : To be implemented
}
