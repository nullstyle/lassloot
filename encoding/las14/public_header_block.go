package las14

import (
	"fmt"
	"strings"
	"text/template"
)

func init() {

	summaryTemplate = template.Must(template.New("Public Header Block Summary").Parse(SummaryTemplateSource))
}

const SummaryTemplateSource = `
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
	summaryTemplate *template.Template
)

func (phb *PublicHeaderBlock) String() string {
	var ret strings.Builder

	err := summaryTemplate.Execute(&ret, phb)
	if err != nil {
		return fmt.Sprintf("header summary invalid: %v", err)
	}

	return ret.String()
}
