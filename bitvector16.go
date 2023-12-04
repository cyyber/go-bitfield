package bitfield

import (
	"math/bits"
)

var _ = Bitfield(Bitvector16{})

// Bitvector16 is a bitfield with a fixed defined size of 16. There is no length bit
// present in the underlying byte array.
type Bitvector16 []byte

const bitvector16ByteSize = 2
const bitvector16BitSize = bitvector16ByteSize * 8

// NewBitvector16 creates a new bitvector of size 16.
func NewBitvector16() Bitvector16 {
	byteArray := [bitvector16ByteSize]byte{}
	return byteArray[:]
}

// BitAt returns the bit value at the given index. If the index requested
// exceeds the number of bits in the bitvector, then this method returns false.
func (b Bitvector16) BitAt(idx uint64) bool {
	// Out of bounds, must be false.
	if idx >= b.Len() || len(b) != bitvector16ByteSize {
		return false
	}

	i := uint8(1 << (idx % 8))
	return b[idx/8]&i == i
}

// SetBitAt will set the bit at the given index to the given value. If the index
// requested exceeds the number of bits in the bitvector, then this method returns
// false.
func (b Bitvector16) SetBitAt(idx uint64, val bool) {
	// Out of bounds, do nothing.
	if idx >= b.Len() || len(b) != bitvector16ByteSize {
		return
	}

	bit := uint8(1 << (idx % 8))
	if val {
		b[idx/8] |= bit
	} else {
		b[idx/8] &^= bit
	}
}

// Len returns the number of bits in the bitvector.
func (b Bitvector16) Len() uint64 {
	return bitvector16BitSize
}

// Count returns the number of 1s in the bitvector.
func (b Bitvector16) Count() uint64 {
	if len(b) == 0 {
		return 0
	}
	c := 0
	for i, bt := range b {
		if i >= bitvector16ByteSize {
			break
		}
		c += bits.OnesCount8(bt)
	}
	return uint64(c)
}

// Bytes returns the bytes data representing the bitvector16.
func (b Bitvector16) Bytes() []byte {
	if len(b) == 0 {
		return []byte{}
	}
	ln := min(len(b), bitvector16ByteSize)
	ret := make([]byte, ln)
	copy(ret, b[:ln])
	return ret[:]
}

// BitIndices returns the list of indices which are set to 1.
func (b Bitvector16) BitIndices() []int {
	indices := make([]int, 0, bitvector16BitSize)
	for i, bt := range b {
		if i >= bitvector16ByteSize {
			break
		}
		for j := 0; j < 8; j++ {
			bit := byte(1 << uint(j))
			if bt&bit == bit {
				indices = append(indices, i*8+j)
			}
		}
	}

	return indices
}
