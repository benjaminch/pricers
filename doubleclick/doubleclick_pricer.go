package doubleclick

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"errors"

	"../helpers"

	"github.com/golang/glog"
)

// DoubleClickPricer implementing price encryption and decryption
// Specs : https://developers.google.com/ad-exchange/rtb/response-guide/decrypt-price
type DoubleClickPricer struct {
	encryptionKey   string
	integrityKey    string
	keyDecodingMode helpers.KeyDecodingMode
	scaleFactor     float64
	isDebugMode     bool
}

func NewDoubleClickPricer(encryptionKey, integrityKey string, keyDecodingMode helpers.KeyDecodingMode, scaleFactor float64, isDebugMode bool) (*DoubleClickPricer, error) {
	var err error
	return &DoubleClickPricer{encryptionKey: encryptionKey, keyDecodingMode: keyDecodingMode, scaleFactor: scaleFactor, isDebugMode: isDebugMode}, err
}

// Encrypt is returned by FormFile when the provided file field name
// is either not present in the request or not a file field.
func (dc *DoubleClickPricer) Encrypt(encryptionKey, integrityKey string, keyDecodingMode helpers.KeyDecodingMode, seed string, price float64, scaleFactor float64, isDebugMode bool) string {
	encryptingFun, _ := helpers.CreateHmac(dc.encryptionKey, dc.keyDecodingMode)
	integrityFun, _ := helpers.CreateHmac(dc.integrityKey, dc.keyDecodingMode)

	// Result
	var (
		iv        [16]byte
		encoded   [8]byte
		signature [4]byte
	)

	if isDebugMode == true {
		fmt.Println("Keys decoding mode : ", keyDecodingMode)
		fmt.Println("Encryption key : ", encryptionKey)
		encryptionKeyHexa, _ := hex.DecodeString(encryptionKey)
		fmt.Println("Encryption key (bytes) : ", []byte(encryptionKeyHexa))
		fmt.Println("Integrity key : ", integrityKey)
		integrityKeyHexa, _ := hex.DecodeString(integrityKey)
		fmt.Println("Integrity key (bytes) : ", []byte(integrityKeyHexa))
	}

	data := helpers.ApplyScaleFactor(price, scaleFactor, isDebugMode)

	// Create Initialization Vector from seed
	sum := md5.Sum([]byte(seed))
	copy(iv[:], sum[:])
	if isDebugMode == true {
		fmt.Println("Seed : ", seed)
		fmt.Println("Initialization vector : ", iv)
	}

	//pad = hmac(e_key, iv), first 8 bytes
	pad := helpers.HmacSum(encryptingFun, iv[:])[:8]
	if isDebugMode == true {
		fmt.Println("// pad = hmac(e_key, iv), first 8 bytes")
		fmt.Println("Pad : ", pad)
	}

	// enc_data = pad <xor> data
	for i := range data {
		encoded[i] = pad[i] ^ data[i]
	}
	if isDebugMode == true {
		fmt.Println("// enc_data = pad <xor> data")
		fmt.Println("Encoded price bytes : ", encoded)
	}

	// signature = hmac(i_key, data || iv), first 4 bytes
	sig := helpers.HmacSum(integrityFun, append(data[:], iv[:]...))[:4]
	copy(signature[:], sig[:])
	if isDebugMode == true {
		fmt.Println("// signature = hmac(i_key, data || iv), first 4 bytes")
		fmt.Println("Signature : ", sig)
	}

	glog.Flush()

	// final_message = WebSafeBase64Encode( iv || enc_price || signature )
	return base64.URLEncoding.EncodeToString(append(append(iv[:], encoded[:]...), signature[:]...))
}

func (dc *DoubleClickPricer) Decrypt(encryptedPrice string) (float64, error) {
	var err error
	encryptingFun, _ := helpers.CreateHmac(dc.encryptionKey, dc.keyDecodingMode)
	integrityFun, _ := helpers.CreateHmac(dc.integrityKey, dc.keyDecodingMode)

	// Decode base64
	decoded, _ := base64.URLEncoding.DecodeString(encryptedPrice)

	// Get elements
	var (
		iv         [16]byte
		p          [8]byte
		signature  [4]byte
		priceMicro [8]byte
	)
	copy(iv[:], decoded[0:16])
	copy(p[:], decoded[16:24])
	copy(signature[:], decoded[24:28])

	// pad = hmac(e_key, iv)
	pad := helpers.HmacSum(encryptingFun, iv[:])[:8]

	// priceMicro = p <xor> pad
	for i := range p {
		priceMicro[i] = pad[i] ^ p[i]
	}

	// conf_sig = hmac(i_key, data || iv)
	sig := helpers.HmacSum(integrityFun, append(priceMicro[:], iv[:]...))[:4]

	// success = (conf_sig == sig)
	for i := range sig {
		if sig[i] != signature[i] {
			return 0, errors.New("Failed to decrypt")
		}
	}
	price := float64(binary.BigEndian.Uint64(priceMicro[:])) / dc.scaleFactor

	glog.Flush()

	return price, err
}
