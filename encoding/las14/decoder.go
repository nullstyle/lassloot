package las14

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"sync"
)

const Las14HeaderSize = 375

// DefaultBudget represents the default read budget provided to a freshly initialized decoder.  1 gigabyte seems like a reasonable limit to decode large well-formed files while still limiting exposure denial of service attacks due to an implementation bug.  See Decoder#safeRead for budget-based reading code.
var DefaultBudget uint = 1000 * (1024 * 1024)

// A Decoder reads and decodes LAS 1.4 files from an input stream.
type Decoder struct {
	r      io.ReadSeeker
	mt     sync.Mutex
	budget uint

	fp  *FirstPassResult
	ret *FullResult
}

// NewDecoder returns a new decoder that reads from r.
func NewDecoder(r io.ReadSeeker) *Decoder {
	return &Decoder{r: r, budget: DefaultBudget}
}

type FirstPassResult struct {
	Header PublicHeaderBlock
}

type FullResult struct {
	FirstPassResult

	pointData []byte
}

func (fr *FullResult) PointDataRecord(idx uint64) *PointDataRecord {
	offset := fr.pointOffset(idx)

	return &PointDataRecord{
		Raw:    fr.pointData[offset : offset+(uint64)(fr.Header.PointDataRecordLength)],
		Format: fr.Header.PointDataRecordFormat,
	}
}

func (fr *FullResult) pointOffset(i uint64) uint64 {
	return i * (uint64)(fr.Header.PointDataRecordLength)
}

func (las *Decoder) FirstPassDecode() (error, *FirstPassResult) {
	las.mt.Lock()
	defer las.mt.Unlock()

	return las.firstPassDecode()
}

type QuerySet string

func (las *Decoder) FullDecode(qs QuerySet) (error, *FullResult) {
	las.mt.Lock()
	defer las.mt.Unlock()

	return las.fullDecode(qs)
}

