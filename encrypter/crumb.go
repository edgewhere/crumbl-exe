package encrypter

import (
	"crumbl/models/core"
	"crumbl/utils"
	"errors"
	"fmt"
)

//--- TYPES

// Crumb holds the encrypted slice, its index and length.
type Crumb struct {
	Encrypted core.Base64
	Index     int
	Length    int
}

//--- METHODS

// String transforms the Crumb into its stringified representation.
// The construct is as follows:
// - the first two characters are the hexadecimal representation of the index;
// - the following four characters are the hexadecimal representation of the length of the encrypted data to follow;
// - the base64-encoded string.
// NB: the condition to only use four characters for the length of the encrypted data implies that
// this encrypted crumb shouldn't be longer than 65 535 characters, ie. 64 ko.
func (c *Crumb) String() string {
	return fmt.Sprintf("%02x", c.Index) + fmt.Sprintf("%04x", c.Length) + c.Encrypted.String()
}

// ToCrumb ...
func ToCrumb(unparsed string) (c Crumb, err error) {
	idx, length, enc, err := Parse(unparsed)
	if err != nil {
		return
	}
	c = Crumb{
		Encrypted: core.Base64(enc),
		Index:     idx,
		Length:    length,
	}
	return
}

// Parse extracts the index, the encrypted length and text from the passed string, or returns an error
func Parse(unparsed string) (index int, length int, enc string, err error) {
	idxHex := unparsed[:2]
	idx, err := utils.HexToInt(idxHex)
	if err != nil {
		return
	}
	lnHex := unparsed[2:6]
	ln, err := utils.HexToInt(lnHex)
	if err != nil {
		return
	}
	e := unparsed[6:]
	if !core.IsBase64String(e) {
		err = errors.New("not a base64-encoded string")
		return
	}
	if ln != len(e) {
		err = errors.New("incompatible lengths")
		return
	}
	index = idx
	length = ln
	enc = e
	return
}
