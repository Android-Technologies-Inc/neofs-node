package blobstor

import (
	"errors"

	"github.com/nspcc-dev/neofs-node/pkg/local_object_storage/blobstor/fstree"
	"github.com/nspcc-dev/neofs-sdk-go/object"
)

// ExistsPrm groups the parameters of Exists operation.
type ExistsPrm struct {
	address
}

// ExistsRes groups resulting values of Exists operation.
type ExistsRes struct {
	exists bool
}

// Exists returns the fact that the object is in BLOB storage.
func (r ExistsRes) Exists() bool {
	return r.exists
}

// Exists checks if object is presented in BLOB storage.
//
// Returns any error encountered that did not allow
// to completely check object existence.
func (b *BlobStor) Exists(prm *ExistsPrm) (*ExistsRes, error) {
	// check presence in shallow dir first (cheaper)
	exists, err := b.existsBig(prm.addr)
	if !exists {
		// TODO: do smth if err != nil

		// check presence in blobovnicza
		exists, err = b.existsSmall(prm.addr)
	}

	if err != nil {
		return nil, err
	}

	return &ExistsRes{
		exists: exists,
	}, err
}

// checks if object is presented in shallow dir.
func (b *BlobStor) existsBig(addr *object.Address) (bool, error) {
	_, err := b.fsTree.Exists(addr)
	if errors.Is(err, fstree.ErrFileNotFound) {
		return false, nil
	}

	return err == nil, err
}

// checks if object is presented in blobovnicza.
func (b *BlobStor) existsSmall(addr *object.Address) (bool, error) {
	// TODO: implement
	return false, nil
}
