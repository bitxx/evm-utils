package util

import (
	"encoding/hex"
	"strings"
)

func HexDecodeString(s string) ([]byte, error) {
	s = strings.TrimPrefix(s, "0x")

	if len(s)%2 != 0 {
		s = "0" + s
	}

	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// HexEncode encodes bytes to a hex string. Contrary to hex.EncodeToString, this function prefixes the hex string
// with "0x"
func HexEncodeToString(b []byte) string {
	return "0x" + hex.EncodeToString(b)
}
