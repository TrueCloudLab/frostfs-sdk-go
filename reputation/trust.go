package reputation

import (
	"errors"
	"fmt"

	"github.com/TrueCloudLab/frostfs-api-go/v2/refs"
	"github.com/TrueCloudLab/frostfs-api-go/v2/reputation"
	frostfscrypto "github.com/TrueCloudLab/frostfs-sdk-go/crypto"
	"github.com/TrueCloudLab/frostfs-sdk-go/version"
)

// Trust represents quantitative assessment of the trust of a participant in the
// NeoFS reputation system.
//
// Trust is mutually compatible with github.com/TrueCloudLab/frostfs-api-go/v2/reputation.Trust
// message. See ReadFromV2 / WriteToV2 methods.
//
// Instances can be created using built-in var declaration.
type Trust struct {
	m reputation.Trust
}

// ReadFromV2 reads Trust from the reputation.Trust message. Returns an
// error if the message is malformed according to the NeoFS API V2 protocol.
//
// See also WriteToV2.
func (x *Trust) ReadFromV2(m reputation.Trust) error {
	if val := m.GetValue(); val < 0 || val > 1 {
		return fmt.Errorf("invalid trust value %v", val)
	}

	peerV2 := m.GetPeer()
	if peerV2 == nil {
		return errors.New("missing peer field")
	}

	var peer PeerID

	err := peer.ReadFromV2(*peerV2)
	if err != nil {
		return fmt.Errorf("invalid peer field: %w", err)
	}

	x.m = m

	return nil
}

// WriteToV2 writes Trust to the reputation.Trust message.
// The message must not be nil.
//
// See also ReadFromV2.
func (x Trust) WriteToV2(m *reputation.Trust) {
	*m = x.m
}

// SetPeer specifies identifier of the participant of the NeoFS reputation system
// to which the Trust relates.
//
// See also Peer.
func (x *Trust) SetPeer(id PeerID) {
	var m reputation.PeerID
	id.WriteToV2(&m)

	x.m.SetPeer(&m)
}

// Peer returns peer identifier set using SetPeer.
//
// Zero Trust returns zero PeerID which is incorrect according to the NeoFS API
// protocol.
func (x Trust) Peer() (res PeerID) {
	m := x.m.GetPeer()
	if m != nil {
		err := res.ReadFromV2(*m)
		if err != nil {
			panic(fmt.Sprintf("unexpected error from ReadFromV2: %v", err))
		}
	}

	return
}

// SetValue sets the Trust value. Value MUST be in range [0;1].
//
// See also Value.
func (x *Trust) SetValue(val float64) {
	if val < 0 || val > 1 {
		panic(fmt.Sprintf("trust value is out-of-range %v", val))
	}

	x.m.SetValue(val)
}

// Value returns value set using SetValue.
//
// Zero Trust has zero value.
func (x Trust) Value() float64 {
	return x.m.GetValue()
}

// PeerToPeerTrust represents trust of one participant of the NeoFS reputation
// system to another one.
//
// Trust is mutually compatible with github.com/TrueCloudLab/frostfs-api-go/v2/reputation.PeerToPeerTrust
// message. See ReadFromV2 / WriteToV2 methods.
//
// Instances can be created using built-in var declaration.
type PeerToPeerTrust struct {
	m reputation.PeerToPeerTrust
}

// ReadFromV2 reads PeerToPeerTrust from the reputation.PeerToPeerTrust message.
// Returns an error if the message is malformed according to the NeoFS API V2
// protocol.
//
// See also WriteToV2.
func (x *PeerToPeerTrust) ReadFromV2(m reputation.PeerToPeerTrust) error {
	trustingV2 := m.GetTrustingPeer()
	if trustingV2 == nil {
		return errors.New("missing trusting peer")
	}

	var trusting PeerID

	err := trusting.ReadFromV2(*trustingV2)
	if err != nil {
		return fmt.Errorf("invalid trusting peer: %w", err)
	}

	trustV2 := m.GetTrust()
	if trustV2 == nil {
		return errors.New("missing trust")
	}

	var trust Trust

	err = trust.ReadFromV2(*trustV2)
	if err != nil {
		return fmt.Errorf("invalid trust: %w", err)
	}

	x.m = m

	return nil
}

