package transformer

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/nspcc-dev/neofs-node/pkg/core/netmap"
	"github.com/nspcc-dev/neofs-node/pkg/core/object"
	objectSDK "github.com/nspcc-dev/neofs-sdk-go/object"
	"github.com/nspcc-dev/neofs-sdk-go/session"
	"github.com/nspcc-dev/neofs-sdk-go/version"
)

type formatter struct {
	prm *FormatterParams

	obj *object.RawObject

	sz uint64
}

// FormatterParams groups NewFormatTarget parameters.
type FormatterParams struct {
	Key *ecdsa.PrivateKey

	NextTarget ObjectTarget

	SessionToken *session.Token

	NetworkState netmap.State
}

// NewFormatTarget returns ObjectTarget instance that finalizes object structure
// and writes it to the next target.
//
// Chunks must be written before the WriteHeader call.
//
// Object changes:
// - sets version to current SDK version;
// - sets payload size to the total length of all written chunks;
// - sets session token;
// - sets number of creation epoch;
// - calculates and sets verification fields (ID, Signature).
func NewFormatTarget(p *FormatterParams) ObjectTarget {
	return &formatter{
		prm: p,
	}
}

func (f *formatter) WriteHeader(obj *object.RawObject) error {
	f.obj = obj

	return nil
}

func (f *formatter) Write(p []byte) (n int, err error) {
	n, err = f.prm.NextTarget.Write(p)

	f.sz += uint64(n)

	return
}

func (f *formatter) Close() (*AccessIdentifiers, error) {
	curEpoch := f.prm.NetworkState.CurrentEpoch()

	f.obj.SetVersion(version.Current())
	f.obj.SetPayloadSize(f.sz)
	f.obj.SetSessionToken(f.prm.SessionToken)
	f.obj.SetCreationEpoch(curEpoch)

	var (
		parID  *objectSDK.ID
		parHdr *objectSDK.Object
	)

	if par := f.obj.Parent(); par != nil && par.Signature() == nil {
		rawPar := objectSDK.NewRawFromV2(par.ToV2())

		rawPar.SetSessionToken(f.prm.SessionToken)
		rawPar.SetCreationEpoch(curEpoch)

		if err := objectSDK.SetIDWithSignature(f.prm.Key, rawPar); err != nil {
			return nil, fmt.Errorf("could not finalize parent object: %w", err)
		}

		parID = rawPar.ID()
		parHdr = rawPar.Object()

		f.obj.SetParent(parHdr)
	}

	if err := objectSDK.SetIDWithSignature(f.prm.Key, f.obj.SDK()); err != nil {
		return nil, fmt.Errorf("could not finalize object: %w", err)
	}

	if err := f.prm.NextTarget.WriteHeader(f.obj); err != nil {
		return nil, fmt.Errorf("could not write header to next target: %w", err)
	}

	if _, err := f.prm.NextTarget.Close(); err != nil {
		return nil, fmt.Errorf("could not close next target: %w", err)
	}

	return new(AccessIdentifiers).
		WithSelfID(f.obj.ID()).
		WithParentID(parID).
		WithParent(parHdr), nil
}
