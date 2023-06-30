package transformer

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"fmt"

	"github.com/TrueCloudLab/frostfs-sdk-go/checksum"
	"github.com/TrueCloudLab/frostfs-sdk-go/object"
	oid "github.com/TrueCloudLab/frostfs-sdk-go/object/id"
	"github.com/TrueCloudLab/frostfs-sdk-go/session"
	"github.com/TrueCloudLab/frostfs-sdk-go/version"
	"github.com/TrueCloudLab/tzhash/tz"
)

type payloadSizeLimiter struct {
	Params

	written, writtenCurrent uint64

	current, parent *object.Object

	currentHashers, parentHashers []payloadChecksumHasher

	previous []oid.ID

	splitID *object.SplitID

	parAttrs []object.Attribute
}

type Params struct {
	Key                    *ecdsa.PrivateKey
	NextTarget             ObjectTarget
	SessionToken           *session.Object
	NetworkState           EpochSource
	MaxSize                uint64
	WithoutHomomorphicHash bool
}

// NewPayloadSizeLimiter returns ObjectTarget instance that restricts payload length
// of the writing object and writes generated objects to targets from initializer.
//
// Calculates and adds homomorphic hash to resulting objects only if withoutHomomorphicHash
// is false.
//
// Objects w/ payload size less or equal than max size remain untouched.
func NewPayloadSizeLimiter(p Params) ObjectTarget {
	return &payloadSizeLimiter{
		Params:  p,
		splitID: object.NewSplitID(),
	}
}

func (s *payloadSizeLimiter) WriteHeader(hdr *object.Object) error {
	s.current = fromObject(hdr)

	s.initialize()

	return nil
}

func (s *payloadSizeLimiter) Write(p []byte) (int, error) {
	if err := s.writeChunk(p); err != nil {
		return 0, err
	}

	return len(p), nil
}

func (s *payloadSizeLimiter) Close() (*AccessIdentifiers, error) {
	return s.release(true)
}

func (s *payloadSizeLimiter) initialize() {
	s.current = fromObject(s.current)

	// if it is an object after the 1st
	if ln := len(s.previous); ln > 0 {
		// initialize parent object once (after 1st object)
		if ln == 1 {
			s.parent = fromObject(s.current)
			s.parentHashers = append(s.parentHashers[:0], s.currentHashers...)

			// return source attributes
			s.parent.SetAttributes(s.parAttrs...)
		}

		// set previous object to the last previous identifier
		s.current.SetPreviousID(s.previous[ln-1])
	}

	s.initializeCurrent()
}

func fromObject(obj *object.Object) *object.Object {
	cnr, _ := obj.ContainerID()

	res := object.New()
	res.SetContainerID(cnr)
	res.SetOwnerID(obj.OwnerID())
	res.SetAttributes(obj.Attributes()...)
	res.SetType(obj.Type())

	// obj.SetSplitID creates splitHeader but we don't need to do it in case
	// of small objects, so we should make nil check.
	if obj.SplitID() != nil {
		res.SetSplitID(obj.SplitID())
	}

	return res
}

func (s *payloadSizeLimiter) initializeCurrent() {
	// create payload hashers
	s.writtenCurrent = 0
	s.initPayloadHashers()
}

func (s *payloadSizeLimiter) initPayloadHashers() {
	s.currentHashers = append(s.currentHashers[:0], payloadChecksumHasher{
		hasher: sha256.New(),
		typ:    checksum.SHA256,
	})

	if !s.WithoutHomomorphicHash {
		s.currentHashers = append(s.currentHashers, payloadChecksumHasher{
			hasher: tz.New(),
			typ:    checksum.TZ,
		})
	}
}

