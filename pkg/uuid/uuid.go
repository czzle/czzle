package uuid

import (
	"crypto/rand"
	"database/sql/driver"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
)

type UUID [16]byte

// Null create a new Null UUID
func Null() UUID {
	return UUID{}
}

const rxStr = "^[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[89ab][a-f0-9]{3}-[a-f0-9]{12}$"

var rx = regexp.MustCompile(rxStr)

// Version returns version number of uuid
func (u UUID) Version() int {
	return int(binary.BigEndian.Uint16(u[6:8]) >> 12)
}

// String returns string value of UUID
func (u UUID) String() string {
	if u.IsNull() {
		return ""
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
}

func (u UUID) MarshalJSON() ([]byte, error) {
	if u.IsNull() {
		return []byte("null"), nil
	}
	return []byte(
		fmt.Sprintf("\"%s\"", u),
	), nil
}

func (u *UUID) UnmarshalJSON(src []byte) error {
	str := string(src)
	if str == "null" {
		*u = Null()
		return nil
	}
	str = strings.Trim(str, "\"")
	*u = FromString(str)
	return nil
}

func (u UUID) Value() (driver.Value, error) {
	if u.IsNull() {
		return nil, nil
	}
	return []byte(u.String()), nil
}

func (u *UUID) Scan(src interface{}) error {
	if src == nil {
		*u = Null()
		return nil
	}
	bytes := src.([]byte)
	str := string(bytes)
	*u = FromString(str)
	return nil
}

func (u UUID) Equals(other interface{}) bool {
	uu, ok := other.(UUID)
	if !ok {
		return false
	}
	for i := 0; i < 16; i++ {
		if u[i] != uu[i] {
			return false
		}
	}
	return true
}

// IsNull checks if UUID is null
func (u UUID) IsNull() bool {
	for _, b := range u {
		if b != 0 {
			return false
		}
	}
	return true
}

// New generates new UUIDv4
func New() UUID {
	buf := make([]byte, 16)
	rand.Read(buf)
	buf[6] = (buf[6] & 0x0f) | 0x40
	var uuid UUID
	copy(uuid[:], buf[:])
	uuid[8] = (uuid[8] & 0x3f) | 0x80
	return uuid
}

// FromString converts string to UUID
func FromString(str string) UUID {
	if str == "" || !rx.MatchString(str) {
		return Null()
	}
	dst := UUID{}
	str = strings.Replace(str, "-", "", -1)
	src := make([]byte, 16)
	hex.Decode(src, []byte(str))
	copy(dst[0:4], src[0:4])
	copy(dst[4:6], src[4:6])
	copy(dst[6:8], src[6:8])
	copy(dst[8:10], src[8:10])
	copy(dst[10:], src[10:])
	return dst
}
