package doubleclick

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"hash"
	"strings"

	"github.com/benjaminch/pricers/helpers"
)

var ErrWrongSize = errors.New("Encrypted price is not 38 chars")
var ErrWrongSignature = errors.New("Failed to decrypt")

// DoubleClickPricer implementing price encryption and decryption
// Specs : https://developers.google.com/ad-exchange/rtb/response-guide/decrypt-price
type DoubleClickPricer struct {
	encryptionKey hash.Hash
	integrityKey  hash.Hash
	scaleFactor   float64
}

// NewDoubleClickPricer returns a DoubleClickPricer struct.
// Keys are either Base64Url (websafe) of hexa. keyDecodingMode
// should be used to specify how keys should be decoded.
// Factor the clear price will be multiplied by before encryption.
// from specs, scaleFactor is 1,000,000, but you can set something else.
// Be aware that the price is stored as an int64 so depending on the digits
// precision you want, picking a scale factor smaller than 1,000,000 may lead
// to price to be rounded and loose some digits precision.
func NewDoubleClickPricer(
	encryptionKey string,
	integrityKey string,
	isBase64Keys bool,
	keyDecodingMode helpers.KeyDecodingMode,
	scaleFactor float64) (*DoubleClickPricer, error) {
	var err error

	encryptionKeyRaw, err := helpers.RawKeyBytes(encryptionKey, isBase64Keys, keyDecodingMode)
	if err != nil {
		return nil, err
	}
	integrityKeyRaw, err := helpers.RawKeyBytes(integrityKey, isBase64Keys, keyDecodingMode)
	if err != nil {
		return nil, err
	}

	encryptingFun := helpers.CreateHmac(encryptionKeyRaw)
	integrityFun := helpers.CreateHmac(integrityKeyRaw)

	return &DoubleClickPricer{
			encryptionKey: encryptingFun,
			integrityKey:  integrityFun,
			scaleFactor:   scaleFactor},
		nil
}

func NewDoubleClickPricerFromRawKeys(
	encryptionKeyRaw []byte,
	integrityKeyRaw []byte,
	scaleFactor float64) *DoubleClickPricer {
	encryptingFun := helpers.CreateHmac(encryptionKeyRaw)
	integrityFun := helpers.CreateHmac(integrityKeyRaw)
	return &DoubleClickPricer{
		encryptionKey: encryptingFun,
		integrityKey:  integrityFun,
		scaleFactor:   scaleFactor}
}

// Encrypt encrypts a clear price and a given seed.
func (dc *DoubleClickPricer) Encrypt(seed string, price float64) string {
	data := helpers.ApplyScaleFactor(price, dc.scaleFactor)

	// Create Initialization Vector from seed
	iv := md5.Sum([]byte(seed))

	//pad = hmac(e_key, iv), first 8 bytes
	pad := helpers.HmacSum(dc.encryptionKey, iv[:], nil)[:8]

	// signature = hmac(i_key, data || iv), first 4 bytes
	signature := helpers.HmacSum(dc.integrityKey, data[:], iv[:])[:4]

	// enc_data = pad <xor> data
	encoded := [8]byte{}
	for i := range data {
		encoded[i] = pad[i] ^ data[i]
	}

	// final_message = WebSafeBase64Encode( iv || enc_price || signature )
	return base64.RawURLEncoding.EncodeToString(append(append(iv[:], encoded[:]...), signature...))
}

// Decrypt decrypts an encrypted price.
func (dc *DoubleClickPricer) Decrypt(encryptedPrice string) (float64, error) {

	// Decode base64 url
	// Just to be safe remove padding if it was added by mistake
	encryptedPrice = strings.TrimRight(encryptedPrice, "=")
	if len(encryptedPrice) != 38 {
		return 0, ErrWrongSize
	}
	decoded, err := base64.RawURLEncoding.DecodeString(encryptedPrice)
	if err != nil {
		return 0, err
	}

	// Get elements
	iv := decoded[0:16]
	p := decoded[16:24]
	signature := decoded[24:28]

	// pad = hmac(e_key, iv)
	pad := helpers.HmacSum(dc.encryptionKey, iv, nil)[:8]

	// priceMicro = p <xor> pad
	priceMicro := [8]byte{}
	for i := range p {
		priceMicro[i] = pad[i] ^ p[i]
	}

	// conf_sig = hmac(i_key, data || iv)
	confirmationSignature := helpers.HmacSum(dc.integrityKey, priceMicro[:], iv)[:4]

	// success = (conf_sig == sig)
	if !bytes.Equal(confirmationSignature, signature) {
		return 0, ErrWrongSignature
	}
	priceInMicros := binary.BigEndian.Uint64(priceMicro[:])
	price := float64(priceInMicros) / dc.scaleFactor
	return price, nil
}
