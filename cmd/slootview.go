package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/nullstyle/lassloot"
	. "github.com/nullstyle/lassloot/cmd/internal/helpers"
	"log"
	"os"
)

var (
	jsonFlag = flag.Bool("json", false, "output result as json")
)

func main() {
	flag.Parse()

	err, pc := lassloot.NewPointCloudFromPath(LasPathFromArgs())
	if err != nil {
		log.Fatalf("failed to create PointCloud: %w", err)
	}

	if *jsonFlag {
		err := json.NewEncoder(os.Stdout).Encode(pc)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		fmt.Println("Header:\n%s\n", pc.Header())
	}
}
