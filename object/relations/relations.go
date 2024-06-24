package relations

import (
	"context"
	"errors"
	"fmt"

	"github.com/TrueCloudLab/frostfs-sdk-go/bearer"
	cid "github.com/TrueCloudLab/frostfs-sdk-go/container/id"
	"github.com/TrueCloudLab/frostfs-sdk-go/object"
	oid "github.com/TrueCloudLab/frostfs-sdk-go/object/id"
	"github.com/TrueCloudLab/frostfs-sdk-go/session"
)

// Tokens contains different tokens to perform requests in Relations implementations.
type Tokens struct {
	Session *session.Object
	Bearer  *bearer.Token
}

type Relations interface {
	// GetSplitInfo tries to get split info by some object id.
	// This method must return split info on any object from split chain as well as on parent/linking object.
	// If object doesn't have any split information returns ErrNoSplitInfo.
	GetSplitInfo(ctx context.Context, cnrID cid.ID, rootID oid.ID, tokens Tokens) (*object.SplitInfo, error)

	// ListChildrenByLinker returns list of children for link object.
	// Result doesn't include link object itself.
	ListChildrenByLinker(ctx context.Context, cnrID cid.ID, linkerID oid.ID, tokens Tokens) ([]oid.ID, error)

	// GetLeftSibling return previous object id in object chain.
	// If no previous object it returns ErrNoLeftSibling.
	GetLeftSibling(ctx context.Context, cnrID cid.ID, objID oid.ID, tokens Tokens) (oid.ID, error)

	// FindSiblingByParentID returns all object that relates to the provided parent id.
	FindSiblingByParentID(ctx context.Context, cnrID cid.ID, parentID oid.ID, tokens Tokens) ([]oid.ID, error)
}

var (
	// ErrNoLeftSibling is an error that must be returned if object doesn't have left sibling in objects chain.
	ErrNoLeftSibling = errors.New("no left siblings")

	// ErrNoSplitInfo is an error that must be returned if requested object isn't virtual.
	ErrNoSplitInfo = errors.New("no split info")
)

// ListAllRelations return all related phy objects for provided root object ID in split-chain order.
// Result doesn't include root object ID itself. If linking object is found its id will be the last one.
func ListAllRelations(ctx context.Context, rels Relations, cnrID cid.ID, rootObjID oid.ID, tokens Tokens) ([]oid.ID, error) {
	return ListRelations(ctx, rels, cnrID, rootObjID, tokens, true)
}

// ListRelations return all related phy objects for provided root object ID in split-chain order.
// Result doesn't include root object ID itself.
func ListRelations(ctx context.Context, rels Relations, cnrID cid.ID, rootObjID oid.ID, tokens Tokens, includeLinking bool) ([]oid.ID, error) {
	splitInfo, err := rels.GetSplitInfo(ctx, cnrID, rootObjID, tokens)
	if err != nil {
		if errors.Is(err, ErrNoSplitInfo) {
			return []oid.ID{}, nil
		}
		return nil, err
	}

	// collect split chain by the descending ease of operations (ease is evaluated heuristically).
	// If any approach fails, we don't try the next since we assume that it will fail too.
	if _, ok := splitInfo.Link(); !ok {
		// the list is expected to contain last part and (probably) split info
		list, err := rels.FindSiblingByParentID(ctx, cnrID, rootObjID, tokens)
		if err != nil {
			return nil, fmt.Errorf("failed to find object children: %w", err)
		}

		for _, id := range list {
			split, err := rels.GetSplitInfo(ctx, cnrID, id, tokens)
			if err != nil {
				if errors.Is(err, ErrNoSplitInfo) {
					continue
				}
				return nil, fmt.Errorf("failed to get split info: %w", err)
			}
			if link, ok := split.Link(); ok {
				splitInfo.SetLink(link)
			}
			if last, ok := split.LastPart(); ok {
				splitInfo.SetLastPart(last)
			}
		}
	}

	if idLinking, ok := splitInfo.Link(); ok {
		children, err := rels.ListChildrenByLinker(ctx, cnrID, idLinking, tokens)
		if err != nil {
			return nil, fmt.Errorf("failed to get linking object's header: %w", err)
		}

		if includeLinking {
			children = append(children, idLinking)
		}
		return children, nil
	}

	idMember, ok := splitInfo.LastPart()
	if !ok {
		return nil, errors.New("missing any data in received object split information")
	}

	chain := []oid.ID{idMember}
	chainSet := map[oid.ID]struct{}{idMember: {}}

	for {
		idMember, err = rels.GetLeftSibling(ctx, cnrID, idMember, tokens)
		if err != nil {
			if errors.Is(err, ErrNoLeftSibling) {
				break
			}
			return nil, fmt.Errorf("failed to read split chain member's header: %w", err)
		}

		if _, ok = chainSet[idMember]; ok {
			return nil, fmt.Errorf("duplicated member in the split chain %s", idMember)
		}

		chain = append([]oid.ID{idMember}, chain...)
		chainSet[idMember] = struct{}{}
	}

	return chain, nil
}
