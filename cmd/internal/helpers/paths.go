package helpers

import (
	"flag"
	"log"
	"path/filepath"
)

func LasPathFromArgs() string {
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

	return lasPath
}