func (s *payloadSizeLimiter) release(finalize bool) (*AccessIdentifiers, error) {
	// Arg finalize is true only when called from Close method.
	// We finalize parent and generate linking objects only if it is more
	// than 1 object in split-chain.
	withParent := finalize && len(s.previous) > 0

	if withParent {
		for i := range s.parentHashers {
			s.parentHashers[i].writeChecksum(s.parent)
		}
		s.parent.SetPayloadSize(s.written)
		s.current.SetParent(s.parent)
	}

	// release current object
	for i := range s.currentHashers {
		s.currentHashers[i].writeChecksum(s.current)
	}

	curEpoch := s.NetworkState.CurrentEpoch()
	ver := version.Current()

	s.current.SetVersion(&ver)
	s.current.SetPayloadSize(s.writtenCurrent)
	s.current.SetSessionToken(s.SessionToken)
	s.current.SetCreationEpoch(curEpoch)

	var (
		parID  *oid.ID
		parHdr *object.Object
	)

	if par := s.current.Parent(); par != nil && par.Signature() == nil {
		rawPar := object.NewFromV2(par.ToV2())

		rawPar.SetSessionToken(s.SessionToken)
		rawPar.SetCreationEpoch(curEpoch)

		if err := object.SetIDWithSignature(*s.Key, rawPar); err != nil {
			return nil, fmt.Errorf("could not finalize parent object: %w", err)
		}

		id, _ := rawPar.ID()
		parID = &id
		parHdr = rawPar

		s.current.SetParent(parHdr)
	}

	if err := object.SetIDWithSignature(*s.Key, s.current); err != nil {
		return nil, fmt.Errorf("could not finalize object: %w", err)
	}

	if err := s.NextTarget.WriteHeader(s.current); err != nil {
		return nil, fmt.Errorf("could not write header to next target: %w", err)
	}

	if _, err := s.NextTarget.Close(); err != nil {
		return nil, fmt.Errorf("could not close next target: %w", err)
	}

	id, _ := s.current.ID()

	ids := &AccessIdentifiers{
		ParentID:     parID,
		SelfID:       id,
		ParentHeader: parHdr,
	}

	// save identifier of the released object
	s.previous = append(s.previous, ids.SelfID)

	if withParent {
		// generate and release linking object
		s.initializeLinking(ids.ParentHeader)
		s.initializeCurrent()

		if _, err := s.release(false); err != nil {
			return nil, fmt.Errorf("could not release linking object: %w", err)
		}
	}

	return ids, nil
}

func (s *payloadSizeLimiter) initializeLinking(parHdr *object.Object) {
	s.current = fromObject(s.current)
	s.current.SetParent(parHdr)
	s.current.SetChildren(s.previous...)
	s.current.SetSplitID(s.splitID)
}

func (s *payloadSizeLimiter) writeChunk(chunk []byte) error {
	for {
		// statement is true if the previous write of bytes reached exactly the boundary.
		if s.written > 0 && s.written%s.MaxSize == 0 {
			if s.written == s.MaxSize {
				s.prepareFirstChild()
			}

			// we need to release current object
			if _, err := s.release(false); err != nil {
				return fmt.Errorf("could not release object: %w", err)
			}

			// initialize another object
			s.initialize()
		}

		var (
			ln         = uint64(len(chunk))
			cut        = ln
			leftToEdge = s.MaxSize - s.written%s.MaxSize
		)

		// write bytes no further than the boundary of the current object
		if ln > leftToEdge {
			cut = leftToEdge
		}

		if err := s.writeHashes(chunk[:cut]); err != nil {
			return fmt.Errorf("could not write chunk to target: %w", err)
		}

		// increase written bytes counter
		s.writtenCurrent += cut
		s.written += cut

		if cut == ln {
			return nil
		}
		// if there are more bytes in buffer we call method again to start filling another object
		chunk = chunk[cut:]
	}
}

func (s *payloadSizeLimiter) writeHashes(chunk []byte) error {
	_, err := s.NextTarget.Write(chunk)
	if err != nil {
		return err
	}

	// The `Write` method of `hash.Hash` never returns an error.
	for i := range s.currentHashers {
		_, _ = s.currentHashers[i].hasher.Write(chunk)
	}

	for i := range s.parentHashers {
		_, _ = s.parentHashers[i].hasher.Write(chunk)
	}

	return nil
}

func (s *payloadSizeLimiter) prepareFirstChild() {
	// initialize split header with split ID on first object in chain
	s.current.InitRelations()
	s.current.SetSplitID(s.splitID)

	// cut source attributes
	s.parAttrs = s.current.Attributes()
	s.current.SetAttributes()

	// attributes will be added to parent in detachParent
}
