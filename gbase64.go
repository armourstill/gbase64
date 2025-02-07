package main

import (
	"encoding/base64"
	"fmt"
	"strings"
)

type Format = string

const (
	STD Format = "std"
	URL Format = "url"
)

var encoders = map[Format]*base64.Encoding{
	URL: base64.URLEncoding,
	STD: base64.StdEncoding,
}

func encode(enc *base64.Encoding, data []byte) []byte {
	encoded := make([]byte, enc.EncodedLen(len(data)))
	enc.Encode(encoded, data)
	return encoded
}

func decode(enc *base64.Encoding, data []byte) ([]byte, error) {
	decoded := make([]byte, enc.DecodedLen(len(data)))
	_, err := enc.Decode(decoded, data)
	if err != nil {
		return nil, err
	}
	return decoded, nil
}

func getencoder(format Format, nopadding bool) (*base64.Encoding, error) {
	enc, ok := encoders[strings.ToLower(format)]
	if !ok {
		return nil, fmt.Errorf("no format named %s", format)
	}
	if nopadding {
		enc = enc.WithPadding(base64.NoPadding)
	}
	return enc, nil
}

func Encode(format Format, nopadding bool, data []byte) ([]byte, error) {
	enc, err := getencoder(format, nopadding)
	if err != nil {
		return nil, err
	}
	return encode(enc, data), nil
}

func Decode(format Format, data []byte) ([]byte, error) {
	nopadding := true
	if l := len(data); l != 0 && data[l-1] == byte(base64.StdPadding) {
		nopadding = false
	}
	enc, err := getencoder(format, nopadding)
	if err != nil {
		return nil, err
	}
	return decode(enc, data)
}
