/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

package dhcplb

import (
	"github.com/facebookgo/ensure"
	"testing"
)

//SOLICIT message wrapped in Relay-Forw
var relayForwBytes = []byte{
	0x0c, // message type = Relay-forw
	0x00, // hop count = 0
	// link address
	0x24, 0x01, 0xdb, 0x00, 0x30, 0x10, 0xc0, 0xfa,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0a,
	// peer address
	0xfe, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x92, 0xe2, 0xba, 0xff, 0xfe, 0x76, 0x33, 0x44,
	0x00, 0x09, // relay message option type
	0x00, 0x28, // option length = 40

	// SOLICIT message as described below
	0x01, 0x00, 0xcd, 0x2e, 0x00, 0x08, 0x00, 0x02, 0xff, 0xff,
	0x00, 0x01, 0x00, 0x0a, 0x00, 0x03, 0x00, 0x01, 0x90, 0xe2,
	0xba, 0x76, 0x33, 0x44, 0x00, 0x03, 0x00, 0x0c, 0x01, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,

	// remote identifier option
	0x00, 0x25, 0x00, 0x0e, 0x00, 0x00, 0x00, 0x09, 0x00, 0x03,
	0x08, 0x00, 0x88, 0x5a, 0x92, 0xde, 0x8a, 0xbc,

	// interface id option
	0x00, 0x12, 0x00, 0x04, 0x09, 0x01, 0x08, 0xca}

//SOLICIT message wrapped in Relay-Forw
var relayForwBytesDuidUUID = []byte{
	0x0c, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0xfe, 0x80, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x26, 0x8a, 0x07, 0xff, 0xfe, 0x56,
	0xdc, 0xa4, 0x00, 0x12, 0x00, 0x06, 0x24, 0x8a,
	0x07, 0x56, 0xdc, 0xa4, 0x00, 0x09, 0x00, 0x5a,
	0x06, 0x7d, 0x9b, 0xca, 0x00, 0x01, 0x00, 0x12,
	0x00, 0x04, 0xb7, 0xfd, 0x0a, 0x8c, 0x1b, 0x14,
	0x10, 0xaa, 0xeb, 0x0a, 0x5b, 0x3f, 0xe8, 0x9d,
	0x0f, 0x56, 0x00, 0x06, 0x00, 0x0a, 0x00, 0x17,
	0x00, 0x18, 0x00, 0x17, 0x00, 0x18, 0x00, 0x01,
	0x00, 0x08, 0x00, 0x02, 0xff, 0xff, 0x00, 0x03,
	0x00, 0x28, 0x07, 0x56, 0xdc, 0xa4, 0x00, 0x00,
	0x0e, 0x10, 0x00, 0x00, 0x15, 0x18, 0x00, 0x05,
	0x00, 0x18, 0x26, 0x20, 0x01, 0x0d, 0xc0, 0x82,
	0x90, 0x63, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0xaf, 0xa0, 0x00, 0x00, 0x1c, 0x20, 0x00, 0x00,
	0x1d, 0x4c}

//ADVERTISE response to SOLICIT wrapped in Relay-Repl
var relayReplBytes = []byte{
	0x0d, // message type = relay-repl
	0x00, // hop count = 0
	// link address
	0x24, 0x01, 0xdb, 0x00, 0x30, 0x10, 0xc0, 0xfa,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0a,
	// peer address
	0xfe, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x92, 0xe2, 0xba, 0xff, 0xfe, 0x76, 0x33, 0x44,
	// interface-id option
	0x00, 0x12, 0x00, 0x04, 0x09, 0x01, 0x08, 0xca,
	0x00, 0x09, // relay message option type
	0x00, 0xba, // option length = 186
	0x02,             // message type = ADVERTISE
	0x00, 0xcd, 0x2e, // XID
	// client identifier option
	0x00, 0x01, 0x00, 0x0a, 0x00, 0x03, 0x00, 0x01, 0x90,
	0xe2, 0xba, 0x76, 0x33, 0x44,
	// server identifier option
	0x00, 0x02, 0x00, 0x0e, 0x00, 0x01, 0x00, 0x01, 0x1e,
	0xde, 0xbd, 0x1c, 0x00, 0x02, 0xc9, 0xca, 0x69, 0x7e,
	// identity association option
	0x00, 0x03, 0x00, 0x28, 0x01, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x6a, 0x24, 0x00, 0x00, 0x72, 0xbf, 0x00, 0x05,
	0x00, 0x18, 0x24, 0x01, 0xdb, 0x00, 0x30, 0x10, 0xc0,
	0xfa, 0xfa, 0xce, 0x00, 0x00, 0x00, 0x28, 0x00, 0x00,
	0x00, 0x00, 0xd4, 0x49, 0x00, 0x00, 0xe2, 0x59,
	// DNS recursive nameserver option
	0x00, 0x17, 0x00, 0x20, 0x24, 0x01, 0xdb, 0x00, 0xee,
	0xf0, 0x0a, 0x53, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x24, 0x01, 0xdb, 0x00, 0xee, 0xf0, 0x0b,
	0x53, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00,
	// domain search list option
	0x18, 0x00, 0x42, 0x02, 0x31, 0x32, 0x04, 0x6c, 0x6c,
	0x61, 0x31, 0x08, 0x66, 0x61, 0x63, 0x65, 0x62, 0x6f,
	0x6f, 0x6b, 0x03, 0x63, 0x6f, 0x6d, 0x00, 0x04, 0x6c,
	0x6c, 0x61, 0x31, 0x08, 0x66, 0x61, 0x63, 0x65, 0x62,
	0x6f, 0x6f, 0x6b, 0x03, 0x63, 0x6f, 0x6d, 0x00, 0x08,
	0x66, 0x61, 0x63, 0x65, 0x62, 0x6f, 0x6f, 0x6b, 0x03,
	0x63, 0x6f, 0x6d, 0x00, 0x05, 0x74, 0x66, 0x62, 0x6e,
	0x77, 0x03, 0x6e, 0x65, 0x74, 0x00}

