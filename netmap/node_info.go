package netmap

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/TrueCloudLab/frostfs-api-go/v2/netmap"
	"github.com/TrueCloudLab/frostfs-api-go/v2/refs"
	frostfscrypto "github.com/TrueCloudLab/frostfs-sdk-go/crypto"
	subnetid "github.com/TrueCloudLab/frostfs-sdk-go/subnet/id"
	"github.com/TrueCloudLab/hrw"
)

// NodeInfo groups information about FrostFS storage node which is reflected
// in the FrostFS network map. Storage nodes advertise this information when
// registering with the FrostFS network. After successful registration, information
// about the nodes is available to all network participants to work with the network
// map (mainly to comply with container storage policies).
//
// NodeInfo is mutually compatible with github.com/TrueCloudLab/frostfs-api-go/v2/netmap.NodeInfo
// message. See ReadFromV2 / WriteToV2 methods.
//
// Instances can be created using built-in var declaration.
type NodeInfo struct {
	m netmap.NodeInfo
}

// reads NodeInfo from netmap.NodeInfo message. If checkFieldPresence is set,
// returns an error on absence of any protocol-required field. Verifies format of any
// presented field according to FrostFS API V2 protocol.
func (x *NodeInfo) readFromV2(m netmap.NodeInfo, checkFieldPresence bool) error {
	var err error

	binPublicKey := m.GetPublicKey()
	if checkFieldPresence && len(binPublicKey) == 0 {
		return errors.New("missing public key")
	}

	if checkFieldPresence && m.NumberOfAddresses() <= 0 {
		return errors.New("missing network endpoints")
	}

	attributes := m.GetAttributes()
	mAttr := make(map[string]struct{}, len(attributes))
	for i := range attributes {
		key := attributes[i].GetKey()
		if key == "" {
			return fmt.Errorf("empty key of the attribute #%d", i)
		} else if _, ok := mAttr[key]; ok {
			return fmt.Errorf("duplicated attbiuted %s", key)
		}

		const subnetPrefix = "__NEOFS__SUBNET_"

		switch {
		case key == attrCapacity:
			_, err = strconv.ParseUint(attributes[i].GetValue(), 10, 64)
			if err != nil {
				return fmt.Errorf("invalid %s attribute: %w", attrCapacity, err)
			}
		case key == attrPrice:
			var err error
			_, err = strconv.ParseUint(attributes[i].GetValue(), 10, 64)
			if err != nil {
				return fmt.Errorf("invalid %s attribute: %w", attrPrice, err)
			}
		case strings.HasPrefix(key, subnetPrefix):
			var id subnetid.ID

			err = id.DecodeString(strings.TrimPrefix(key, subnetPrefix))
			if err != nil {
				return fmt.Errorf("invalid key to the subnet attribute %s: %w", key, err)
			}

			if val := attributes[i].GetValue(); val != "True" && val != "False" {
				return fmt.Errorf("invalid value of the subnet attribute %s: %w", val, err)
			}
		default:
			if attributes[i].GetValue() == "" {
				return fmt.Errorf("empty value of the attribute %s", key)
			}
		}
	}

	x.m = m

	return nil
}

// ReadFromV2 reads NodeInfo from the netmap.NodeInfo message. Checks if the
// message conforms to FrostFS API V2 protocol.
//
// See also WriteToV2.
func (x *NodeInfo) ReadFromV2(m netmap.NodeInfo) error {
	return x.readFromV2(m, true)
}

// WriteToV2 writes NodeInfo to the netmap.NodeInfo message. The message MUST NOT
// be nil.
//
// See also ReadFromV2.
func (x NodeInfo) WriteToV2(m *netmap.NodeInfo) {
	*m = x.m
}

// Marshal encodes NodeInfo into a binary format of the FrostFS API protocol
// (Protocol Buffers with direct field order).
//
// See also Unmarshal.
func (x NodeInfo) Marshal() []byte {
	var m netmap.NodeInfo
	x.WriteToV2(&m)

	return m.StableMarshal(nil)
}

