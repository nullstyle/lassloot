package lassloot

import (
	"fmt"
	"github.com/nullstyle/lassloot/encoding/las14"
	"log"
	"os"
)

// PointCloud is the primary API for interacting with LAS files provided by this library.
// It represents a high level interface that wraps a lower level parse result,
// providing caching and higher level algorithms.
type PointCloud struct {
	fr *las14.FullResult
}

func NewPointCloudFromPath(path string) (error, *PointCloud) {
	f, err := os.Open(path)
	if err != nil {
		return err, nil
	}
	defer func() {
		err := f.Close()
		if err != nil {
			log.Printf("error occurred closing las file: %w", err)
		}
	}()

	d := las14.NewDecoder(f)
	err, fr := d.FullDecode("")

	return nil, &PointCloud{fr}
}

func (pc *PointCloud) Header() *Header {
	return &Header{
		RawHeader: pc.fr.Header,
	}
}

func (pc *PointCloud) Len() uint64 {
	return (uint64)(pc.fr.Header.LegacyNumberOfPointRecords)
}

func (pc *PointCloud) PointSize() int {
	return (int)(pc.fr.Header.PointDataRecordLength)
}

func (pc *PointCloud) PointAt(idx uint64) (error, *Point) {

	if idx > pc.Len() {
		return fmt.Errorf("index %d too high", idx), nil
	}

	return nil, &Point{
		pc:  pc,
		PDR: pc.fr.PointDataRecord(idx),
	}
}

// LocalizeXYZ scales and offsets the provided xyz according to the pointcloud's offset and scale as reflected in the
// header.
func (pc *PointCloud) LocalizeXYZ(ix int64, iy int64, iz int64) (ox float64, oy float64, oz float64) {
	h := pc.fr.Header
	ox = ((float64)(ix) * h.XScaleFactor) + h.XOffset
	oy = ((float64)(iy) * h.YScaleFactor) + h.YOffset
	oz = ((float64)(iz) * h.ZScaleFactor) + h.ZOffset
	return
}

// ScaleXYZ scales (but doesnt offset) the provided xyz according to the pointcloud's scale as reflected in
// the header.
func (pc *PointCloud) ScaleXYZ(ix int64, iy int64, iz int64) (ox float64, oy float64, oz float64) {
	h := pc.fr.Header
	ox = ((float64)(ix) * h.XScaleFactor)
	oy = ((float64)(iy) * h.YScaleFactor)
	oz = ((float64)(iz) * h.ZScaleFactor)
	return
}

type Header struct {
	RawHeader las14.PublicHeaderBlock
	pc        *PointCloud
}

type Point struct {
	PDR *las14.PointDataRecord
	pc  *PointCloud
}

func (p *Point) XYZ() (x float64, y float64, z float64) {
	pd := p.PDR.Get()
	x, y, z = p.pc.LocalizeXYZ(pd.XYZ())

	return
}

func (p *Point) TruncatedXYZ() (x int64, y int64, z int64) {
	pd := p.PDR.Get()
	fx, fy, fz := p.pc.LocalizeXYZ(pd.XYZ())

	x = (int64)(fx)
	y = (int64)(fy)
	z = (int64)(fz)
	return
}

func (p *Point) UnscaledXYZ() (x int64, y int64, z int64) {
	pd := p.PDR.Get()
	x, y, z = pd.XYZ()

	return
}

func (p *Point) UnoffsetXYZ() (x float64, y float64, z float64) {
	pd := p.PDR.Get()
	x, y, z = p.pc.ScaleXYZ(pd.XYZ())

	return
}
