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

type Header struct {
	RawHeader las14.PublicHeaderBlock
	pc        *PointCloud
}

type Point struct {
	PDR *las14.PointDataRecord
	pc  *PointCloud
}

func (p *Point) CSV() []string {
	data := p.PDR.Get()
	x, y, z := data.At()
	return []string{
		fmt.Sprintf("%d", x),
		fmt.Sprintf("%d", y),
		fmt.Sprintf("%d", z),
	}
}