// Unmarshal decodes FrostFS API protocol binary format into the NodeInfo
// (Protocol Buffers with direct field order). Returns an error describing
// a format violation.
//
// See also Marshal.
func (x *NodeInfo) Unmarshal(data []byte) error {
	var m netmap.NodeInfo

	err := m.Unmarshal(data)
	if err != nil {
		return err
	}

	return x.readFromV2(m, false)
}

// MarshalJSON encodes NodeInfo into a JSON format of the FrostFS API protocol
// (Protocol Buffers JSON).
//
// See also UnmarshalJSON.
func (x NodeInfo) MarshalJSON() ([]byte, error) {
	var m netmap.NodeInfo
	x.WriteToV2(&m)

	return m.MarshalJSON()
}

// UnmarshalJSON decodes FrostFS API protocol JSON format into the NodeInfo
// (Protocol Buffers JSON). Returns an error describing a format violation.
//
// See also MarshalJSON.
func (x *NodeInfo) UnmarshalJSON(data []byte) error {
	var m netmap.NodeInfo

	err := m.UnmarshalJSON(data)
	if err != nil {
		return err
	}

	return x.readFromV2(m, false)
}

// SetPublicKey sets binary-encoded public key bound to the node. The key
// authenticates the storage node, so it MUST be unique within the network.
//
// Argument MUST NOT be mutated, make a copy first.
//
// See also PublicKey.
func (x *NodeInfo) SetPublicKey(key []byte) {
	x.m.SetPublicKey(key)
}

// PublicKey returns value set using SetPublicKey.
//
// Zero NodeInfo has no public key, which is incorrect according to
// FrostFS system requirements.
//
// Return value MUST not be mutated, make a copy first.
func (x NodeInfo) PublicKey() []byte {
	return x.m.GetPublicKey()
}

// StringifyPublicKey returns HEX representation of PublicKey.
func StringifyPublicKey(node NodeInfo) string {
	return frostfscrypto.StringifyKeyBinary(node.PublicKey())
}

// SetNetworkEndpoints sets list to the announced node's network endpoints.
// Node MUSt have at least one announced endpoint. List MUST be unique.
// Endpoints are used for communication with the storage node within FrostFS
// network. It is expected that node serves storage node services on these
// endpoints (it also adds a wait on their network availability).
//
// Argument MUST NOT be mutated, make a copy first.
//
// See also IterateNetworkEndpoints.
func (x *NodeInfo) SetNetworkEndpoints(v ...string) {
	x.m.SetAddresses(v...)
}

// NumberOfNetworkEndpoints returns number of network endpoints announced by the node.
//
// See also SetNetworkEndpoints.
func (x NodeInfo) NumberOfNetworkEndpoints() int {
	return x.m.NumberOfAddresses()
}

// IterateNetworkEndpoints iterates over network endpoints announced by the
// node and pass them into f. Breaks iteration on f's true return. Handler
// MUST NOT be nil.
//
// Zero NodeInfo contains no endpoints which is incorrect according to
// FrostFS system requirements.
//
// See also SetNetworkEndpoints.
func (x NodeInfo) IterateNetworkEndpoints(f func(string) bool) {
	x.m.IterateAddresses(f)
}

// IterateNetworkEndpoints is an extra-sugared function over IterateNetworkEndpoints
// method which allows to unconditionally iterate over all node's network endpoints.
func IterateNetworkEndpoints(node NodeInfo, f func(string)) {
	node.IterateNetworkEndpoints(func(addr string) bool {
		f(addr)
		return false
	})
}

// assert NodeInfo type provides hrw.Hasher required for HRW sorting.
var _ hrw.Hasher = NodeInfo{}