// WriteToV2 writes PeerToPeerTrust to the reputation.PeerToPeerTrust message.
// The message must not be nil.
//
// See also ReadFromV2.
func (x PeerToPeerTrust) WriteToV2(m *reputation.PeerToPeerTrust) {
	*m = x.m
}

// SetTrustingPeer specifies the peer from which trust comes in terms of the
// NeoFS reputation system.
//
// See also TrustingPeer.
func (x *PeerToPeerTrust) SetTrustingPeer(id PeerID) {
	var m reputation.PeerID
	id.WriteToV2(&m)

	x.m.SetTrustingPeer(&m)
}

// TrustingPeer returns peer set using SetTrustingPeer.
//
// Zero PeerToPeerTrust has no trusting peer which is incorrect according
// to the NeoFS API protocol.
func (x PeerToPeerTrust) TrustingPeer() (res PeerID) {
	m := x.m.GetTrustingPeer()
	if m != nil {
		err := res.ReadFromV2(*m)
		if err != nil {
			panic(fmt.Sprintf("unexpected error from PeerID.ReadFromV2: %v", err))
		}
	}

	return
}

// SetTrust sets trust value of the trusting peer to another participant
// of the NeoFS reputation system.
//
// See also Trust.
func (x *PeerToPeerTrust) SetTrust(t Trust) {
	var tV2 reputation.Trust
	t.WriteToV2(&tV2)

	x.m.SetTrust(&tV2)
}

// Trust returns trust set using SetTrust.
//
// Zero PeerToPeerTrust returns zero Trust which is incorect according to the
// NeoFS API protocol.
func (x PeerToPeerTrust) Trust() (res Trust) {
	m := x.m.GetTrust()
	if m != nil {
		err := res.ReadFromV2(*m)
		if err != nil {
			panic(fmt.Sprintf("unexpected error from Trust.ReadFromV2: %v", err))
		}
	}

	return
}

// GlobalTrust represents the final assessment of trust in the participant of
// the NeoFS reputation system obtained taking into account all other participants.
//
// GlobalTrust is mutually compatible with github.com/TrueCloudLab/frostfs-api-go/v2/reputation.GlobalTrust
// message. See ReadFromV2 / WriteToV2 methods.
//
// To submit GlobalTrust value in NeoFS zero instance SHOULD be declared,
// initialized using Init method and filled using dedicated methods.
type GlobalTrust struct {
	m reputation.GlobalTrust
}

// ReadFromV2 reads GlobalTrust from the reputation.GlobalTrust message.
// Returns an error if the message is malformed according to the NeoFS API V2
// protocol.
//
// See also WriteToV2.
func (x *GlobalTrust) ReadFromV2(m reputation.GlobalTrust) error {
	if m.GetVersion() == nil {
		return errors.New("missing version")
	}

	if m.GetSignature() == nil {
		return errors.New("missing signature")
	}

	body := m.GetBody()
	if body == nil {
		return errors.New("missing body")
	}

	managerV2 := body.GetManager()
	if managerV2 == nil {
		return errors.New("missing manager")
	}

	var manager PeerID

	err := manager.ReadFromV2(*managerV2)
	if err != nil {
		return fmt.Errorf("invalid manager: %w", err)
	}

	trustV2 := body.GetTrust()
	if trustV2 == nil {
		return errors.New("missing trust")
	}

	var trust Trust

	err = trust.ReadFromV2(*trustV2)
	if err != nil {
		return fmt.Errorf("invalid trust: %w", err)
	}

	x.m = m

	return nil
}

// WriteToV2 writes GlobalTrust to the reputation.GlobalTrust message.
// The message must not be nil.
//
// See also ReadFromV2.
func (x GlobalTrust) WriteToV2(m *reputation.GlobalTrust) {
	*m = x.m
}

