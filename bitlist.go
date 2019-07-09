package bitfield

import (
	"math/bits"
)

var _ = Bitfield(Bitlist{})

// Bitlist is a bitfield implementation backed by an array of bytes. The most
// significant bit in the array of bytes indicates the start position of the
// bitfield.
//
// Examples of the underlying byte array as bitlist:
//  byte{0b00001000} is a bitlist with 3 bits which are all zero. bits=[0,0,0]
//  byte{0b00011111} is a bitlist with 4 bits which are all one.  bits=[1,1,1,1]
//  byte{0b00011000, 0b00000001} is a bitlist with 8 bits.        bits=[0,0,0,1,1,0,0,0]
//  byte{0b00011000, 0b00000010} is a bitlist with 9 bits.        bits=[0,0,0,0,1,1,0,0,0]
type Bitlist []byte

// NewBitlist creates a new bitlist of size N.
func NewBitlist(n uint64) Bitlist {
	ret := make(Bitlist, n/8+1)

	// Set most significant bit for length bit.
	i := uint8(1 << (n % 8))
	ret[n/8] |= i

	return ret
}

// BitAt returns the bit value at the given index. If the index requested
// exceeds the number of bits in the bitlist, then this method returns false.
func (b Bitlist) BitAt(idx uint64) bool {
	// Out of bounds, must be false.
	upperBounds := b.Len()
	if idx >= upperBounds {
		return false
	}

	i := uint8(1 << (idx % 8))
	return b[idx/8]&i == i
}

// SetBitAt will set the bit at the given index to the given value. If the index
// requested exceeds the number of bits in the bitlist, then this method returns
// false.
func (b Bitlist) SetBitAt(idx uint64, val bool) {
	// Out of bounds, do nothing.
	upperBounds := b.Len()
	if idx >= upperBounds {
		return
	}

	bit := uint8(1 << (idx % 8))
	if val {
		b[idx/8] |= bit
	} else {
		b[idx/8] &^= bit
	}

}

// Len of the bitlist returns the number of bits available in the underlying
// byte array.
func (b Bitlist) Len() uint64 {
	if len(b) == 0 {
		return 0
	}
	// The most significant bit is present in the last byte in the array.
	last := b[len(b)-1]

	// Determine the position of the most significant bit.
	msb := bits.Len8(last)

	// The absolute position of the most significant bit will be the number of
	// bits in the preceding bytes plus the position of the most significant
	// bit. Subtract this value by 1 to determine the length of the bitlist.
	return uint64(8*(len(b)-1) + msb - 1)
}

// Bytes returns the trimmed underlying byte array without the length bit. The
// leading zeros in the bitlist will be trimmed to the smallest byte length
// representation of the bitlist. This may produce an empty byte slice if all
// bits were zero.
func (b Bitlist) Bytes() []byte {
	if len(b) == 0 {
		return []byte{}
	}

	ret := make([]byte, len(b))
	copy(ret, b)

	// Clear the most significant bit (the length bit).
	msb := uint8(bits.Len8(ret[len(ret)-1])) - 1
	clearBit := uint8(1 << msb)
	ret[len(ret)-1] &^= clearBit

	// Clear any leading zero bytes.
	newLen := len(ret)
	for i := len(ret) - 1; i >= 0; i-- {
		if ret[i] != 0x00 {
			break
		}
		newLen = i
	}

	return ret[:newLen]
}

// Count returns the number of 1s in the bitlist.
func (b Bitlist) Count() uint64 {
	c := 0

	for _, bt := range b {
		c += bits.OnesCount8(bt)
	}

	if c > 0 {
		c-- // Remove length bit from count.
	}

	return uint64(c)
}
