package core

import (
	"encoding/base64"
)

// Base64 is the base64 string representation of a byte array.
// NB: if the passed string is not a valid base64 representation, it will not throw an error but rather returns empty or nil items when methods are called.
type Base64 string

// Bytes ...
func (b Base64) Bytes() []byte {
	if b.String() == "" {
		return nil
	}
	dec, err := base64.StdEncoding.DecodeString(string(b))
	if err != nil {
		return nil
	}
	return dec
}

// String ...
func (b Base64) String() string {
	str := string(b)
	if !IsBase64String(str) {
		return ""
	}
	return str
}

// ToBase64 ...
func ToBase64(bytes []byte) Base64 {
	if bytes == nil {
		return Base64("")
	}
	return Base64(base64.StdEncoding.EncodeToString(bytes))
}

// IsBase64String ...
func IsBase64String(str string) bool {
	_, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return false
	}
	return true
}

// Decoded returns the underlying base64-decoded byte array
func (b Base64) Decoded() []byte {
	if b.String() == "" {
		return nil
	}
	dec, err := base64.StdEncoding.DecodeString(b.String())
	if err != nil {
		return nil
	}
	return dec
}