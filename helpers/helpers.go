package helpers

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"hash"
	"strings"

	"github.com/golang/glog"
)

type KeyDecodingMode string

const (
	Utf8 KeyDecodingMode = "utf-8"
	Hexa                 = "hexa"
)

func CreateHmac(key string, mode KeyDecodingMode) (hash.Hash, error) {
	var err error
	var k []byte

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

func HmacSum(hmac hash.Hash, buf []byte) []byte {
	hmac.Reset()
	hmac.Write(buf)
	return hmac.Sum(nil)
}

func AddBase64Padding(encryptedPrice string) string {
	var base64 string

	base64 = encryptedPrice

	if i := len(base64) % 4; i != 0 {
		base64 += strings.Repeat("=", 4-i)
	}

	return base64
}

func ApplyScaleFactor(price float64, scaleFactor float64, isDebugMode bool) [8]byte {
	scaledPrice := [8]byte{}
	binary.BigEndian.PutUint64(scaledPrice[:], uint64(price*scaleFactor))

	if isDebugMode == true || glog.V(2) {
		glog.Info(fmt.Sprintf("Micro price bytes: %v", scaledPrice))
	}

	return scaledPrice
}