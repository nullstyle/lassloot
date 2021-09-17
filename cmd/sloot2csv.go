package main

import (
	"encoding/csv"
	"flag"
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
	w.Write([]string{"x", "y", "z"})
	l := pc.Len()
	for i := (uint64)(0); i < l; i++ {
		err, point := pc.PointAt(i)
		if err != nil {
			log.Fatalf("failed to get point: %w", err)
		}
		w.Write(point.CSV())
	}

	if err := w.Error(); err != nil {
		log.Fatalln("error writing csv:", err)
	}

	w.Flush()
}
