package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/nullstyle/lassloot"
	"log"
	"os"
	"path/filepath"
)

var jsonFlag = flag.Bool("json", false, "output result as json")

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		log.Fatalln("too many args")
	}

	path, err := filepath.EvalSymlinks(args[0])
	if err != nil {
		log.Fatalln(err)
	}

	lasPath, err := filepath.Abs(path)
	if err != nil {
		log.Fatalln(err)
	}

	err, pc := lassloot.NewPointCloudFromPath(lasPath)
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
