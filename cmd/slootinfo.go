package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/nullstyle/lassloot"
	"log"
	"os"
	"text/template"

	. "github.com/nullstyle/lassloot/cmd/internal/helpers"
)

var (
	jsonFlag = flag.Bool("json", false, "output result as json")
)

const reponseTemplateSource = `
	HeaderSize                               = {{.HeaderSize}}
	OffsetToPointData                        = {{.OffsetToPointData}}
	StartOfWaveformDataPacketRecord          = {{.StartOfWaveformDataPacketRecord}}
	StartOfFirstExtendedVariableLengthRecord = {{.StartOfFirstExtendedVariableLengthRecord}}

	LegacyNumberOfPointRecords               = {{.LegacyNumberOfPointRecords}}
	NumberOfPointRecords                     = {{.NumberOfPointRecords}}

	NumberOfVariableLengthRecords            = {{.NumberOfVariableLengthRecords}}
	NumberOfExtendedVariableLengthRecords    = {{.NumberOfExtendedVariableLengthRecords}}
`

var (
	responseTemplate *template.Template
)

func main() {
	flag.Parse()
	responseTemplate = template.Must(template.New("Info Response").Parse(reponseTemplateSource))

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

type infoResponse struct {
}

func (ir *infoResponse) String() {

}