func TestXidSolicit(t *testing.T) {
	//SOLICIT message copied from Wireshark as hex bytes
	//XID = 0x00cd2e
	bytes := []byte{0x01, //message type
		0x00, 0xcd, 0x2e, // XID
		0x00, 0x08, 0x00, 0x02, 0xff, 0xff, // elapsed time option
		// client identifier
		0x00, 0x01, 0x00, 0x0a, 0x00, 0x03, 0x00, 0x01,
		0x90, 0xe2, 0xba, 0x76, 0x33, 0x44,
		// identity association
		0x00, 0x03, 0x00, 0x0c, 0x01, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	packet := Packet6(bytes)
	xid, err := packet.XID()
	if err != nil {
		t.Fatalf("Failed to extract XId: %s", err)
	}
	// XID should be 0x00cd2e
	if xid != 0x00cd2e {
		t.Fatalf("Expected xid 0x%x but got 0x%x", 0x00cd2e, xid)
	}
}

func TestXidRelayForw(t *testing.T) {
	packet := Packet6(relayForwBytes)
	xid, err := packet.XID()
	if err != nil {
		t.Fatalf("Failed to extract XId: %s", err)
	}
	// XID should be 0x00cd2e
	if xid != 0x00cd2e {
		t.Fatalf("Expected xid 0x%x but got 0x%x", 0x00cd2e, xid)
	}
}

func TestXidRelayReply(t *testing.T) {
	packet := Packet6(relayReplBytes)
	xid, err := packet.XID()
	if err != nil {
		t.Fatalf("%s", err)
	}
	// XID should be 0x00cd2e
	if xid != 0x00cd2e {
		t.Fatalf("Expected xid 0x%x but got 0x%x", 0x00cd2e, xid)
	}
}

func TestDuid(t *testing.T) {
	packet := Packet6(relayForwBytes)
	duid, err := packet.Duid()
	if err != nil {
		t.Fatalf("Error extracting duid: %s", err)
	}
	expected := []byte{0x00, 0x03, 0x00, 0x01, 0x90, 0xe2, 0xba, 0x76, 0x33, 0x44}
	ensure.DeepEqual(t, duid, expected)
}

func TestDuidUUID(t *testing.T) {
	packet := Packet6(relayForwBytesDuidUUID)
	duid, err := packet.Duid()
	if err != nil {
		t.Fatalf("Error extracting duid: %s", err)
	}
	expected := []byte{0x00, 0x04, 0xb7, 0xfd, 0x0a, 0x8c, 0x1b, 0x14,
		0x10, 0xaa, 0xeb, 0x0a, 0x5b, 0x3f, 0xe8, 0x9d,
		0x0f, 0x56}
	ensure.DeepEqual(t, duid, expected)
	mac, errMac := packet.Mac()
	if errMac != nil {
		t.Fatalf("Error extracting mac from peer-address relayinfo: %s", errMac)
	}
	if FormatID(mac) != "24:8a:07:56:dc:a4" {
		t.Fatalf("Expected mac %s but got %s", "24:8a:07:56:de:b0", FormatID(mac))
	}
}

// basic sanity check that packet decodes the correct way after being
// encapsulated in a relay-forward
func TestEncapsulateSanity(t *testing.T) {
	packet := Packet6(relayForwBytes)
	startXid, err := packet.XID()
	if err != nil {
		t.Fatalf("Failed to extract XId: %s", err)
	}

	encapsulated := packet.Encapsulate(nil)
	// dhcp6message type should be SOLICIT
	msg, err := encapsulated.dhcp6message()
	msgType, _ := msg.Type()
	if msgType != Solicit {
		t.Fatalf("Expected type %s, got %s", Solicit, msgType)
	}
	// XID should be the same after encapsulating
	encXid, err := encapsulated.XID()
	if err != nil {
		t.Fatalf("Failed to extract XId from encapsulated message")
	}
	if startXid != encXid {
		t.Fatalf("Expected xid 0x%x but got 0x%x", startXid, encXid)
	}
}

func TestUnwindXid(t *testing.T) {
	packet := Packet6(relayReplBytes)
	msg, err := packet.dhcp6message()
	if err != nil {
		t.Fatalf("Failed to extract message")
	}
	// type should be ADVERTISE
	msgType, _ := msg.Type()
	if msgType != Advertise {
		t.Fatalf("Expected type %s but got %s", Advertise, msgType)
	}
	// XID should be 0x00cd2e
	xid, err := msg.XID()
	if err != nil {
		t.Fatalf("Failed to extract XId: %s", err)
	}
	if xid != 0x00cd2e {
		t.Fatalf("Expected xid 0x%x but got 0x%x", 0x00cd2e, xid)
	}
}

// Test_MalformedPacket tests to make sure the parser can properly handle a packet
// with an option that specifies its length such that it ends out-of-bounds
func TestMalformedPacket(t *testing.T) {
	// copied as bytes from Wireshark, then modified the RelayMessage option length
	bytes := []byte{
		0x0c, 0x00, 0x24, 0x01, 0xdb, 0x00, 0x30, 0x10, 0xb0, 0x8a, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x0a, 0xfe, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x0b, 0xab, 0xff, 0xfe, 0x8a,
		0x6d, 0xf2, 0x00, 0x09, 0x00, 0x50 /*was 0x32*/, 0x01, 0x8d, 0x3e, 0x24, 0x00, 0x01, 0x00, 0x0e, 0x00, 0x01,
		0x00, 0x01, 0x0c, 0x71, 0x3d, 0x0e, 0x00, 0x0b, 0xab, 0x8a, 0x6d, 0xf2, 0x00, 0x08, 0x00, 0x02,
		0x00, 0x00, 0x00, 0x03, 0x00, 0x0c, 0xee, 0xbf, 0xfb, 0x6e, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0x00, 0x06, 0x00, 0x02, 0x00, 0x17, 0x00, 0x25, 0x00, 0x0e, 0x00, 0x00, 0x00, 0x09,
		0x00, 0x03, 0x08, 0x00, 0xf0, 0x7f, 0x06, 0xd6, 0x4c, 0x3c, 0x00, 0x12, 0x00, 0x04, 0x09, 0x01,
		0x08, 0x5a,
	}
	packet := Packet6(bytes)
	_, err := packet.dhcp6message()
	if err == nil {
		t.Fatalf("Should be unable to extract dhcp6message, but did not fail")
	}
}

func TestStackedRelayInfo(t *testing.T) {
	// taken from a production tcpdump
	bytes := []byte{
		0x0c, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x24, 0x01, 0xdb, 0x00, 0x01, 0x11,
		0x70, 0x32, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x0a, 0x00, 0x09, 0x00, 0x60, 0x0c, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0xfe, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x26, 0x8a, 0x07, 0xff, 0xfe, 0x87, 0xf9, 0x7f,
		0x00, 0x12, 0x00, 0x06, 0x24, 0x8a, 0x07, 0x87,
		0xf9, 0x7f, 0x00, 0x09, 0x00, 0x30, 0x01, 0x50,
		0xb4, 0x93, 0x00, 0x01, 0x00, 0x0a, 0x00, 0x03,
		0x00, 0x01, 0x24, 0x8a, 0x07, 0x87, 0xf9, 0x7f,
		0x00, 0x03, 0x00, 0x0c, 0x00, 0x00, 0x2b, 0x67,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x08, 0x00, 0x02, 0xff, 0xff, 0x00, 0x06,
		0x00, 0x04, 0x00, 0x17, 0x00, 0x18,
	}
	packet := Packet6(bytes)
	msg, err := packet.dhcp6message()
	if err != nil {
		t.Fatalf("Failed to extract message")
	}
	// type should be SOLICIT
	msgType, _ := msg.Type()
	if msgType != Solicit {
		t.Fatalf("Expected type %s but got %s", Solicit, msgType)
	}
	// XID should be 0x00cd2e
	xid, err := msg.XID()
	if err != nil {
		t.Fatalf("Failed to extract XId: %s", err)
	}
	if xid != 0x50b493 {
		t.Fatalf("Expected xid 0x%x but got 0x%x", 0x50b493, xid)
	}
	// Hops should be 2
	hops, err := packet.Hops()
	if err != nil {
		t.Fatalf("Failed to extract XId: %s", err)
	}
	if hops != 1 {
		t.Fatalf("Expected 1 hop, got %d", hops)
	}
	addr, err := packet.GetInnerMostPeerAddr()
	if err != nil {
		t.Fatalf("Failed running GetInnerMostPeerAddr: %s", err)
	}
	if addr.String() != "fe80::268a:7ff:fe87:f97f" {
		t.Fatalf(
			"Expected %s but got %s", "fe80::268a:7ff:fe87:f97f", addr.String())
	}
}
