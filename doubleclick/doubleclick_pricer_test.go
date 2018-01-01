package doubleclick

import (
	"testing"

	"github.com/benjaminch/openrtb-pricers/helpers"
	testshelpers "github.com/benjaminch/openrtb-pricers/tests_helpers"
)

func BuildNewPricer(encryptionKey string, integrityKey string, isBase64Keys bool, keyDecodingMode helpers.KeyDecodingMode, scaleFactor float64, isDebugMode bool) (*Pricer, error) {
	return NewPricer(encryptionKey, integrityKey, isBase64Keys, keyDecodingMode, scaleFactor, isDebugMode)
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
	var pricer *Pricer
	var err error
	pricer, err = BuildNewPricer(
		"652f83ada0545157a1b7fb0c0e09f59e7337332fe7abd4eb10449b8ee6c39135",
		"bd0a3dfb82ad95c5e63e159a62f73c6aca98ba2495322194759d512d77eb2bb5",
		false, // Keys are not base64
		helpers.Hexa,
		1000000,
		false,
	)

	if err != nil {
		t.Error("Error creating new Pricer : ", err)
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
	// Create a pricer with:
	// - UTF-8 keys
	// - Price scale factor as micro
	// - No debug mode
	var pricer *Pricer
	var err error
	pricer, err = BuildNewPricer(
		"6356770B3C111C07F778AFD69F16643E9110090FD4C479D91181EED2523788F1",
		"3588BF6D387E8AEAD4EEC66798255369AF47BFD48B056E8934CEFEF3609C469E",
		false, // Keys are not base64
		helpers.Utf8,
		1000000,
		false,
	)

	if err != nil {
		t.Error("Error creating new Pricer : ", err)
	}

	// Encrypted prices we will try to decrypt
	var pricesTestCase = []PriceTestCase{
		NewPriceTestCase("u7iq5XwQTNpAyThDrV5tuJXw-Y_IXQgkMA3RFA", 1.465),
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

func TestDecryptWithDebug(t *testing.T) {
	// TODO : To be implemented
}

func TestEncryptWithHexaKeys(t *testing.T) {
	// Create a pricer with:
	// - HEX keys
	// - Price scale factor as micro
	// - No debug mode
	var pricer *Pricer
	var err error
	pricer, err = BuildNewPricer(
		"652f83ada0545157a1b7fb0c0e09f59e7337332fe7abd4eb10449b8ee6c39135",
		"bd0a3dfb82ad95c5e63e159a62f73c6aca98ba2495322194759d512d77eb2bb5",
		false, // Keys are not base64
		helpers.Hexa,
		1000000,
		false,
	)

	if err != nil {
		t.Error("Error creating new Pricer : ", err)
	}

	// Clear prices we will try to encrypt
	var pricesTestCase = []PriceTestCase{
		NewPriceTestCase("1B2M2Y8AsgTpgAmY7PhCfgDo9mJGavHOuu-2SA==", 1.354),
		NewPriceTestCase("1B2M2Y8AsgTpgAmY7PhCfgDo9mJGTyiewYLbwg==", 3.24),
		NewPriceTestCase("1B2M2Y8AsgTpgAmY7PhCfgDo9mJGcRqedwjz2g==", 1),
		NewPriceTestCase("1B2M2Y8AsgTpgAmY7PhCfgDo9mJGc8xOTOXIGA==", 0.89),
		NewPriceTestCase("1B2M2Y8AsgTpgAmY7PhCfgDo9mJDi7nevR9kUw==", 100),
		NewPriceTestCase("1B2M2Y8AsgTpgAmY7PhCfgDo9mJGfn_O2Zdh_g==", 0.01),
	}

	for _, price := range pricesTestCase {
		var result string
		var err error
		result, err = pricer.Encrypt("", price.clear, false)
		if err != nil {
			t.Errorf("Encryption failed. Error : %s", err)
		}
		if result != price.encrypted {
			t.Errorf("Encryption failed. Should be : %s but was : %s", price.encrypted, result)
		}
	}
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

func TestEncryptDecryptWithHexaKeys(t *testing.T) {
	// Create a pricer with:
	// - HEX keys
	// - Price scale factor as micro
	// - No debug mode
	var pricer *Pricer
	var err error
	pricer, err = BuildNewPricer(
		"652f83ada0545157a1b7fb0c0e09f59e7337332fe7abd4eb10449b8ee6c39135",
		"bd0a3dfb82ad95c5e63e159a62f73c6aca98ba2495322194759d512d77eb2bb5",
		false, // Keys are not base64
		helpers.Hexa,
		1000000,
		false,
	)

	if err != nil {
		t.Error("Error creating new Pricer : ", err)
	}

	// Clear prices to encrypt
	var pricesTestCase = []PriceTestCase{
		NewPriceTestCase("", 1.465),
		NewPriceTestCase("", 0),
		NewPriceTestCase("", 100),
		NewPriceTestCase("", 1.45676),
		NewPriceTestCase("", 1.0),
		NewPriceTestCase("", 1000),
	}

	for _, price := range pricesTestCase {
		var decrypted float64
		var encrypted string
		var err error

		// Encrypt
		encrypted, err = pricer.Encrypt("", price.clear, false)
		if err != nil {
			t.Errorf("Encryption failed. Error : %s", err)
		}

		// Decrypt
		decrypted, err = pricer.Decrypt(encrypted, false)
		if err != nil {
			t.Errorf("Decryption failed. Error : %s", err)
		}

		// Assert that the decrypted price is the one with encrypted in a first place
		if !testshelpers.FloatEquals(decrypted, price.clear) {
			t.Errorf("Decryption failed. Should be : %f but was : %f", price.clear, decrypted)
		}
	}
}

func TestEncryptDecryptWithUtf8Keys(t *testing.T) {
	// Create a pricer with:
	// - UTF-8 keys
	// - Price scale factor as micro
	// - No debug mode
	var pricer *Pricer
	var err error
	pricer, err = BuildNewPricer(
		"6356770B3C111C07F778AFD69F16643E9110090FD4C479D91181EED2523788F1",
		"3588BF6D387E8AEAD4EEC66798255369AF47BFD48B056E8934CEFEF3609C469E",
		false, // Keys are not base64
		helpers.Utf8,
		1000000,
		false,
	)

	if err != nil {
		t.Error("Error creating new Pricer : ", err)
	}

	// Clear prices to encrypt
	var pricesTestCase = []PriceTestCase{
		NewPriceTestCase("", 1.465),
		NewPriceTestCase("", 0),
		NewPriceTestCase("", 100),
		NewPriceTestCase("", 1.45676),
		NewPriceTestCase("", 1.0),
		NewPriceTestCase("", 1000),
	}

	for _, price := range pricesTestCase {
		var decrypted float64
		var encrypted string
		var err error

		// Encrypt
		encrypted, err = pricer.Encrypt("", price.clear, false)
		if err != nil {
			t.Errorf("Encryption failed. Error : %s", err)
		}

		// Decrypt
		decrypted, err = pricer.Decrypt(encrypted, false)
		if err != nil {
			t.Errorf("Decryption failed. Error : %s", err)
		}

		// Assert that the decrypted price is the one with encrypted in a first place
		if !testshelpers.FloatEquals(decrypted, price.clear) {
			t.Errorf("Decryption failed. Should be : %f but was : %f", price.clear, decrypted)
		}
	}
}
