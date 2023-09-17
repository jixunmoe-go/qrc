package qrc

import (
	"bytes"
	"compress/zlib"
	"io"

	"github.com/jixunmoe-go/qrc/internal/des"
	"github.com/jixunmoe-go/qrc/internal/qmc"
)

var qrc_qmc_magic = []uint8{0x98, 0x25, 0xB0, 0xAC, 0xE3, 0x02, 0x83, 0x68, 0xE8, 0xFC, 0x6C}

var des_key1 = []byte("!@#)(NHL")
var des_key2 = []byte("123ZXC!@")
var des_key3 = []byte("!@#)(*$%")

func isMagicEq(data []uint8) bool {
	if len(data) < len(qrc_qmc_magic) {
		return false
	}

	for i := 0; i < len(qrc_qmc_magic); i++ {
		if data[i] != qrc_qmc_magic[i] {
			return false
		}
	}
	return true
}

func DecodeQRC(data []uint8) ([]uint8, error) {
	var result []uint8
	// if data starts with qrc_qmc_magic
	if isMagicEq(data) {
		result = qmc.QmcDecode(data)
		result = result[len(qrc_qmc_magic):]
	} else {
		result = data
	}

	des.New(des_key1, false).TransformBytes(result)
	des.New(des_key2, true).TransformBytes(result)
	des.New(des_key3, false).TransformBytes(result)

	zlib_reader, err := zlib.NewReader(bytes.NewReader(result))
	if err != nil {
		return nil, err
	}
	return io.ReadAll(zlib_reader)
}