func (las *Decoder) firstPassDecode() (error, *FirstPassResult) {
	if las.fp != nil {
		return nil, las.fp
	}

	var header PublicHeaderBlock
	cur, err := las.r.Seek(0, io.SeekStart)
	if err != nil {
		return fmt.Errorf("failed to seek to header: %w", err), nil
	}

	if cur != 0 {
		panic("invalid offset returned when seeking to start of file")
	}

	actualSig := make([]byte, 4)
	n, err := las.safeRead(actualSig)
	if err != nil {
		return fmt.Errorf("failed to read header: %w", err), nil
	}

	if n != 4 {
		return fmt.Errorf("failed to read full las header: only %d bytes read", n), nil
	}

	if (string)(actualSig) != HeaderMagicBytes {
		return fmt.Errorf("invalid file signature: read %s", actualSig), nil
	}

	// begin decoding actual header data

	fileSourceID := make([]byte, 2)
	n, err = las.safeRead(fileSourceID)
	if err != nil {
		return fmt.Errorf("failed to read header: %w", err), nil
	}
	header.FileSourceID = binary.LittleEndian.Uint16(fileSourceID)

	geb := make([]byte, 2)
	n, err = las.safeRead(geb)
	if err != nil {
		return fmt.Errorf("failed to read header: %w", err), nil
	}
	header.GlobalEncoding = (GlobalEncodingBitField)(binary.LittleEndian.Uint16(geb))

	projectID := make([]byte, 16)
	n, err = las.safeRead(projectID)
	if err != nil {
		return fmt.Errorf("failed to read header: %w", err), nil
	}
	header.ProjectID.Data1 = binary.LittleEndian.Uint32(projectID[0:4])
	header.ProjectID.Data2 = binary.LittleEndian.Uint16(projectID[4:6])
	header.ProjectID.Data3 = binary.LittleEndian.Uint16(projectID[6:8])
	header.ProjectID.Data4 = binary.LittleEndian.Uint64(projectID[8:])

	versions := make([]byte, 2)
	n, err = las.safeRead(versions)
	if err != nil {
		return fmt.Errorf("failed to read header: %w", err), nil
	}
	header.VersionMajor = versions[0]
	header.VersionMinor = versions[1]

	systemID := make([]byte, 32)
	n, err = las.safeRead(systemID)
	if err != nil {
		return fmt.Errorf("failed to read header: %w", err), nil
	}
	copy(header.SystemID[0:32], systemID)

	generatingSoftware := make([]byte, 32)
	n, err = las.safeRead(generatingSoftware)
	if err != nil {
		return fmt.Errorf("failed to read header: %w", err), nil
	}
	copy(header.GeneratingSoftware[0:32], generatingSoftware)

	createDOY := make([]byte, 2)
	n, err = las.safeRead(createDOY)
	if err != nil {
		return fmt.Errorf("failed to read header: %w", err), nil
	}
	header.FileCreationDayOfYear = binary.LittleEndian.Uint16(createDOY)

	createY := make([]byte, 2)
	n, err = las.safeRead(createY)
	if err != nil {
		return fmt.Errorf("failed to read header: %w", err), nil
	}
	header.FileCreationYear = binary.LittleEndian.Uint16(createY)

	headerSize := make([]byte, 2)
	n, err = las.safeRead(headerSize)
	if err != nil {
		return fmt.Errorf("failed to read header: %w", err), nil
	}
	header.HeaderSize = binary.LittleEndian.Uint16(headerSize)
	if header.HeaderSize != Las14HeaderSize {
		return fmt.Errorf("unrecognized header size: lassloot only handles LAS 1.4 headers of 375 bytes"), nil
	}

	offsetToPoints := make([]byte, 4)
	n, err = las.safeRead(offsetToPoints)
	if err != nil {
		return fmt.Errorf("failed to read header: %w", err), nil
	}
	header.OffsetToPointData = binary.LittleEndian.Uint32(offsetToPoints)

	nVLR := make([]byte, 4)
	n, err = las.safeRead(nVLR)
	if err != nil {
		return fmt.Errorf("failed to read header: %w", err), nil
	}
	header.NumberOfVariableLengthRecords = binary.LittleEndian.Uint32(nVLR)

	pdrFormat := make([]byte, 1)
	n, err = las.safeRead(pdrFormat)
	if err != nil {
		return fmt.Errorf("failed to read header: %w", err), nil
	}
	header.PointDataRecordFormat = (PointDataFormat)(pdrFormat[0])

	pdrLength := make([]byte, 2)
	n, err = las.safeRead(pdrLength)
	if err != nil {
		return fmt.Errorf("failed to read header: %w", err), nil
	}
	header.PointDataRecordLength = binary.LittleEndian.Uint16(pdrLength)

	lnpr := make([]byte, 4)
	n, err = las.safeRead(lnpr)
	if err != nil {
		return fmt.Errorf("failed to read header: %w", err), nil
	}
	header.LegacyNumberOfPointRecords = binary.LittleEndian.Uint32(lnpr)

	lnpbr := make([]byte, binary.Size(header.LegacyNumberOfPointsByReturn))
	n, err = las.safeRead(lnpbr)
	if err != nil {
		return fmt.Errorf("failed to read header: %w", err), nil
	}
	header.LegacyNumberOfPointsByReturn[0] = binary.LittleEndian.Uint32(lnpbr[0:4])
	header.LegacyNumberOfPointsByReturn[1] = binary.LittleEndian.Uint32(lnpbr[4:8])
	header.LegacyNumberOfPointsByReturn[2] = binary.LittleEndian.Uint32(lnpbr[8:12])
	header.LegacyNumberOfPointsByReturn[3] = binary.LittleEndian.Uint32(lnpbr[12:16])
	header.LegacyNumberOfPointsByReturn[4] = binary.LittleEndian.Uint32(lnpbr[16:])

	xscale := make([]byte, 8)
	n, err = las.safeRead(xscale)
	if err != nil {
		return fmt.Errorf("failed to read header: %w", err), nil
	}
	header.XScaleFactor = math.Float64frombits(binary.LittleEndian.Uint64(xscale))

	yscale := make([]byte, 8)
	n, err = las.safeRead(yscale)
	if err != nil {
		return fmt.Errorf("failed to read header: %w", err), nil
	}
	header.YScaleFactor = math.Float64frombits(binary.LittleEndian.Uint64(yscale))

	zscale := make([]byte, 8)
	n, err = las.safeRead(zscale)
	if err != nil {
		return fmt.Errorf("failed to read header: %w", err), nil
	}
	header.ZScaleFactor = math.Float64frombits(binary.LittleEndian.Uint64(zscale))

	xoffset := make([]byte, 8)
	n, err = las.safeRead(xoffset)
	if err != nil {
		return fmt.Errorf("failed to read header: %w", err), nil
	}
	header.XOffset = math.Float64frombits(binary.LittleEndian.Uint64(xoffset))

	yoffset := make([]byte, 8)
	n, err = las.safeRead(yoffset)
	if err != nil {
		return fmt.Errorf("failed to read header: %w", err), nil
	}
	header.YOffset = math.Float64frombits(binary.LittleEndian.Uint64(yoffset))

	zoffset := make([]byte, 8)
	n, err = las.safeRead(zoffset)
	if err != nil {
		return fmt.Errorf("failed to read header: %w", err), nil
	}
	header.ZOffset = math.Float64frombits(binary.LittleEndian.Uint64(zoffset))

	maxx := make([]byte, 8)
	n, err = las.safeRead(maxx)
	if err != nil {
		return fmt.Errorf("failed to read header: %w", err), nil
	}
	header.MaxX = math.Float64frombits(binary.LittleEndian.Uint64(maxx))

	minx := make([]byte, 8)
	n, err = las.safeRead(minx)
	if err != nil {
		return fmt.Errorf("failed to read header: %w", err), nil
	}
	header.MinX = math.Float64frombits(binary.LittleEndian.Uint64(minx))

	maxy := make([]byte, 8)
	n, err = las.safeRead(maxy)
	if err != nil {
		return fmt.Errorf("failed to read header: %w", err), nil
	}
	header.MaxY = math.Float64frombits(binary.LittleEndian.Uint64(maxy))

	miny := make([]byte, 8)
	n, err = las.safeRead(miny)
	if err != nil {
		return fmt.Errorf("failed to read header: %w", err), nil
	}
	header.MinY = math.Float64frombits(binary.LittleEndian.Uint64(miny))

	maxz := make([]byte, 8)
	n, err = las.safeRead(maxz)
	if err != nil {
		return fmt.Errorf("failed to read header: %w", err), nil
	}
	header.MaxZ = math.Float64frombits(binary.LittleEndian.Uint64(maxz))

	minz := make([]byte, 8)
	n, err = las.safeRead(minz)
	if err != nil {
		return fmt.Errorf("failed to read header: %w", err), nil
	}
	header.MinZ = math.Float64frombits(binary.LittleEndian.Uint64(minz))

	startOfWaveform := make([]byte, 8)
	n, err = las.safeRead(startOfWaveform)
	if err != nil {
		return fmt.Errorf("failed to read header: %w", err), nil
	}
	header.StartOfWaveformDataPacketRecord = binary.LittleEndian.Uint64(startOfWaveform)

	startOfEVLR := make([]byte, 8)
	n, err = las.safeRead(startOfEVLR)
	if err != nil {
		return fmt.Errorf("failed to read header: %w", err), nil
	}
	header.StartOfFirstExtendedVariableLengthRecord = binary.LittleEndian.Uint64(startOfEVLR)

	nEVLR := make([]byte, 8)
	n, err = las.safeRead(nEVLR)
	if err != nil {
		return fmt.Errorf("failed to read header: %w", err), nil
	}
	header.NumberOfExtendedVariableLengthRecords = binary.LittleEndian.Uint64(nEVLR)

	npr := make([]byte, 8)
	n, err = las.safeRead(npr)
	if err != nil {
		return fmt.Errorf("failed to read header: %w", err), nil
	}
	header.NumberOfPointRecords = binary.LittleEndian.Uint64(npr)

	npbr := make([]byte, binary.Size(header.NumberOfPointsByReturn))
	n, err = las.safeRead(npbr)
	if err != nil {
		return fmt.Errorf("failed to read header: %w", err), nil
	}
	header.NumberOfPointsByReturn[0] = binary.BigEndian.Uint64(npbr[0:8])
	header.NumberOfPointsByReturn[1] = binary.BigEndian.Uint64(npbr[8:16])
	header.NumberOfPointsByReturn[2] = binary.BigEndian.Uint64(npbr[16:24])
	header.NumberOfPointsByReturn[3] = binary.BigEndian.Uint64(npbr[24:32])
	header.NumberOfPointsByReturn[4] = binary.BigEndian.Uint64(npbr[32:40])
	header.NumberOfPointsByReturn[5] = binary.BigEndian.Uint64(npbr[40:48])
	header.NumberOfPointsByReturn[6] = binary.BigEndian.Uint64(npbr[48:56])
	header.NumberOfPointsByReturn[7] = binary.BigEndian.Uint64(npbr[56:64])
	header.NumberOfPointsByReturn[8] = binary.BigEndian.Uint64(npbr[64:72])
	header.NumberOfPointsByReturn[9] = binary.BigEndian.Uint64(npbr[72:80])
	header.NumberOfPointsByReturn[10] = binary.BigEndian.Uint64(npbr[80:88])
	header.NumberOfPointsByReturn[11] = binary.BigEndian.Uint64(npbr[88:96])
	header.NumberOfPointsByReturn[12] = binary.BigEndian.Uint64(npbr[96:104])
	header.NumberOfPointsByReturn[13] = binary.BigEndian.Uint64(npbr[104:112])
	header.NumberOfPointsByReturn[14] = binary.BigEndian.Uint64(npbr[112:120])

	las.fp = &FirstPassResult{
		Header: header,
	}
	return nil, las.fp
}

