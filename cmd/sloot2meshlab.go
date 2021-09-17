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

func main() {
	flag.Parse()

	err, pc := lassloot.NewPointCloudFromPath(LasPathFromArgs())
	if err != nil {
		log.Fatalf("failed to create PointCloud: %w", err)
	}

	w := csv.NewWriter(os.Stdout)

	l := pc.Len()
	for i := (uint64)(0); i < l; i++ {
		err, point := pc.PointAt(i)
		if err != nil {
			log.Fatalf("failed to get point: %w", err)
		}
		err = w.Write(pointToMeshlab(point))
		if err != nil {
			log.Fatalf("failed to write point: %w", err)
		}
	}

	if err := w.Error(); err != nil {
		log.Fatalln("error writing csv:", err)
	}

	w.Flush()
}

func pointToMeshlab(p *lassloot.Point) []string {
	if p == nil {
		panic("don't pass a nil dipshit")
	}

	x, y, z := p.UnoffsetXYZ()
	return []string{
		// NOTE:  Meshlab uses Y for height, where LAS files are Z for height
		fmt.Sprintf("%f", x),
		fmt.Sprintf("%f", z),
		fmt.Sprintf("%f", y),
	}
}
