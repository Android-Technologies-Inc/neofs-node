package shard

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"os"

	"github.com/nspcc-dev/neofs-node/pkg/core/object"
)

// ErrInvalidMagic is returned when dump format is invalid.
var ErrInvalidMagic = errors.New("invalid magic")

// RestorePrm groups the parameters of Restore operation.
type RestorePrm struct {
	path         string
	stream       io.Reader
	ignoreErrors bool
}

// WithPath is a Restore option to set the destination path.
func (p *RestorePrm) WithPath(path string) *RestorePrm {
	p.path = path
	return p
}

// WithStream is a Restore option to set the stream to read objects from.
// It takes priority over `WithPath` option.
func (p *RestorePrm) WithStream(r io.Reader) *RestorePrm {
	p.stream = r
	return p
}

// WithIgnoreErrors is a Restore option which allows to ignore errors encountered during restore.
// Corrupted objects will not be processed.
func (p *RestorePrm) WithIgnoreErrors(ignore bool) *RestorePrm {
	p.ignoreErrors = ignore
	return p
}

// RestoreRes groups the result fields of Restore operation.
type RestoreRes struct {
	count  int
	failed int
}

// Count return amount of object written.
func (r *RestoreRes) Count() int {
	return r.count
}

// FailCount return amount of object skipped.
func (r *RestoreRes) FailCount() int {
	return r.failed
}

// Restore restores objects from the dump prepared by Dump.
//
// Returns any error encountered.
func (s *Shard) Restore(prm *RestorePrm) (*RestoreRes, error) {
	// Disallow changing mode during restore.
	s.m.RLock()
	defer s.m.RUnlock()

	if s.info.Mode != ModeReadWrite {
		return nil, ErrReadOnlyMode
	}

	r := prm.stream
	if r == nil {
		f, err := os.OpenFile(prm.path, os.O_RDONLY, os.ModeExclusive)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		r = f
	}

	var m [4]byte
	_, _ = io.ReadFull(r, m[:])
	if !bytes.Equal(m[:], dumpMagic) {
		return nil, ErrInvalidMagic
	}

	var count, failCount int
	var data []byte
	var size [4]byte
	for {
		// If there are less than 4 bytes left, `Read` returns nil error instead of
		// io.ErrUnexpectedEOF, thus `ReadFull` is used.
		_, err := io.ReadFull(r, size[:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}

		sz := binary.LittleEndian.Uint32(size[:])
		if uint32(cap(data)) < sz {
			data = make([]byte, sz)
		} else {
			data = data[:sz]
		}

		_, err = r.Read(data)
		if err != nil {
			return nil, err
		}

		obj := object.New()
		err = obj.Unmarshal(data)
		if err != nil {
			if prm.ignoreErrors {
				failCount++
				continue
			}
			return nil, err
		}

		_, err = s.Put(new(PutPrm).WithObject(obj))
		if err != nil {
			return nil, err
		}

		count++
	}

	return &RestoreRes{count: count, failed: failCount}, nil
}
