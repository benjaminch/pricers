# openrtb-pricers
[![Build Status](https://travis-ci.org/benjaminch/openrtb-pricers.svg?branch=master)](https://travis-ci.org/benjaminch/openrtb-pricers)
[![GoDoc](https://godoc.org/github.com/benjaminch/openrtb-pricers?status.svg)](https://godoc.org/github.com/benjaminch/openrtb-pricers)

## Overview
This library supports RTB development for Open RTB common price encryption in Golang.

## Installation
```bash
$ go get github.com/benjaminch/openrtb-pricers
```

## Supported encryption protocols
### Google Private Data
Specs https://developers.google.com/ad-exchange/rtb/response-guide/decrypt-price
#### Examples
##### Creating a new Google Private Data Pricer
```golang
import "github.com/benjaminch/openrtb-pricers/doubleclick"

var pricer *doubleclick.DoubleClickPricer
var err error
pricer, err = doubleclick.NewDoubleClickPricer(
    "skU7Ax_NL5pPAFyKdkfZjZz2-VhIN8bjj1rVFOaJ_5o=",  // Encryption key
    "arO23ykdNqUQ5LEoQ0FVmPkBd7xB5CO89PDZlSjpFxo=",  // Integrity key
    true,                                            // Keys are base64
    helpers.Utf8,                                    // Keys should be ingested as Utf-8
    1000000,                                         // Price scale Factor Micro
    false,                                           // No debug
)
```
##### Encrypting a clear price
```golang
import "github.com/benjaminch/openrtb-pricers/doubleclick"

var result string
var err error
result, err = pricer.Encrypt(
    "",    // Seed
    1,     // Clear price
    false  // No debug
)
if err != nil {
    err = errors.New("Encryption failed. Error : %s", err)
}
```
##### Decrypting an encrypted price
```golang
import "github.com/benjaminch/openrtb-pricers/doubleclick"

var result float64
var err error
result, err = pricer.Decrypt(
    "WEp8nQAAAAADG-y45xxIC1tMWuTjzmDW6HtroQ",  // Encrypted price
    false,                                     // No debug
)
if err != nil {
    err = errors.New("Decryption failed. Error : %s", err)
}
```
## Todos
- [ ] Complete documentation:
  - [ ] How to use the Pricer Builder (describing all params)
  - [ ] How to use the Pricer Encrypt function (describing all params)
  - [ ] How to use the Pricer Decrypt function (describing all params)
- [ ] Complete tests for helpers
- [ ] Complete tests for Google Private Data, including various key formats (hex, utf8, base64, etc.)
- [ ] Add other most common price encryption protocols (AES, etc.)
