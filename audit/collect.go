package audit

import (
	"context"
	"fmt"

	"github.com/TrueCloudLab/frostfs-sdk-go/checksum"
	cid "github.com/TrueCloudLab/frostfs-sdk-go/container/id"
	"github.com/TrueCloudLab/frostfs-sdk-go/object"
	oid "github.com/TrueCloudLab/frostfs-sdk-go/object/id"
	"github.com/TrueCloudLab/frostfs-sdk-go/object/relations"
	"github.com/TrueCloudLab/frostfs-sdk-go/storagegroup"
	"github.com/TrueCloudLab/tzhash/tz"
)

type Collector interface {
	Head(ctx context.Context, addr oid.Address) (*object.Object, error)
	relations.Relations
}

// CollectMembers creates new storage group structure and fills it
// with information about members collected via HeadReceiver.
//
// Resulting storage group consists of physically stored objects only.
func CollectMembers(ctx context.Context, collector Collector, cnr cid.ID, members []oid.ID, tokens relations.Tokens, calcHomoHash bool) (*storagegroup.StorageGroup, error) {
	var (
		err        error
		sumPhySize uint64
		phyMembers []oid.ID
		phyHashes  [][]byte
		addr       oid.Address
		sg         storagegroup.StorageGroup
	)

	addr.SetContainer(cnr)

	for i := range members {
		if phyMembers, err = relations.ListRelations(ctx, collector, cnr, members[i], tokens, false); err != nil {
			return nil, err
		}

		for _, phyMember := range phyMembers {
			addr.SetObject(phyMember)
			leaf, err := collector.Head(ctx, addr)
			if err != nil {
				return nil, fmt.Errorf("head phy member '%s': %w", phyMember.EncodeToString(), err)
			}

			sumPhySize += leaf.PayloadSize()
			cs, _ := leaf.PayloadHomomorphicHash()

			if calcHomoHash {
				phyHashes = append(phyHashes, cs.Value())
			}
		}
	}

	sg.SetMembers(phyMembers)
	sg.SetValidationDataSize(sumPhySize)

	if calcHomoHash {
		sumHash, err := tz.Concat(phyHashes)
		if err != nil {
			return nil, err
		}

		var cs checksum.Checksum
		tzHash := [64]byte{}
		copy(tzHash[:], sumHash)
		cs.SetTillichZemor(tzHash)

		sg.SetValidationDataHash(cs)
	}

	return &sg, nil
}