// Hash implements hrw.Hasher interface.
//
// Hash is needed to support weighted HRW therefore sort function sorts nodes
// based on their public key. Hash isn't expected to be used directly.
func (x NodeInfo) Hash() uint64 {
	return hrw.Hash(x.m.GetPublicKey())
}

// less declares "less than" comparison between two NodeInfo instances:
// x1 is less than x2 if it has less Hash().
//
// Method is needed for internal placement needs.
func less(x1, x2 NodeInfo) bool {
	return x1.Hash() < x2.Hash()
}

func (x *NodeInfo) setNumericAttribute(key string, num uint64) {
	x.SetAttribute(key, strconv.FormatUint(num, 10))
}

// SetPrice sets the storage cost declared by the node. By default, zero
// price is announced.
func (x *NodeInfo) SetPrice(price uint64) {
	x.setNumericAttribute(attrPrice, price)
}

// Price returns price set using SetPrice.
//
// Zero NodeInfo has zero price.
func (x NodeInfo) Price() uint64 {
	val := x.Attribute(attrPrice)
	if val == "" {
		return 0
	}

	price, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("unexpected price parsing error %s: %v", val, err))
	}

	return price
}

// SetCapacity sets the storage capacity declared by the node. By default, zero
// capacity is announced.
func (x *NodeInfo) SetCapacity(capacity uint64) {
	x.setNumericAttribute(attrCapacity, capacity)
}

// capacity returns capacity set using SetCapacity.
//
// Zero NodeInfo has zero capacity.
func (x NodeInfo) capacity() uint64 {
	val := x.Attribute(attrCapacity)
	if val == "" {
		return 0
	}

	capacity, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("unexpected capacity parsing error %s: %v", val, err))
	}

	return capacity
}

const attrUNLOCODE = "UN-LOCODE"

// SetLOCODE specifies node's geographic location in UN/LOCODE format. Each
// storage node MUST declare it for entrance to the FrostFS network. Node MAY
// declare the code of the nearest location as needed, for example, when it is
// impossible to unambiguously attribute the node to any location from UN/LOCODE
// database.
//
// See also LOCODE.
func (x *NodeInfo) SetLOCODE(locode string) {
	x.SetAttribute(attrUNLOCODE, locode)
}

// LOCODE returns node's location code set using SetLOCODE.
//
// Zero NodeInfo has empty location code which is invalid according to
// FrostFS API system requirement.
func (x NodeInfo) LOCODE() string {
	return x.Attribute(attrUNLOCODE)
}

// SetCountryCode sets code of the country in ISO 3166-1_alpha-2 to which
// storage node belongs (or the closest one).
//
// SetCountryCode is intended only for processing the network registration
// request by the Inner Ring. Other parties SHOULD NOT use it.
func (x *NodeInfo) SetCountryCode(countryCode string) {
	x.SetAttribute("CountryCode", countryCode)
}

// SetCountryName sets short name of the country in ISO-3166 format to which
// storage node belongs (or the closest one).
//
// SetCountryName is intended only for processing the network registration
// request by the Inner Ring. Other parties SHOULD NOT use it.
func (x *NodeInfo) SetCountryName(country string) {
	x.SetAttribute("Country", country)
}

// SetLocationName sets storage node's location name from "NameWoDiacritics"
// column in the UN/LOCODE record corresponding to the specified LOCODE.
//
// SetLocationName is intended only for processing the network registration
// request by the Inner Ring. Other parties SHOULD NOT use it.
func (x *NodeInfo) SetLocationName(location string) {
	x.SetAttribute("Location", location)
}

// SetSubdivisionCode sets storage node's subdivision code from "SubDiv" column in
// the UN/LOCODE record corresponding to the specified LOCODE.
//
// SetSubdivisionCode is intended only for processing the network registration
// request by the Inner Ring. Other parties SHOULD NOT use it.
func (x *NodeInfo) SetSubdivisionCode(subDiv string) {
	x.SetAttribute("SubDivCode", subDiv)
}

