package doubleclick

import (
	"testing"

	"github.com/benjaminch/openrtb-pricers/helpers"
	testshelpers "github.com/benjaminch/openrtb-pricers/tests_helpers"
)

func BuildNewDoubleClickPricer(encryptionKey string, integrityKey string, isBase64Keys bool, keyDecodingMode helpers.KeyDecodingMode, scaleFactor float64, isDebugMode bool) (*DoubleClickPricer, error) {
	return NewDoubleClickPricer(encryptionKey, integrityKey, isBase64Keys, keyDecodingMode, scaleFactor, isDebugMode)
}

type PriceTestCase struct {
	encrypted   string
	clear       float64
	scaleFactor float64
}

func NewPriceTestCase(encrypted string, clear float64, scaleFactor float64) PriceTestCase {
	return PriceTestCase{encrypted: encrypted, clear: clear, scaleFactor: scaleFactor}
}

func TestDecryptGoogleOfficialExamples(t *testing.T) {
	// From specs examples
	// https://developers.google.com/ad-exchange/rtb/response-guide/decrypt-price

	var pricer *DoubleClickPricer
	var err error
	pricer, err = BuildNewDoubleClickPricer(
		"ZS-DraBUUVeht_sMDgn1nnM3My_nq9TrEESbjubDkTU",
		"vQo9-4KtlcXmPhWaYvc8asqYuiSVMiGUdZ1RLXfrK7U",
		true, // Keys are base64
		helpers.Utf8,
		1000000,
		false,
	)

	if err != nil {
		t.Error("Error creating new Pricer : ", err)
	}

	// Encrypted prices we will try to decrypt
	var pricesTestCase = []PriceTestCase{
		NewPriceTestCase("anCGGFJApcfB6ZGc6mindhpTrYXHY4ONo7lXpg", 1.354, 1000000),
		NewPriceTestCase("ce131TRp7waIZI2qOiRr2DMm2sSIeGh_wIAwVQ", 3.24, 1000000),
		NewPriceTestCase("K6tfPnPvN_5E2xS3GssrFYeouJJRkBQqxR_FxQ", 1, 1000000),
		NewPriceTestCase("lEzCWnwgB21Dy2_H43PKZeZaNDstZZElZRFTDQ", 0.89, 1000000),
		NewPriceTestCase("L91lB6giyIXh2o4CeUf0F7sCXozKWRXAUeMUfg", 100, 1000000),
		NewPriceTestCase("8WY0BgWbds1eEVNFkrXVIr1GU08iueKrP0wXfw", 0.01, 1000000),
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
		NewPriceTestCase("anCGGFJApcfB6ZGc6mindhpTrYXHY4ONo7lXpg", 1.354, 1000000),
		NewPriceTestCase("ce131TRp7waIZI2qOiRr2DMm2sSIeGh_wIAwVQ", 3.24, 1000000),
		NewPriceTestCase("K6tfPnPvN_5E2xS3GssrFYeouJJRkBQqxR_FxQ", 1, 1000000),
		NewPriceTestCase("lEzCWnwgB21Dy2_H43PKZeZaNDstZZElZRFTDQ", 0.89, 1000000),
		NewPriceTestCase("L91lB6giyIXh2o4CeUf0F7sCXozKWRXAUeMUfg", 100, 1000000),
		NewPriceTestCase("8WY0BgWbds1eEVNFkrXVIr1GU08iueKrP0wXfw", 0.01, 1000000),
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
	var pricer *DoubleClickPricer
	var err error
	pricer, err = BuildNewDoubleClickPricer(
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
		NewPriceTestCase("u7iq5XwQTNpAyThDrV5tuJXw-Y_IXQgkMA3RFA", 1.465, 1000000),
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

func TestDecryptWithScaleFactor(t *testing.T) {
	// Encrypted prices we will try to decrypt
	// with several scale factors
	var pricesTestCase = []PriceTestCase{
		NewPriceTestCase("1B2M2Y8AsgTpgAmY7PhCfgDo9mJGVwr-Q_z9Cw==", 1.354, 2000000),
		NewPriceTestCase("1B2M2Y8AsgTpgAmY7PhCfgDo9mJGNHC-ozlJ9Q==", 3.24, 1500000),
		NewPriceTestCase("1B2M2Y8AsgTpgAmY7PhCfgDo9mJGf95-RX4DPw==", 1, 100000),
		NewPriceTestCase("1B2M2Y8AsgTpgAmY7PhCfgDo9mJGc8xOTOXIGA==", 0.89, 1000000),
		NewPriceTestCase("1B2M2Y8AsgTpgAmY7PhCfgDo9mJEhKheuuMVqg==", 100, 500000),
		NewPriceTestCase("1B2M2Y8AsgTpgAmY7PhCfgDo9mJGfn-cvGvYQg==", 0.01, 1005000),
	}

	for _, priceTestCase := range pricesTestCase {
		// Create a pricer with:
		// - HEX keys
		// - Price scale factor as micro
		// - No debug mode
		var pricer *DoubleClickPricer
		var err error
		pricer, err = BuildNewDoubleClickPricer(
			"652f83ada0545157a1b7fb0c0e09f59e7337332fe7abd4eb10449b8ee6c39135",
			"bd0a3dfb82ad95c5e63e159a62f73c6aca98ba2495322194759d512d77eb2bb5",
			false, // Keys are not base64
			helpers.Hexa,
			priceTestCase.scaleFactor,
			false,
		)

		if err != nil {
			t.Error("Error creating new Pricer : ", err)
		}

		var result float64
		result, err = pricer.Decrypt(priceTestCase.encrypted, false)
		if err != nil {
			t.Errorf("Decryption failed. Error : %s", err)
		}
		if result != priceTestCase.clear {
			t.Errorf("Decryption failed. Should be : %f but was : %f (scale factor: %f)", priceTestCase.clear, result, priceTestCase.scaleFactor)
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
	var pricer *DoubleClickPricer
	var err error
	pricer, err = BuildNewDoubleClickPricer(
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
		NewPriceTestCase("1B2M2Y8AsgTpgAmY7PhCfgDo9mJGavHOuu-2SA==", 1.354, 1000000),
		NewPriceTestCase("1B2M2Y8AsgTpgAmY7PhCfgDo9mJGTyiewYLbwg==", 3.24, 1000000),
		NewPriceTestCase("1B2M2Y8AsgTpgAmY7PhCfgDo9mJGcRqedwjz2g==", 1, 1000000),
		NewPriceTestCase("1B2M2Y8AsgTpgAmY7PhCfgDo9mJGc8xOTOXIGA==", 0.89, 1000000),
		NewPriceTestCase("1B2M2Y8AsgTpgAmY7PhCfgDo9mJDi7nevR9kUw==", 100, 1000000),
		NewPriceTestCase("1B2M2Y8AsgTpgAmY7PhCfgDo9mJGfn_O2Zdh_g==", 0.01, 1000000),
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

/*
func TestEncryptWithUtf8Keys(t *testing.T) {
	// Create a pricer with:
	// - UTF-8 keys
	// - Price scale factor as micro
	// - No debug mode
	var pricer *DoubleClickPricer
	var err error
	pricer, err = BuildNewDoubleClickPricer(
		"ZS_Cg8KtwqBUUVfCocK3w7sMDgnDtcKeczczL8OnwqvDlMOrEETCm8KOw6bDg8KRNQ",
		"wr0KPcO7woLCrcKVw4XDpj4Vwppiw7c8asOKwpjCuiTClTIhwpR1wp1RLXfDqyvCtQ",
		false, // Keys are base64
		helpers.Utf8,
		1000000,
		false,
	)

	if err != nil {
		t.Error("Error creating new Pricer : ", err)
	}

	// Clear prices we will try to encrypt
	var pricesTestCase = []PriceTestCase{
		NewPriceTestCase("1B2M2Y8AsgTpgAmY7PhCfgDo9mJGavHOuu-2SA==", 1.354, 1000000),
		NewPriceTestCase("1B2M2Y8AsgTpgAmY7PhCfgDo9mJGTyiewYLbwg==", 3.24, 1000000),
		NewPriceTestCase("1B2M2Y8AsgTpgAmY7PhCfgDo9mJGcRqedwjz2g==", 1, 1000000),
		NewPriceTestCase("1B2M2Y8AsgTpgAmY7PhCfgDo9mJGc8xOTOXIGA==", 0.89, 1000000),
		NewPriceTestCase("1B2M2Y8AsgTpgAmY7PhCfgDo9mJDi7nevR9kUw==", 100, 1000000),
		NewPriceTestCase("1B2M2Y8AsgTpgAmY7PhCfgDo9mJGfn_O2Zdh_g==", 0.01, 1000000),
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
*/

func TestEncryptWithScaleFactor(t *testing.T) {

	// Clear prices we will try to encrypt
	// with several scale factors
	var pricesTestCase = []PriceTestCase{
		NewPriceTestCase("1B2M2Y8AsgTpgAmY7PhCfgDo9mJGVwr-Q_z9Cw==", 1.354, 2000000),
		NewPriceTestCase("1B2M2Y8AsgTpgAmY7PhCfgDo9mJGNHC-ozlJ9Q==", 3.24, 1500000),
		NewPriceTestCase("1B2M2Y8AsgTpgAmY7PhCfgDo9mJGf95-RX4DPw==", 1, 100000),
		NewPriceTestCase("1B2M2Y8AsgTpgAmY7PhCfgDo9mJGc8xOTOXIGA==", 0.89, 1000000),
		NewPriceTestCase("1B2M2Y8AsgTpgAmY7PhCfgDo9mJEhKheuuMVqg==", 100, 500000),
		NewPriceTestCase("1B2M2Y8AsgTpgAmY7PhCfgDo9mJGfn-cvGvYQg==", 0.01, 1005000),
	}

	for _, priceTestCase := range pricesTestCase {
		// Create a pricer with:
		// - HEX keys
		// - Price scale factor as micro
		// - No debug mode
		var pricer *DoubleClickPricer
		var err error
		pricer, err = BuildNewDoubleClickPricer(
			"652f83ada0545157a1b7fb0c0e09f59e7337332fe7abd4eb10449b8ee6c39135",
			"bd0a3dfb82ad95c5e63e159a62f73c6aca98ba2495322194759d512d77eb2bb5",
			false, // Keys are not base64
			helpers.Hexa,
			priceTestCase.scaleFactor,
			false,
		)

		if err != nil {
			t.Error("Error creating new Pricer : ", err)
		}

		var result string
		result, err = pricer.Encrypt("", priceTestCase.clear, false)
		t.Logf("%f : %s", priceTestCase.clear, result)
		if err != nil {
			t.Errorf("Encryption failed. Error : %s", err)
		}
		if result != priceTestCase.encrypted {
			t.Errorf("Encryption failed. Should be : %s but was : %s (scale factor: %f)", priceTestCase.encrypted, result, priceTestCase.scaleFactor)
		}
	}
}

func TestEncryptWithDebug(t *testing.T) {
	// TODO : To be implemented
}

func TestEncryptDecryptWithHexaKeys(t *testing.T) {
	// Create a pricer with:
	// - HEX keys
	// - Price scale factor as micro
	// - No debug mode
	var pricer *DoubleClickPricer
	var err error
	pricer, err = BuildNewDoubleClickPricer(
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
		NewPriceTestCase("", 1.465, 1000000),
		NewPriceTestCase("", 0, 1000000),
		NewPriceTestCase("", 100, 1000000),
		NewPriceTestCase("", 1.45676, 1000000),
		NewPriceTestCase("", 1.0, 1000000),
		NewPriceTestCase("", 1000, 1000000),
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
	var pricer *DoubleClickPricer
	var err error
	pricer, err = BuildNewDoubleClickPricer(
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
		NewPriceTestCase("", 1.465, 1000000),
		NewPriceTestCase("", 0, 1000000),
		NewPriceTestCase("", 100, 1000000),
		NewPriceTestCase("", 1.45676, 1000000),
		NewPriceTestCase("", 1.0, 1000000),
		NewPriceTestCase("", 1000, 1000000),
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

func TestEncryptDecryptWithSeed(t *testing.T) {
	// Seeds to test
	var seedsToTest = []string{"", "test", "a", "b", "azertyuiopmlkjhgfdsqwxcvbn"}

	// Prices to be encrypted / decrypted
	var pricesTestCase = []PriceTestCase{
		NewPriceTestCase("", 1.465, 1000000),
		NewPriceTestCase("", 0, 1000000),
		NewPriceTestCase("", 100, 1000000),
		NewPriceTestCase("", 1.45676, 1000000),
		NewPriceTestCase("", 1.0, 1000000),
		NewPriceTestCase("", 1000, 1000000),
	}

	// Create a pricer with:
	// - HEX keys
	// - Price scale factor as micro
	// - No debug mode
	var pricer *DoubleClickPricer
	var err error
	pricer, err = BuildNewDoubleClickPricer(
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

	var encryptedPrices []string

	for _, seed := range seedsToTest {

		for _, price := range pricesTestCase {
			var decrypted float64
			var encrypted string
			var err error

			// Encrypt
			encrypted, err = pricer.Encrypt(seed, price.clear, false)
			if err != nil {
				t.Errorf("Encryption failed. Error : %s", err)
			}
			// Store encrypted prices for further tests
			encryptedPrices = append(encryptedPrices, encrypted)

			// Decrypt
			decrypted, err = pricer.Decrypt(encrypted, false)
			if err != nil {
				t.Errorf("Decryption failed. Error : %s", err)
			}

			// Assert that the decrypted price is the one with encrypted in a first place
			if !testshelpers.FloatEquals(decrypted, price.clear) {
				t.Errorf("Decryption failed. Should be : %f but was : %f (seed : %s)", price.clear, decrypted, seed)
			}
		}
	}

	// Checking that every single encrypted prices are different to each others
	// checking only the same prices with all seeds
	initialPricesCount := len(pricesTestCase)
	for i, encryptedPrice := range encryptedPrices {
		if i < initialPricesCount-1 && encryptedPrice == encryptedPrices[i+initialPricesCount] {
			t.Errorf("Encrypted prices are the same but they shouldn't : price : %f, encrypted : %s, seed : %s", pricesTestCase[i].clear, encryptedPrice, seedsToTest[i%initialPricesCount])
		}
	}
}

func TestEncryptDecryptWithScaleFactor(t *testing.T) {
	// Scale factors to test
	var scaleFactorsToTest = []float64{0.1, 1, 10, 50, 100, 10000}

	// Prices we will try to encrypt / decrypt
	// Because by design, the price is stored as int64 before encryption
	// we need to encrypt big prices in order to be able to test a full
	// range of scale factors (from 0.1 to 10,000).
	// Not doing so will end up having rounding issue.
	var pricesTestCase = []PriceTestCase{
		NewPriceTestCase("", 13540, 1000000),
		NewPriceTestCase("", 3240, 1000000),
		NewPriceTestCase("", 10, 1000000),
		NewPriceTestCase("", 890, 1000000),
		NewPriceTestCase("", 1000, 1000000),
	}

	for _, scaleFactor := range scaleFactorsToTest {

		// Create a pricer with:
		// - HEX keys
		// - Price scale factor as micro
		// - No debug mode
		var pricer *DoubleClickPricer
		var err error
		pricer, err = BuildNewDoubleClickPricer(
			"652f83ada0545157a1b7fb0c0e09f59e7337332fe7abd4eb10449b8ee6c39135",
			"bd0a3dfb82ad95c5e63e159a62f73c6aca98ba2495322194759d512d77eb2bb5",
			false, // Keys are not base64
			helpers.Hexa,
			scaleFactor,
			false,
		)

		if err != nil {
			t.Error("Error creating new Pricer : ", err)
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
				t.Errorf("Decryption failed. Should be : %f but was : %f (scale factor : %f)", price.clear, decrypted, scaleFactor)
			}
		}
	}
}
