// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

//go:build !js
// +build !js

package stun

import (
	"net"
	"testing"
)

func BenchmarkFingerprint_AddTo(b *testing.B) {
	b.ReportAllocs()
	msg := new(Message)
	s := NewSoftware("software")
	addr := &XORMappedAddress{
		IP: net.IPv4(213, 1, 223, 5),
	}
	addAttr(b, msg, addr)
	addAttr(b, msg, s)
	b.SetBytes(int64(len(msg.Raw)))
	for i := 0; i < b.N; i++ {
		Fingerprint.AddTo(msg) //nolint:errcheck,gosec
		msg.WriteLength()
		msg.Length -= attributeHeaderSize + fingerprintSize
		msg.Raw = msg.Raw[:msg.Length+messageHeaderSize]
		msg.Attributes = msg.Attributes[:len(msg.Attributes)-1]
	}
}

func TestFingerprint_Check(t *testing.T) {
	m := new(Message)
	addAttr(t, m, NewSoftware("software"))
	m.WriteHeader()
	Fingerprint.AddTo(m) //nolint:errcheck,gosec
	m.WriteHeader()
	if err := Fingerprint.Check(m); err != nil {
		t.Error(err)
	}
	m.Raw[3]++
	if err := Fingerprint.Check(m); err == nil {
		t.Error("should error")
	}
}

func TestFingerprint_CheckBad(t *testing.T) {
	m := new(Message)
	addAttr(t, m, NewSoftware("software"))
	m.WriteHeader()
	if err := Fingerprint.Check(m); err == nil {
		t.Error("should error")
	}
	m.Add(AttrFingerprint, []byte{1, 2, 3})
	if !IsAttrSizeInvalid(Fingerprint.Check(m)) {
		t.Error("IsAttrSizeInvalid should be true")
	}
}

func BenchmarkFingerprint_Check(b *testing.B) {
	b.ReportAllocs()
	m := new(Message)
	s := NewSoftware("software")
	addr := &XORMappedAddress{
		IP: net.IPv4(213, 1, 223, 5),
	}
	addAttr(b, m, addr)
	addAttr(b, m, s)
	m.WriteHeader()
	Fingerprint.AddTo(m) //nolint:errcheck,gosec
	m.WriteHeader()
	b.SetBytes(int64(len(m.Raw)))
	for i := 0; i < b.N; i++ {
		if err := Fingerprint.Check(m); err != nil {
			b.Fatal(err)
		}
	}
}