// SetSubdivisionName sets storage node's subdivision name in ISO 3166-2 format.
//
// SetSubdivisionName is intended only for processing the network registration
// request by the Inner Ring. Other parties SHOULD NOT use it.
func (x *NodeInfo) SetSubdivisionName(subDiv string) {
	x.SetAttribute("SubDiv", subDiv)
}

// SetContinentName sets name of the storage node's continent from
// Seven-Continent model.
//
// SetContinentName is intended only for processing the network registration
// request by the Inner Ring. Other parties SHOULD NOT use it.
func (x *NodeInfo) SetContinentName(continent string) {
	x.SetAttribute("Continent", continent)
}

// Enumeration of well-known attributes.
const (
	// attrPrice is a key to the node attribute that indicates the
	// price in GAS tokens for storing one GB of data during one Epoch.
	attrPrice = "Price"

	// attrCapacity is a key to the node attribute that indicates the
	// total available disk space in Gigabytes.
	attrCapacity = "Capacity"

	// attrExternalAddr is a key for the attribute storing node external addresses.
	attrExternalAddr = "ExternalAddr"
	// sepExternalAddr is a separator for multi-value ExternalAddr attribute.
	sepExternalAddr = ","
)

// SetExternalAddresses sets multi-addresses to use
// to connect to this node from outside.
//
// Panics if addr is an empty list.
func (x *NodeInfo) SetExternalAddresses(addr ...string) {
	x.SetAttribute(attrExternalAddr, strings.Join(addr, sepExternalAddr))
}

// ExternalAddresses returns list of multi-addresses to use
// to connect to this node from outside.
func (x NodeInfo) ExternalAddresses() []string {
	a := x.Attribute(attrExternalAddr)
	if len(a) == 0 {
		return nil
	}

	return strings.Split(a, sepExternalAddr)
}

// NumberOfAttributes returns number of attributes announced by the node.
//
// See also SetAttribute.
func (x NodeInfo) NumberOfAttributes() int {
	return len(x.m.GetAttributes())
}

// IterateAttributes iterates over all node attributes and passes the into f.
// Handler MUST NOT be nil.
func (x NodeInfo) IterateAttributes(f func(key, value string)) {
	a := x.m.GetAttributes()
	for i := range a {
		f(a[i].GetKey(), a[i].GetValue())
	}
}

// SetAttribute sets value of the node attribute value by the given key.
// Both key and value MUST NOT be empty.
func (x *NodeInfo) SetAttribute(key, value string) {
	if key == "" {
		panic("empty key in SetAttribute")
	} else if value == "" {
		panic("empty value in SetAttribute")
	}

	a := x.m.GetAttributes()
	for i := range a {
		if a[i].GetKey() == key {
			a[i].SetValue(value)
			return
		}
	}

	a = append(a, netmap.Attribute{})
	a[len(a)-1].SetKey(key)
	a[len(a)-1].SetValue(value)

	x.m.SetAttributes(a)
}

// Attribute returns value of the node attribute set using SetAttribute by the
// given key. Returns empty string if attribute is missing.
func (x NodeInfo) Attribute(key string) string {
	a := x.m.GetAttributes()
	for i := range a {
		if a[i].GetKey() == key {
			return a[i].GetValue()
		}
	}

	return ""
}

// SortAttributes sorts node attributes set using SetAttribute lexicographically.
// The method is only needed to make NodeInfo consistent, e.g. for signing.
func (x *NodeInfo) SortAttributes() {
	as := x.m.GetAttributes()
	if len(as) == 0 {
		return
	}

	sort.Slice(as, func(i, j int) bool {
		switch strings.Compare(as[i].GetKey(), as[j].GetKey()) {
		case -1:
			return true
		case 1:
			return false
		default:
			return as[i].GetValue() < as[j].GetValue()
		}
	})

	x.m.SetAttributes(as)
}

// EnterSubnet writes storage node's intention to enter the given subnet.
//
// Zero NodeInfo belongs to zero subnet.
func (x *NodeInfo) EnterSubnet(id subnetid.ID) {
	x.changeSubnet(id, true)
}

