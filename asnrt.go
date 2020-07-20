package asnrt

import (
	"reflect"
)

// Constants

const (
	// Encoding Rules
	BASIC_ENCODING_RULES            = 0
	CANONICAL_ENCODING_RULES        = 1
	DISTINGUISHED_ENCODING_RULES    = 2
	UNALIGNED_PACKED_ENCODING_RULES = 3
	ALIGNED_PACKED_ENCODING_RULES   = 4

	OCTET_ENCODING_RULES = 8
)

// Type definitions
type Any interface{}

type Null struct{}

var NULL Null = Null{}

type Enum interface {
	Name() string
	Ordinal() int
}

type BitString struct {
	Bytes      []byte // bits packed into bytes.
	UnusedBits uint8  // unused bits in tailing byte.
}

type OctetString []byte

type ObjectIdentifier []int

type Buffer interface {
	EncodingRules() byte
	Array() []byte
	Flip()
}

type AsnType interface {
	TypeId() uint8
	Tag(value reflect.Value) int
	MatchTag(tag int) bool
}

type AsnModule interface {
	GetType(id int) AsnType
	GetValue(id int, v Any)
	GetValueSet(id int, vs Any)
	GetObject(id int, o Any)
	GetObjectSet(id int, os Any)
}

// API Methods

type Runtime struct {
	Allocate func(length int, encodingRules byte) (Buffer, error)
	Wrap     func(array []byte, encodingRules byte) (Buffer, error)
	Encode   func(b Buffer, v Any, t AsnType) error
	Decode   func(b Buffer, v Any, t AsnType) error
	LoadMeta func(content []byte) (AsnModule, error)
}

var SET_MASK []byte = []byte{0x80, 0x40, 0x20, 0x10, 0x08, 0x04, 0x02, 0x01}
var CLEAR_MASK []byte = []byte{0x7F, 0xBF, 0xDF, 0xEF, 0xF7, 0xFB, 0xFD, 0xFE}
var FIRST_MASK []byte = []byte{0xFF, 0x7F, 0x3F, 0x1F, 0x0F, 0x07, 0x03, 0x01, 0x00}
var SECOND_MASK []byte = []byte{0x00, 0x80, 0xC0, 0xE0, 0xF0, 0xF8, 0xFC, 0xFE, 0xFF}

func (self *BitString) Length() int {
	return (len(self.Bytes) << 3) - int(self.UnusedBits)
}

func (self *BitString) Bit(bit int) bool {
	offset := bit % 8
	index := bit / 8
	if self.Bytes == nil {
		return false
	}
	if index >= len(self.Bytes) {
		return false
	}
	if index == len(self.Bytes)-1 {
		if offset >= int(8-self.UnusedBits) {
			return false
		}
	}
	return (self.Bytes[index] & SET_MASK[offset]) != 0
}

func (self *BitString) SetBit(bit int, b bool) {
	size := self.Length()

	offset := bit % 8
	index := bit / 8
	if index >= len(self.Bytes) {
		self.Bytes = append(self.Bytes, make([]byte, len(self.Bytes)-index)...)
	}
	if b {
		self.Bytes[index] |= SET_MASK[offset]
	} else {
		self.Bytes[index] &= CLEAR_MASK[offset]
	}

	if bit >= size {
		size = bit + 1
		index = size % 8
		if index == 0 {
			self.UnusedBits = 0
		} else {
			self.UnusedBits = (uint8)(8 - index)
		}
	}
}

func (self *BitString) TrimTailingZeros() {
	if self.Bytes != nil && len(self.Bytes) > 0 {
		size := self.Length()
		var i int
		for i = size - 1; i >= 0; i-- {
			if self.Bit(i) {
				break
			}
		}
		bit := i + 1
		numBytes := ((bit - 1) >> 3) + 1
		numBits := bit % 8
		self.UnusedBits = 0
		if numBits != 0 {
			self.UnusedBits = uint8(8 - numBits)
		}
		self.Bytes = self.Bytes[0:numBytes]
		self.ClearUnusedBits()
	}
}

func (self *BitString) ClearUnusedBits() {
	if self.Bytes != nil && len(self.Bytes) > 0 {
		self.Bytes[len(self.Bytes)-1] &= SECOND_MASK[8-self.UnusedBits]
	}
}
