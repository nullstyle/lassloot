package las14

// HeaderMagicBytes are the characters that occur at the beginning of a LAS header block.
// See the "File Signature" definition on page 5 of the OGC version of the LAS 1.4 spec
const HeaderMagicBytes = "LASF"

// protocol types

// PublicHeaderBlock represents a parsed header block for an ASPRS LAS 1.4 File
type PublicHeaderBlock struct {
	FileSourceID                             uint16
	GlobalEncoding                           GlobalEncodingBitField
	ProjectID                                ProjectID
	VersionMajor                             byte
	VersionMinor                             byte
	SystemID                                 SystemID
	GeneratingSoftware                       [32]byte
	FileCreationDayOfYear                    uint16
	FileCreationYear                         uint16
	HeaderSize                               uint16
	OffsetToPointData                        uint32
	NumberOfVariableLengthRecords            uint32
	PointDataRecordFormat                    byte
	PointDataRecordLength                    uint16
	LegacyNumberOfPointRecords               uint32
	LegacyNumberOfPointByReturn              [5]uint32
	XScaleFactor                             float64
	YScaleFactor                             float64
	ZScaleFactor                             float64
	XOffset                                  float64
	YOffset                                  float64
	ZOffset                                  float64
	MaxX                                     float64
	MinX                                     float64
	MaxY                                     float64
	MinY                                     float64
	MaxZ                                     float64
	MinZ                                     float64
	StartOfWaveformDataPacketRecord          uint64
	StartOfFirstExtendedVariableLengthRecord uint64
	NumberOfExtendedVariableLengthRecords    uint64
	NumberOfPointRecords                     uint64
	NumberOfPointsByReturn                   [15]uint64
}

type GlobalEncodingBitField uint16
type ProjectID struct {
}

type SystemID [32]byte

type VariableLengthRecord struct{}
type PointDataRecord struct{}
type ExtendedVariableLengthRecord struct{}
