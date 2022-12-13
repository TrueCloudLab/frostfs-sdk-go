package ns

import (
	"net"

	cid "github.com/TrueCloudLab/frostfs-sdk-go/container/id"
)

// DNS looks up NeoFS names using system DNS.
//
// See also net package.
type DNS struct{}

// ResolveContainerName looks up for DNS TXT records for the given domain name
// and returns the first one which represents valid container ID in a string format.
// Otherwise, returns an error.
//
// See also net.LookupTXT.
func (x *DNS) ResolveContainerName(name string) (id cid.ID, err error) {
	records, err := net.LookupTXT(name)
	if err != nil {
		return
	}

	for i := range records {
		err = id.DecodeString(records[i])
		if err == nil {
			return
		}
	}

	err = errNotFound

	return
}
