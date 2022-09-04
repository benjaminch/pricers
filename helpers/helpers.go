package helpers

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"strings"
)

// KeyDecodingMode : Describing how keys should be decoded.
type KeyDecodingMode string

// String : Returns the KeyDecodingMode string representation.
func (kd KeyDecodingMode) String() string {
	return string(kd)
}

const (
	// Utf8 : Key should be decoded as utf-8 string.
	Utf8 KeyDecodingMode = "utf-8"
	// Hexa : Key should be decoded as hexa string.
	Hexa KeyDecodingMode = "hexa"
)

// ParseKeyDecodingMode : Parses KeyDecodingMode from string.
func ParseKeyDecodingMode(input string) (KeyDecodingMode, error) {
	var err error
	var parsed KeyDecodingMode

	if input == "" {
		err = errors.New("input is empty, cannot parse empty input")
	} else {
		switch input {
		case Utf8.String():
			parsed = Utf8
			break
		case Hexa.String():
			parsed = Hexa
			break
		default:
			err = errors.New("input doesn't match to any key decoding mode")
		}
	}

	return parsed, err
}

// CreateHmac : Returns Hash from input string.
func CreateHmac(key string, isBase64 bool, mode KeyDecodingMode) (hash.Hash, error) {
	var err error
	var b64DecodedKey []byte
	var k []byte

	if isBase64 {
		b64DecodedKey, err = base64.URLEncoding.DecodeString(AddBase64Padding(key))
		if err == nil {
			// If no error, then use the base 64 decoded key
			key = string(b64DecodedKey[:])
		}
	}

	if mode == Utf8 {
		k = []byte(key)
	} else {
		k, err = hex.DecodeString(key)
	}

	if err != nil {
		return nil, err
	}

	return hmac.New(sha1.New, k), nil
}

// HmacSum : Returns Hmac sum bytes.
func HmacSum(hmac hash.Hash, buf []byte) []byte {
	hmac.Reset()
	hmac.Write(buf)
	return hmac.Sum(nil)
}

// AddBase64Padding : Returns base 64 string adding extra padding if needed.
func AddBase64Padding(base64Input string) string {
	var base64 string

	base64 = base64Input

	if i := len(base64) % 4; i != 0 {
		base64 += strings.Repeat("=", 4-i)
	}

	return base64
}

// ApplyScaleFactor : Applies a scale factor to a given price.
// Scaled price will be represented on 8 bytes.
func ApplyScaleFactor(price float64, scaleFactor float64, isDebugMode bool) [8]byte {
	scaledPrice := [8]byte{}
	binary.BigEndian.PutUint64(scaledPrice[:], uint64(price*scaleFactor))

	if isDebugMode {
		fmt.Printf("Micro price bytes: %v", scaledPrice)
	}

	return scaledPrice
}
