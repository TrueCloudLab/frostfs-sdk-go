package version

import (
	"fmt"

	"github.com/TrueCloudLab/frostfs-api-go/v2/refs"
)

// Version represents revision number in SemVer scheme.
//
// Version is mutually compatible with github.com/TrueCloudLab/frostfs-api-go/v2/refs.Version
// message. See ReadFromV2 / WriteToV2 methods.
//
// Instances can be created using built-in var declaration.
//
// Note that direct typecast is not safe and may result in loss of compatibility:
//
//	_ = Version(refs.Version{}) // not recommended
type Version refs.Version

const sdkMjr, sdkMnr = 2, 13

// Current returns Version instance that initialized to the
// latest supported FrostFS API revision number in SDK.
func Current() (v Version) {
	v.SetMajor(sdkMjr)
	v.SetMinor(sdkMnr)
	return v
}

// Major returns major number of the revision.
func (v *Version) Major() uint32 {
	return (*refs.Version)(v).GetMajor()
}

// SetMajor sets major number of the revision.
func (v *Version) SetMajor(val uint32) {
	(*refs.Version)(v).SetMajor(val)
}

// Minor returns minor number of the revision.
func (v *Version) Minor() uint32 {
	return (*refs.Version)(v).GetMinor()
}

// SetMinor sets minor number of the revision.
func (v *Version) SetMinor(val uint32) {
	(*refs.Version)(v).SetMinor(val)
}

// WriteToV2 writes Version to the refs.Version message.
// The message must not be nil.
//
// See also ReadFromV2.
func (v Version) WriteToV2(m *refs.Version) {
	*m = (refs.Version)(v)
}

// ReadFromV2 reads Version from the refs.Version message. Checks if the message
// conforms to FrostFS API V2 protocol.
//
// See also WriteToV2.
func (v *Version) ReadFromV2(m refs.Version) error {
	*v = Version(m)
	return nil
}

// String implements fmt.Stringer.
//
// String is designed to be human-readable, and its format MAY differ between
// SDK versions.
func (v Version) String() string {
	return EncodeToString(v)
}

// EncodeToString encodes version according to format from specification:
// semver formatted value without patch and with v prefix, e.g. 'v2.1'.
func EncodeToString(v Version) string {
	return fmt.Sprintf("v%d.%d", v.Major(), v.Minor())
}

// Equal returns true if versions are identical.
func (v Version) Equal(v2 Version) bool {
	return v.Major() == v2.Major() &&
		v.Minor() == v2.Minor()
}
