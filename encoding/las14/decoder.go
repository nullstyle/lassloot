package las14

import (
	"fmt"
	"io"
	"sync"
)

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
}

func (las *Decoder) FirstPassDecode() (error, *FirstPassResult) {
	las.mt.Lock()
	defer las.mt.Unlock()

	return las.firstPassDecode()
}

func (las *Decoder) FullDecode() (error, *FullResult) {
	las.mt.Lock()
	defer las.mt.Unlock()

	return las.fullDecode()
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
		return fmt.Errorf("failed to read full las header: only %i bytes read", n), nil
	}

	if (string)(actualSig) != HeaderMagicBytes {
		return fmt.Errorf("invalid file signature: read %s", actualSig), nil
	}

	las.fp = &FirstPassResult{
		Header: header,
	}

	return nil, las.fp
}

func (las *Decoder) fullDecode() (error, *FullResult) {
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
	// TODO

	return nil, &FullResult{
		*fp,
	}
}

func (las *Decoder) safeRead(p []byte) (n int, err error) {
	requested := uint(len(p))
	if requested > las.budget {
		return 0, fmt.Errorf("read budget exhausted: %i byte read requested", requested)
	}
	// NOTE: we deduct the requested byte count rather than the actually-read byte count from the budget because we are crotchety bastards #dealwithit
	las.budget -= requested

	n, err = las.r.Read(p)
	if err != nil {
		return n, fmt.Errorf("safe read failed: %w", err)
	}

	return n, nil
}
