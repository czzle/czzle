package vd

import (
	"testing"
)

func TestIPv4(t *testing.T) {
	err := IPv4("88.222.172.43").Validate(nil, "")
	if err != nil {
		t.Fatal(err)
	}
}

func TestIPv6(t *testing.T) {
	err := IPv6("2001:0db8:85a3:0000:0000:8a2e:0370:7334").Validate(nil, "")
	if err != nil {
		t.Fatal(err)
	}
}
