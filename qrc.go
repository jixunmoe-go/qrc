package qrc

import "github.com/jixunmoe-go/qrc/internal/qrc"

func DecodeQRC(data []uint8) ([]uint8, error) {
	return qrc.DecodeQRC(data)
}
