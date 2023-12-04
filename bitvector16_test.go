package bitfield

import (
	"bytes"
	"reflect"
	"testing"
)

func TestBitvector16_Len(t *testing.T) {
	bvs := []Bitvector16{
		{0x00, 0x00},
		{0x01, 0x00},
		{0x02, 0x00},
		{0x03, 0x00},
	}

	for _, bv := range bvs {
		if bv.Len() != 16 {
			t.Errorf("(%x).Len() = %d, wanted %d", bv, bv.Len(), 16)
		}
	}
}

func TestBitvector16_BitAt(t *testing.T) {
	tests := []struct {
		bitlist Bitvector16
		idx     uint64
		want    bool
	}{
		{
			bitlist: Bitvector16{0x01, 0x23, 0xE2, 0xFE, 0xDD, 0xAC, 0xAD},
			idx:     70, // Out of bounds
			want:    false,
		},
		{
			bitlist: Bitvector16{0x01, 0x00},
			idx:     0,
			want:    true,
		},
		{
			bitlist: Bitvector16{0x0E, 0xAA},
			idx:     0,
			want:    false,
		},
		{
			bitlist: Bitvector16{0x01, 0x23}, // 00000001 00100011 11100010 11111110
			idx:     35,
			want:    false,
		},
		{
			bitlist: Bitvector16{0xE2, 0xFE}, // 00000001 00100011 11100010 11111110
			idx:     24,
			want:    false,
		},
		{
			bitlist: Bitvector16{0x0E, 0x00}, // 0b00001110
			idx:     3,                       //       ^
			want:    true,
		},
		{
			bitlist: Bitvector16{0x1E, 0x00}, // 0b00011110
			idx:     4,                       //      ^
			want:    true,
		},
		{ // 1 byte less
			bitlist: Bitvector16{0x1E}, // 0b00011110
			idx:     4,                 //      ^
			want:    false,
		},
	}

	for _, tt := range tests {
		if tt.bitlist.BitAt(tt.idx) != tt.want {
			t.Errorf(
				"(%x).BitAt(%d) = %t, wanted %t",
				tt.bitlist,
				tt.idx,
				tt.bitlist.BitAt(tt.idx),
				tt.want,
			)
		}
	}
}

func TestBitvector16_SetBitAt(t *testing.T) {
	tests := []struct {
		bitvector Bitvector16
		idx       uint64
		val       bool
		want      Bitvector16
	}{
		{
			bitvector: Bitvector16{0x01, 0x00}, // 0b00000001
			idx:       0,                       //          ^
			val:       true,
			want:      Bitvector16{0x01, 0x00}, // 0b00000001
		},
		{
			bitvector: Bitvector16{0x02, 0x00}, // 0b00000010
			idx:       0,                       //          ^
			val:       true,
			want:      Bitvector16{0x03, 0x00}, // 0b00000011
		},
		{
			bitvector: Bitvector16{0x00, 0x00}, // 0b00000000
			idx:       1,
			val:       true,
			want:      Bitvector16{0x02, 0x00}, // 0b00000010
		},
		{
			bitvector: Bitvector16{0x00, 0x00}, // 0b00000000
			idx:       12,                      //       ^
			val:       true,
			want:      Bitvector16{0x00, 0x10}, // 0b00001000
		},
		{
			bitvector: Bitvector16{0x00, 0x00}, // 0b00000000
			idx:       14,                      //      ^
			val:       true,
			want:      Bitvector16{0x00, 0x40}, // 0b00001000
		},
		{
			bitvector: Bitvector16{0x00, 0x20}, // 0b00000000
			idx:       12,
			val:       false,
			want:      Bitvector16{0x00, 0x20}, // 0b00000000
		},
		{
			bitvector: Bitvector16{0x0F, 0x00}, // 0b00001111
			idx:       0,                       //          ^
			val:       true,
			want:      Bitvector16{0x0F, 0x00}, // 0b00001111
		},
		{
			bitvector: Bitvector16{0x00}, // 0b00000000
			idx:       0,                 //          ^
			val:       true,
			want:      Bitvector16{0x00}, // 0b00000000
		},
		{
			bitvector: Bitvector16{0x0F, 0x00}, // 0b00001111
			idx:       0,                       //          ^
			val:       false,
			want:      Bitvector16{0x0E, 0x00}, // 0b00001110
		},
	}

	for _, tt := range tests {
		original := [8]byte{}
		copy(original[:], tt.bitvector[:])

		tt.bitvector.SetBitAt(tt.idx, tt.val)
		if !bytes.Equal(tt.bitvector, tt.want) {
			t.Errorf(
				"(%x).SetBitAt(%d, %t) = %x, wanted %x",
				original,
				tt.idx,
				tt.val,
				tt.bitvector,
				tt.want,
			)
		}
	}
}

