package des

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleEncryption(t *testing.T) {
	input := []uint8{
		0xFD, 0x0E, 0x64, 0x06, 0x65, 0xBE, 0x74, 0x13, //
		0x77, 0x63, 0x3B, 0x02, 0x45, 0x4E, 0x70, 0x7A, //
	}
	expected := []uint8{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6}
	des := New([]byte("TEST!KEY"), true)
	assert.True(t, des.TransformBytes(input))
	assert.Equal(t, expected, input)
}

func TestSimpleDecryption(t *testing.T) {
	input := []uint8{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6}
	expected := []uint8{
		0xFD, 0x0E, 0x64, 0x06, 0x65, 0xBE, 0x74, 0x13, //
		0x77, 0x63, 0x3B, 0x02, 0x45, 0x4E, 0x70, 0x7A, //
	}
	des := New([]byte("TEST!KEY"), false)
	assert.True(t, des.TransformBytes(input))
	assert.Equal(t, expected, input)
}
