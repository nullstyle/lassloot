package lassloot

import (
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

type Header struct {
	RawHeader las14.PublicHeaderBlock
	pc        *PointCloud
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
	err, fr := d.FullDecode()

	return nil, &PointCloud{fr}
}
func (pc *PointCloud) Header() *Header {
	return &Header{
		RawHeader: pc.fr.Header,
	}
}