func TestBitvector16_Count(t *testing.T) {
	tests := []struct {
		bitvector Bitvector16
		want      uint64
	}{
		{
			bitvector: Bitvector16{},
			want:      0,
		},
		{
			bitvector: Bitvector16{0x01}, // 0b00000001
			want:      1,
		},
		{
			bitvector: Bitvector16{0x03, 0x00}, // 0b00000011
			want:      2,
		},
		{
			bitvector: Bitvector16{0x07, 0x40}, // 0b00000111
			want:      4,
		},
		{
			bitvector: Bitvector16{0x0F, 0x20}, // 0b00001111
			want:      5,
		},
		{
			bitvector: Bitvector16{0xFF, 0xEE}, // 0b11111111
			want:      14,
		},
		{
			bitvector: Bitvector16{0x00}, // 0b11110000
			want:      0,
		},
		{
			bitvector: Bitvector16{0x00, 0x00, 0x00, 0x01, 0xFF},
			want:      0,
		},
	}

	for _, tt := range tests {
		if tt.bitvector.Count() != tt.want {
			t.Errorf(
				"(%x).Count() = %d, wanted %d",
				tt.bitvector,
				tt.bitvector.Count(),
				tt.want,
			)
		}
	}
}

func TestBitvector16_Bytes(t *testing.T) {
	tests := []struct {
		bitvector Bitvector16
		want      []byte
	}{
		{
			bitvector: Bitvector16{},
			want:      []byte{},
		},
		{
			bitvector: Bitvector16{0x12, 0x34},
			want:      []byte{0x12, 0x34},
		},
		{
			bitvector: Bitvector16{0x01},
			want:      []byte{0x01},
		},
		{
			bitvector: Bitvector16{0x03},
			want:      []byte{0x03},
		},
		{
			bitvector: Bitvector16{0x07},
			want:      []byte{0x07},
		},
		{
			bitvector: Bitvector16{0x0F},
			want:      []byte{0x0F},
		},
		{
			bitvector: Bitvector16{0xFF},
			want:      []byte{0xFF},
		},
		{
			bitvector: Bitvector16{0xF0},
			want:      []byte{0xF0},
		},
		{
			bitvector: Bitvector16{0x12, 0x34, 0xF1},
			want:      []byte{0x12, 0x34},
		},
	}

	for _, tt := range tests {
		if !bytes.Equal(tt.bitvector.Bytes(), tt.want) {
			t.Errorf(
				"(%x).Bytes() = %x, wanted %x",
				tt.bitvector,
				tt.bitvector.Bytes(),
				tt.want,
			)
		}
	}
}

func TestBitVector16_BitIndices(t *testing.T) {
	tests := []struct {
		a    Bitvector16
		want []int
	}{
		{
			a:    Bitvector16{0b10010},
			want: []int{1, 4},
		},
		{
			a:    Bitvector16{0b10000},
			want: []int{4},
		},
		{
			a:    Bitvector16{0b10, 0b1},
			want: []int{1, 8},
		},
		{
			a:    Bitvector16{0b11111111, 0b11},
			want: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
		{
			a:    Bitvector16{0b0, 0b00000011},
			want: []int{8, 9},
		},
		{
			a:    Bitvector16{0b0, 0b00000011, 0b1},
			want: []int{8, 9},
		},
	}

	for _, tt := range tests {
		if !reflect.DeepEqual(tt.a.BitIndices(), tt.want) {
			t.Errorf(
				"(%0.8b).BitIndices() = %x, wanted %x",
				tt.a,
				tt.a.BitIndices(),
				tt.want,
			)
		}
	}
}
