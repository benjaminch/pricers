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

	"github.com/golang/glog"
)

// KeyDecodingMode : Describing how keys should be decoded.
type KeyDecodingMode string

// String : Returns the KeyDecodingMode string representation.
//
// Returns
//    - string:    String representation of KeyDecodingMode object.
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
//
// ParseKeyDecodingMode requires the following params:
//    - input: 		        keydecodingmode string to parse.
//
// Returns
//    - KeyDecodingMode:    KeyDecodingMode object representation of the input string.
//    - error:              Error occuring while parsing the KeyDecodingMode or nil if no error.
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
//
// CreateHmac requires the following params:
//    - key: 				String key to transform.
//    - isBase64: 			Is key is base 64?
//    - KeyDecodingMode:    How should the key should be decoded?
//
// Returns
//    - hash.Hash:          Key hash.
//    - error:              Error occuring while creating the key or nil if no error.
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
//
// HmacSum requires the following params:
//    - hmac: 		    hmac hash.
//    - buf: 			bytes buffer to write in the hmac.
//
// Returns
//    - byte[]:         Hmac sum bytes.
func HmacSum(hmac hash.Hash, buf []byte) []byte {
	hmac.Reset()
	hmac.Write(buf)
	return hmac.Sum(nil)
}

// AddBase64Padding : Returns base 64 string adding extra padding if needed.
//
// AddBase64Padding requires the following params:
//    - base64Input: 	Base64 string to be padded.
//
// Returns
//    - string:          Padded Base64 string.
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
//
// ApplyScaleFactor requires the following params:
//    - price: 			Price to be scale.
//    - scaleFactor: 	Factor to be applied to the price.
//    - isDebugMode: 	Print debug logs.
//
// Returns
//    - [8]byte:        Scaled priced on 8 bytes.
func ApplyScaleFactor(price float64, scaleFactor float64, isDebugMode bool) [8]byte {
	scaledPrice := [8]byte{}
	binary.BigEndian.PutUint64(scaledPrice[:], uint64(price*scaleFactor))

	if isDebugMode == true || glog.V(2) {
		glog.Info(fmt.Sprintf("Micro price bytes: %v", scaledPrice))
	}

	return scaledPrice
}