func (las *Decoder) fullDecode(qs QuerySet) (error, *FullResult) {
	var err error
	var fp *FirstPassResult

	// populate fp
	if las.fp != nil {
		fp = las.fp
	} else {
		err, fp = las.firstPassDecode()
		if err != nil {
			return fmt.Errorf("full decode failed: invoked first pass decode failed with %w", err), nil
		}
	}

	// populate full result
	_, err = las.r.Seek((int64)(fp.Header.OffsetToPointData), io.SeekStart)
	if err != nil {
		return fmt.Errorf("failed to seek to start of points: %w", err), nil
	}

	pointDataSize := (uint32)(fp.Header.PointDataRecordLength) * fp.Header.LegacyNumberOfPointRecords
	pointData := make([]byte, pointDataSize)
	n, err := las.safeRead(pointData)
	if err != nil {
		return fmt.Errorf("failed to read point data: %w", err), nil
	}

	if n != (int)(pointDataSize) {
		return fmt.Errorf("could not read full point data: %w", err), nil
	}

	return nil, &FullResult{
		FirstPassResult: *fp,
		pointData:       pointData,
	}
}

func (las *Decoder) safeRead(p []byte) (n int, err error) {
	requested := uint(len(p))
	if requested > las.budget {
		return 0, fmt.Errorf("read budget exhausted: %d byte read requested", requested)
	}
	// NOTE: we deduct the requested byte count rather than the actually-read byte count from the budget because we are crotchety bastards #dealwithit
	las.budget -= requested

	n, err = las.r.Read(p)
	if err != nil {
		return n, fmt.Errorf("safe read failed: %w", err)
	}

	return n, nil
}
