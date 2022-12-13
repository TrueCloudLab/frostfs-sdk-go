package frostfscrypto

import "encoding/hex"

// StringifyKeyBinary returns string with HEX representation of source.
// Format can be changed and it's unsafe to rely on it beyond human-readable output.
func StringifyKeyBinary(src []byte) string {
	return hex.EncodeToString(src)
}
