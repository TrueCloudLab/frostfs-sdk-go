package transformer

import (
	objectSDK "github.com/TrueCloudLab/frostfs-sdk-go/object"
	"github.com/nspcc-dev/neo-go/pkg/util/slice"
)

type chanTarget struct {
	header  *objectSDK.Object
	payload []byte
	ch      chan<- *objectSDK.Object
}

// NewChannelTarget returns ObjectTarget which writes
// object parts to a provided channel.
func NewChannelTarget(ch chan<- *objectSDK.Object) ObjectTarget {
	return &chanTarget{
		ch: ch,
	}
}

// WriteHeader implements the ObjectTarget interface.
func (c *chanTarget) WriteHeader(object *objectSDK.Object) error {
	c.header = object
	return nil
}

// Write implements the ObjectTarget interface.
func (c *chanTarget) Write(p []byte) (n int, err error) {
	c.payload = append(c.payload, p...)
	return len(p), nil
}

// Close implements the ObjectTarget interface.
func (c *chanTarget) Close() (*AccessIdentifiers, error) {
	if len(c.payload) != 0 {
		c.header.SetPayload(slice.Copy(c.payload))
	}
	c.ch <- c.header

	c.header = nil
	c.payload = nil
	return new(AccessIdentifiers), nil
}