// Init initializes all internal data of the GlobalTrust required by NeoFS API
// protocol. Init MUST be called when creating a new global trust instance.
// Init SHOULD NOT be called multiple times. Init SHOULD NOT be called if
// the GlobalTrust instance is used for decoding only.
func (x *GlobalTrust) Init() {
	var ver refs.Version
	version.Current().WriteToV2(&ver)

	x.m.SetVersion(&ver)
}

func (x *GlobalTrust) setBodyField(setter func(*reputation.GlobalTrustBody)) {
	if x != nil {
		body := x.m.GetBody()
		if body == nil {
			body = new(reputation.GlobalTrustBody)
			x.m.SetBody(body)
		}

		setter(body)
	}
}

// SetManager sets identifier of the NeoFS reputation system's participant which
// performed trust estimation.
//
// See also Manager.
func (x *GlobalTrust) SetManager(id PeerID) {
	var m reputation.PeerID
	id.WriteToV2(&m)

	x.setBodyField(func(body *reputation.GlobalTrustBody) {
		body.SetManager(&m)
	})
}

// Manager returns peer set using SetManager.
//
// Zero GlobalTrust has zero manager which is incorrect according to the
// NeoFS API protocol.
func (x GlobalTrust) Manager() (res PeerID) {
	m := x.m.GetBody().GetManager()
	if m != nil {
		err := res.ReadFromV2(*m)
		if err != nil {
			panic(fmt.Sprintf("unexpected error from ReadFromV2: %v", err))
		}
	}

	return
}

// SetTrust sets the global trust score of the network to a specific network
// member.
//
// See also Trust.
func (x *GlobalTrust) SetTrust(trust Trust) {
	var m reputation.Trust
	trust.WriteToV2(&m)

	x.setBodyField(func(body *reputation.GlobalTrustBody) {
		body.SetTrust(&m)
	})
}

// Trust returns trust set using SetTrust.
//
// Zero GlobalTrust return zero Trust which is incorrect according to the
// NeoFS API protocol.
func (x GlobalTrust) Trust() (res Trust) {
	m := x.m.GetBody().GetTrust()
	if m != nil {
		err := res.ReadFromV2(*m)
		if err != nil {
			panic(fmt.Sprintf("unexpected error from ReadFromV2: %v", err))
		}
	}

	return
}

// Sign calculates and writes signature of the GlobalTrust data. Returns
// signature calculation errors.
//
// Zero GlobalTrust is unsigned.
//
// Note that any GlobalTrust mutation is likely to break the signature, so it is
// expected to be calculated as a final stage of GlobalTrust formation.
//
// See also VerifySignature.
func (x *GlobalTrust) Sign(signer frostfscrypto.Signer) error {
	var sig frostfscrypto.Signature

	err := sig.Calculate(signer, x.m.GetBody().StableMarshal(nil))
	if err != nil {
		return fmt.Errorf("calculate signature: %w", err)
	}

	var sigv2 refs.Signature
	sig.WriteToV2(&sigv2)

	x.m.SetSignature(&sigv2)

	return nil
}

// VerifySignature checks if GlobalTrust signature is presented and valid.
//
// Zero GlobalTrust fails the check.
//
// See also Sign.
func (x GlobalTrust) VerifySignature() bool {
	sigV2 := x.m.GetSignature()
	if sigV2 == nil {
		return false
	}

	var sig frostfscrypto.Signature

	return sig.ReadFromV2(*sigV2) == nil && sig.Verify(x.m.GetBody().StableMarshal(nil))
}

// Marshal encodes GlobalTrust into a binary format of the NeoFS API protocol
// (Protocol Buffers with direct field order).
//
// See also Unmarshal.
func (x GlobalTrust) Marshal() []byte {
	return x.m.StableMarshal(nil)
}

// Unmarshal decodes NeoFS API protocol binary format into the GlobalTrust
// (Protocol Buffers with direct field order). Returns an error describing
// a format violation.
//
// See also Marshal.
func (x *GlobalTrust) Unmarshal(data []byte) error {
	return x.m.Unmarshal(data)
}