// ExitSubnet writes storage node's intention to exit the given subnet.
func (x *NodeInfo) ExitSubnet(id subnetid.ID) {
	x.changeSubnet(id, false)
}

func (x *NodeInfo) changeSubnet(id subnetid.ID, isMember bool) {
	var (
		idv2 refs.SubnetID
		info netmap.NodeSubnetInfo
	)

	id.WriteToV2(&idv2)

	info.SetID(&idv2)
	info.SetEntryFlag(isMember)

	netmap.WriteSubnetInfo(&x.m, info)
}

// ErrRemoveSubnet is returned when a node needs to leave the subnet.
var ErrRemoveSubnet = netmap.ErrRemoveSubnet

// IterateSubnets iterates over all subnets the node belongs to and passes the IDs to f.
// Handler MUST NOT be nil.
//
// If f returns ErrRemoveSubnet, then removes subnet entry. Note that this leads to an
// instant mutation of NodeInfo. Breaks on any other non-nil error and returns it.
//
// Returns an error if subnet incorrectly enabled/disabled.
// Returns an error if the node is not included to any subnet by the end of the loop.
//
// See also EnterSubnet, ExitSubnet.
func (x NodeInfo) IterateSubnets(f func(subnetid.ID) error) error {
	var id subnetid.ID

	return netmap.IterateSubnets(&x.m, func(idv2 refs.SubnetID) error {
		err := id.ReadFromV2(idv2)
		if err != nil {
			return fmt.Errorf("invalid subnet: %w", err)
		}

		err = f(id)
		if errors.Is(err, ErrRemoveSubnet) {
			return netmap.ErrRemoveSubnet
		}

		return err
	})
}

var errAbortSubnetIter = errors.New("abort subnet iterator")

// BelongsToSubnet is a helper function over the IterateSubnets method which
// checks whether a node belongs to a subnet.
//
// Zero NodeInfo belongs to zero subnet only.
func BelongsToSubnet(node NodeInfo, id subnetid.ID) bool {
	err := node.IterateSubnets(func(id_ subnetid.ID) error {
		if id.Equals(id_) {
			return errAbortSubnetIter
		}

		return nil
	})

	return errors.Is(err, errAbortSubnetIter)
}

// SetOffline sets the state of the node to "offline". When a node updates
// information about itself in the network map, this action is interpreted as
// an intention to leave the network.
func (x *NodeInfo) SetOffline() {
	x.m.SetState(netmap.Offline)
}

// IsOffline checks if the node is in the "offline" state.
//
// Zero NodeInfo has undefined state which is not offline (note that it does not
// mean online).
//
// See also SetOffline.
func (x NodeInfo) IsOffline() bool {
	return x.m.GetState() == netmap.Offline
}

// SetOnline sets the state of the node to "online". When a node updates
// information about itself in the network map, this
// action is interpreted as an intention to enter the network.
//
// See also IsOnline.
func (x *NodeInfo) SetOnline() {
	x.m.SetState(netmap.Online)
}

// IsOnline checks if the node is in the "online" state.
//
// Zero NodeInfo has undefined state which is not online (note that it does not
// mean offline).
//
// See also SetOnline.
func (x NodeInfo) IsOnline() bool {
	return x.m.GetState() == netmap.Online
}

// SetMaintenance sets the state of the node to "maintenance". When a node updates
// information about itself in the network map, this
// state declares temporal unavailability for a node.
//
// See also IsMaintenance.
func (x *NodeInfo) SetMaintenance() {
	x.m.SetState(netmap.Maintenance)
}

// IsMaintenance checks if the node is in the "maintenance" state.
//
// Zero NodeInfo has undefined state.
//
// See also SetMaintenance.
func (x NodeInfo) IsMaintenance() bool {
	return x.m.GetState() == netmap.Maintenance
}
