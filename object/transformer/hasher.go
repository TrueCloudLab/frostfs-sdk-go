package transformer

import (
	"crypto/sha256"
	"hash"

	"github.com/TrueCloudLab/frostfs-sdk-go/checksum"
	objectSDK "github.com/TrueCloudLab/frostfs-sdk-go/object"
	"github.com/TrueCloudLab/tzhash/tz"
)

type payloadChecksumHasher struct {
	hasher hash.Hash
	typ    checksum.Type
}

func (h payloadChecksumHasher) writeChecksum(obj *objectSDK.Object) {
	switch h.typ {
	case checksum.SHA256:
		csSHA := [sha256.Size]byte{}
		h.hasher.Sum(csSHA[:0])

		var cs checksum.Checksum
		cs.SetSHA256(csSHA)

		obj.SetPayloadChecksum(cs)
	case checksum.TZ:
		csTZ := [tz.Size]byte{}
		h.hasher.Sum(csTZ[:0])

		var cs checksum.Checksum
		cs.SetTillichZemor(csTZ)

		obj.SetPayloadHomomorphicHash(cs)
	default:
		panic("unreachable")
	}
}
