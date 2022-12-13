package object

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/TrueCloudLab/frostfs-api-go/v2/object"
	"github.com/TrueCloudLab/frostfs-sdk-go/checksum"
	frostfscrypto "github.com/TrueCloudLab/frostfs-sdk-go/crypto"
	oid "github.com/TrueCloudLab/frostfs-sdk-go/object/id"
)

var (
	errCheckSumMismatch = errors.New("payload checksum mismatch")
	errCheckSumNotSet   = errors.New("payload checksum is not set")
	errIncorrectID      = errors.New("incorrect object identifier")
)

// CalculatePayloadChecksum calculates and returns checksum of
// object payload bytes.
func CalculatePayloadChecksum(payload []byte) checksum.Checksum {
	var res checksum.Checksum
	checksum.Calculate(&res, checksum.SHA256, payload)

	return res
}

// CalculateAndSetPayloadChecksum calculates checksum of current
// object payload and writes it to the object.
func CalculateAndSetPayloadChecksum(obj *Object) {
	obj.SetPayloadChecksum(
		CalculatePayloadChecksum(obj.Payload()),
	)
}

// VerifyPayloadChecksum checks if payload checksum in the object
// corresponds to its payload.
func VerifyPayloadChecksum(obj *Object) error {
	actual := CalculatePayloadChecksum(obj.Payload())

	cs, set := obj.PayloadChecksum()
	if !set {
		return errCheckSumNotSet
	}

	if !bytes.Equal(cs.Value(), actual.Value()) {
		return errCheckSumMismatch
	}

	return nil
}

// CalculateID calculates identifier for the object.
func CalculateID(obj *Object) (oid.ID, error) {
	var id oid.ID
	id.SetSHA256(sha256.Sum256(obj.ToV2().GetHeader().StableMarshal(nil)))

	return id, nil
}

// CalculateAndSetID calculates identifier for the object
// and writes the result to it.
func CalculateAndSetID(obj *Object) error {
	id, err := CalculateID(obj)
	if err != nil {
		return err
	}

	obj.SetID(id)

	return nil
}

// VerifyID checks if identifier in the object corresponds to
// its structure.
func VerifyID(obj *Object) error {
	id, err := CalculateID(obj)
	if err != nil {
		return err
	}

	oID, set := obj.ID()
	if !set {
		return errOIDNotSet
	}

	if !id.Equals(oID) {
		return errIncorrectID
	}

	return nil
}

// CalculateAndSetSignature signs id with provided key and sets that signature to
// the object.
func CalculateAndSetSignature(key ecdsa.PrivateKey, obj *Object) error {
	oID, set := obj.ID()
	if !set {
		return errOIDNotSet
	}

	sig, err := oID.CalculateIDSignature(key)
	if err != nil {
		return err
	}

	obj.SetSignature(&sig)

	return nil
}

// VerifyIDSignature verifies object ID signature.
func (o *Object) VerifyIDSignature() bool {
	m := (*object.Object)(o)

	sigV2 := m.GetSignature()
	if sigV2 == nil {
		return false
	}

	idV2 := m.GetObjectID()
	if idV2 == nil {
		return false
	}

	var sig frostfscrypto.Signature

	return sig.ReadFromV2(*sigV2) == nil && sig.Verify(idV2.StableMarshal(nil))
}

// SetIDWithSignature sets object identifier and signature.
func SetIDWithSignature(key ecdsa.PrivateKey, obj *Object) error {
	if err := CalculateAndSetID(obj); err != nil {
		return fmt.Errorf("could not set identifier: %w", err)
	}

	if err := CalculateAndSetSignature(key, obj); err != nil {
		return fmt.Errorf("could not set signature: %w", err)
	}

	return nil
}

// SetVerificationFields calculates and sets all verification fields of the object.
func SetVerificationFields(key ecdsa.PrivateKey, obj *Object) error {
	CalculateAndSetPayloadChecksum(obj)

	return SetIDWithSignature(key, obj)
}

// CheckVerificationFields checks all verification fields of the object.
func CheckVerificationFields(obj *Object) error {
	if err := CheckHeaderVerificationFields(obj); err != nil {
		return fmt.Errorf("invalid header structure: %w", err)
	}

	if err := VerifyPayloadChecksum(obj); err != nil {
		return fmt.Errorf("invalid payload checksum: %w", err)
	}

	return nil
}

var errInvalidSignature = errors.New("invalid signature")

// CheckHeaderVerificationFields checks all verification fields except payload.
func CheckHeaderVerificationFields(obj *Object) error {
	if !obj.VerifyIDSignature() {
		return errInvalidSignature
	}

	if err := VerifyID(obj); err != nil {
		return fmt.Errorf("invalid identifier: %w", err)
	}

	return nil
}
