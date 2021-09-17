package las14

import (
	"encoding/binary"
	"fmt"
)

func init() {

}

type pdr6 struct {
	pdr *PointDataRecord
}

func (p *pdr6) XYZ() (x int64, y int64, z int64) {
	x = (int64)(binary.LittleEndian.Uint32(p.pdr.Raw[0:4]))
	y = (int64)(binary.LittleEndian.Uint32(p.pdr.Raw[4:8]))
	z = (int64)(binary.LittleEndian.Uint32(p.pdr.Raw[8:12]))
	return
}

func (p *pdr6) Intensity() uint16 {
	return binary.LittleEndian.Uint16(p.pdr.Raw[12:14])
}

func (p *pdr6) Classification() byte {
	return p.pdr.Raw[16]
}

var _ = (*pdr6)(nil)

func (pdr *PointDataRecord) Get() PointData {
	switch pdr.Format {
	case 6:
		return &pdr6{pdr}
	default:
		panic(fmt.Sprintf("unhandled format encountered: %d", pdr.Format))
	}
}
