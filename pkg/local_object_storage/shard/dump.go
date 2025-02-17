package shard

import (
	"encoding/binary"
	"errors"
	"io"
	"os"

	"github.com/nspcc-dev/neofs-node/pkg/local_object_storage/blobstor"
	"github.com/nspcc-dev/neofs-node/pkg/local_object_storage/writecache"
)

var dumpMagic = []byte("NEOF")

// DumpPrm groups the parameters of Dump operation.
type DumpPrm struct {
	path         string
	stream       io.Writer
	ignoreErrors bool
}

// WithPath is an Dump option to set the destination path.
func (p *DumpPrm) WithPath(path string) *DumpPrm {
	p.path = path
	return p
}

// WithStream is an Dump option to set the destination stream.
// It takes priority over `path` option.
func (p *DumpPrm) WithStream(r io.Writer) *DumpPrm {
	p.stream = r
	return p
}

// WithIgnoreErrors is an Dump option to allow ignore all errors during iteration.
// This includes invalid blobovniczas as well as corrupted objects.
func (p *DumpPrm) WithIgnoreErrors(ignore bool) *DumpPrm {
	p.ignoreErrors = ignore
	return p
}

// DumpRes groups the result fields of Dump operation.
type DumpRes struct {
	count int
}

// Count return amount of object written.
func (r *DumpRes) Count() int {
	return r.count
}

var ErrMustBeReadOnly = errors.New("shard must be in read-only mode")

// Dump dumps all objects from the shard to a file or stream.
//
// Returns any error encountered.
func (s *Shard) Dump(prm *DumpPrm) (*DumpRes, error) {
	s.m.RLock()
	defer s.m.RUnlock()

	if s.info.Mode != ModeReadOnly {
		return nil, ErrMustBeReadOnly
	}

	w := prm.stream
	if w == nil {
		f, err := os.OpenFile(prm.path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0640)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		w = f
	}

	_, err := w.Write(dumpMagic)
	if err != nil {
		return nil, err
	}

	var count int

	if s.hasWriteCache() {
		err := s.writeCache.Iterate(new(writecache.IterationPrm).WithHandler(func(data []byte) error {
			var size [4]byte
			binary.LittleEndian.PutUint32(size[:], uint32(len(data)))
			if _, err := w.Write(size[:]); err != nil {
				return err
			}

			if _, err := w.Write(data); err != nil {
				return err
			}

			count++
			return nil
		}).WithIgnoreErrors(prm.ignoreErrors))
		if err != nil {
			return nil, err
		}
	}

	var pi blobstor.IteratePrm

	if prm.ignoreErrors {
		pi.IgnoreErrors()
	}
	pi.SetIterationHandler(func(elem blobstor.IterationElement) error {
		data := elem.ObjectData()

		var size [4]byte
		binary.LittleEndian.PutUint32(size[:], uint32(len(data)))
		if _, err := w.Write(size[:]); err != nil {
			return err
		}

		if _, err := w.Write(data); err != nil {
			return err
		}

		count++
		return nil
	})

	if _, err := s.blobStor.Iterate(pi); err != nil {
		return nil, err
	}

	return &DumpRes{count: count}, nil
}
