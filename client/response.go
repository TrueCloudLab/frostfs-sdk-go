package client

import "github.com/TrueCloudLab/frostfs-api-go/v2/session"

// ResponseMetaInfo groups meta information about any NeoFS API response.
type ResponseMetaInfo struct {
	key []byte

	epoch uint64
}

type responseV2 interface {
	GetMetaHeader() *session.ResponseMetaHeader
	GetVerificationHeader() *session.ResponseVerificationHeader
}

// ResponderKey returns responder's public key in a binary format.
//
// Result must not be mutated.
func (x ResponseMetaInfo) ResponderKey() []byte {
	return x.key
}

// Epoch returns local NeoFS epoch of the server.
func (x ResponseMetaInfo) Epoch() uint64 {
	return x.epoch
}
