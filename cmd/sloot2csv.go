package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/nullstyle/lassloot"
	. "github.com/nullstyle/lassloot/cmd/internal/helpers"
	"log"
	"os"
)

var (
	unscaledFlag = flag.Bool("unscaled", false, "output result unscaled by file's scale factors")
	unoffsetFlag = flag.Bool("unoffset", false, "output result scaled, but not offset by file's scale factors")
)

func main() {
	flag.Parse()

	err, pc := lassloot.NewPointCloudFromPath(LasPathFromArgs())
	if err != nil {
		log.Fatalf("failed to create PointCloud: %w", err)
	}

	w := csv.NewWriter(os.Stdout)
	err = w.Write([]string{"x", "y", "z"})
	if err != nil {
		log.Fatalf("failed to write csv header: %w", err)
	}

	l := pc.Len()
	for i := (uint64)(0); i < l; i++ {
		err, point := pc.PointAt(i)
		if err != nil {
			log.Fatalf("failed to get point: %w", err)
		}
		err = w.Write(pointToCSV(point))
		if err != nil {
			log.Fatalf("failed to write csv header: %w", err)
		}
	}

	if err := w.Error(); err != nil {
		log.Fatalln("error writing csv:", err)
	}

	w.Flush()
}

func pointToCSV(p *lassloot.Point) []string {
	if p == nil {
		panic("don't pass a nil dipshit")
	}

	switch {
	case *unoffsetFlag:
		x, y, z := p.UnoffsetXYZ()
		return []string{
			fmt.Sprintf("%f", x),
			fmt.Sprintf("%f", y),
			fmt.Sprintf("%f", z),
		}
	case *unscaledFlag:
		x, y, z := p.UnscaledXYZ()
		return []string{
			fmt.Sprintf("%d", x),
			fmt.Sprintf("%d", y),
			fmt.Sprintf("%d", z),
		}
	default:
		x, y, z := p.XYZ()
		return []string{
			fmt.Sprintf("%f", x),
			fmt.Sprintf("%f", y),
			fmt.Sprintf("%f", z),
		}
	}
}
