package doubleclick

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"github.com/benjaminch/pricers/helpers"
	"hash"
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
	pad := helpers.HmacSum(dc.encryptionKey, iv[:], nil, nil)[:8]

	// signature = hmac(i_key, data || iv), first 4 bytes
	signature := helpers.HmacSum(dc.integrityKey, data[:], iv[:], nil)[:4]

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
	buf := make([]byte, 56)
	priceInMicros, err := dc.DecryptRaw([]byte(encryptedPrice), buf)
	price := float64(priceInMicros) / dc.scaleFactor
	return price, err
}

// DecryptRaw decrypts an encrypted price.
// It returns the price as integer in micros without applying a scaleFactor
// You must pass a buffer (56 bytes) for decoder so that can reused again to avoid allocation
func (dc *DoubleClickPricer) DecryptRaw(encryptedPrice []byte, buf []byte) (uint64, error) {
	// first 28 bytes of buf are used to decode price, second 28 for a hmac sum buffer
	decoded := buf[:28]
	hmacBuf := buf[28:]

	// Decode base64 url
	// Just to be safe remove padding if it was added by mistake
	encryptedPrice = bytes.TrimRight(encryptedPrice, "=")
	if len(encryptedPrice) != 38 {
		return 0, ErrWrongSize
	}
	_, err := base64.RawURLEncoding.Decode(decoded, encryptedPrice)
	if err != nil {
		return 0, err
	}

	// Get elements
	iv := decoded[0:16]
	p := binary.BigEndian.Uint64(decoded[16:24])
	signature := binary.BigEndian.Uint32(decoded[24:28])

	// pad = hmac(e_key, iv)
	padBytes := helpers.HmacSum(dc.encryptionKey, iv, nil, hmacBuf)[:8]
	pad := binary.BigEndian.Uint64(padBytes)

	// priceMicro = p <xor> pad
	priceInMicros := pad ^ p
	priceMicro := [8]byte{}
	binary.BigEndian.PutUint64(priceMicro[:], priceInMicros)

	// conf_sig = hmac(i_key, data || iv)
	confirmationSignatureBytes := helpers.HmacSum(dc.integrityKey, priceMicro[:], iv, hmacBuf)[:4]
	confirmationSignature := binary.BigEndian.Uint32(confirmationSignatureBytes)

	// success = (conf_sig == sig)
	if confirmationSignature != signature {
		return 0, ErrWrongSignature
	}
	return priceInMicros, nil
}
