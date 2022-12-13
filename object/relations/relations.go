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
	// GetSplitInfo tries to get split info by root object id.
	// If object isn't virtual it returns ErrNoSplitInfo.
	GetSplitInfo(ctx context.Context, cnrID cid.ID, rootID oid.ID, tokens Tokens) (*object.SplitInfo, error)

	// ListChildrenByLinker returns list of children for link object.
	// Result doesn't include link object itself.
	ListChildrenByLinker(ctx context.Context, cnrID cid.ID, linkerID oid.ID, tokens Tokens) ([]oid.ID, error)

	// GetLeftSibling return previous object id in object chain.
	// If no previous object it returns ErrNoLeftSibling.
	GetLeftSibling(ctx context.Context, cnrID cid.ID, objID oid.ID, tokens Tokens) (oid.ID, error)

	// FindSiblingBySplitID returns all objects that relates to the provided split id.
	FindSiblingBySplitID(ctx context.Context, cnrID cid.ID, splitID *object.SplitID, tokens Tokens) ([]oid.ID, error)

	// FindSiblingByParentID returns all object that relates to the provided parent id.
	FindSiblingByParentID(ctx context.Context, cnrID cid.ID, parentID oid.ID, tokens Tokens) ([]oid.ID, error)
}

var (
	// ErrNoLeftSibling is an error that must be returned if object doesn't have left sibling in objects chain.
	ErrNoLeftSibling = errors.New("no left siblings")

	// ErrNoSplitInfo is an error that must be returned if requested object isn't virtual.
	ErrNoSplitInfo = errors.New("no split info")
)

// ListAllRelations return all related phy objects for provided root object ID.
// Result doesn't include root object ID itself.
func ListAllRelations(ctx context.Context, rels Relations, cnrID cid.ID, rootObjID oid.ID, tokens Tokens) ([]oid.ID, error) {
	splitInfo, err := rels.GetSplitInfo(ctx, cnrID, rootObjID, tokens)
	if err != nil {
		if errors.Is(err, ErrNoSplitInfo) {
			return []oid.ID{}, nil
		}
		return nil, err
	}

	// collect split chain by the descending ease of operations (ease is evaluated heuristically).
	// If any approach fails, we don't try the next since we assume that it will fail too.
	if idLinking, ok := splitInfo.Link(); ok {
		children, err := rels.ListChildrenByLinker(ctx, cnrID, idLinking, tokens)
		if err != nil {
			return nil, fmt.Errorf("failed to get linking object's header: %w", err)
		}

		// include linking object
		return append(children, idLinking), nil
	}

	if idSplit := splitInfo.SplitID(); idSplit != nil {
		members, err := rels.FindSiblingBySplitID(ctx, cnrID, idSplit, tokens)
		if err != nil {
			return nil, fmt.Errorf("failed to search objects by split ID: %w", err)
		}
		return members, nil
	}

	idMember, ok := splitInfo.LastPart()
	if !ok {
		return nil, errors.New("missing any data in received object split information")
	}

	chain := []oid.ID{idMember}
	chainSet := map[oid.ID]struct{}{idMember: {}}

	// prmHead.SetRawFlag(false)
	// split members are almost definitely singular, but don't get hung up on it

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

		chain = append(chain, idMember)
		chainSet[idMember] = struct{}{}
	}

	list, err := rels.FindSiblingByParentID(ctx, cnrID, rootObjID, tokens)
	if err != nil {
		return nil, fmt.Errorf("failed to find object children: %w", err)
	}

	for i := range list {
		if _, ok = chainSet[list[i]]; !ok {
			chain = append(chain, list[i])
		}
	}

	return chain, nil
}
