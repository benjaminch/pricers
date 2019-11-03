# pricers
[![Build Status](https://travis-ci.org/benjaminch/pricers.svg?branch=master)](https://travis-ci.org/benjaminch/pricers)
[![GoDoc](https://godoc.org/github.com/benjaminch/pricers?status.svg)](https://godoc.org/github.com/benjaminch/pricers)
[![Go Report Card](https://goreportcard.com/badge/github.com/benjaminch/pricers)](https://goreportcard.com/report/github.com/benjaminch/pricers)
[![Maintainability](https://api.codeclimate.com/v1/badges/95e0f8491d86d90c6da6/maintainability)](https://codeclimate.com/github/benjaminch/pricers/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/95e0f8491d86d90c6da6/test_coverage)](https://codeclimate.com/github/benjaminch/pricers/test_coverage)

## Overview
This library supports RTB development for Open RTB common price encryption in Golang.

## Installation
```bash
$ go get github.com/benjaminch/pricers
```

## Supported encryption protocols
### Google Private Data
Specs https://developers.google.com/ad-exchange/rtb/response-guide/decrypt-price
#### Examples
##### Creating a new Google Private Data Pricer
```golang
import "github.com/benjaminch/pricers/doubleclick"

var pricer *doubleclick.DoubleClickPricer
var err error
pricer, err = doubleclick.NewDoubleClickPricer(
    "ZS-DraBUUVeht_sMDgn1nnM3My_nq9TrEESbjubDkTU",   // Encryption key
    "vQo9-4KtlcXmPhWaYvc8asqYuiSVMiGUdZ1RLXfrK7U",   // Integrity key
    true,                                            // Keys are base64
    helpers.Utf8,                                    // Keys should be ingested as Utf-8
    1000000,                                         // Price scale Factor Micro
    false,                                           // No debug
)
```
##### Encrypting a clear price
```golang
import "github.com/benjaminch/pricers/doubleclick"

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
import "github.com/benjaminch/pricers/doubleclick"

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
- [ ] Re-organize directory layout following https://github.com/golang-standards/project-layout
- [ ] Complete documentation:
  - [ ] How to use the Pricer Builder (describing all params)
  - [ ] How to use the Pricer Encrypt function (describing all params)
  - [ ] How to use the Pricer Decrypt function (describing all params)
- [ ] Complete tests for helpers
- [ ] Complete tests for Google Private Data, including various key formats (hex, utf8, base64, etc.)
- [ ] Add other most common price encryption protocols (AES, etc.)
   - [ ] BlowFish
   - [ ] Symetric Algorithm
   - [ ] XOR
