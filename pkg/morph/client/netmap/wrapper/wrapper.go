package wrapper

import (
	"fmt"

	"github.com/nspcc-dev/neo-go/pkg/encoding/fixedn"
	"github.com/nspcc-dev/neo-go/pkg/util"
	"github.com/nspcc-dev/neofs-node/pkg/morph/client"
	"github.com/nspcc-dev/neofs-node/pkg/morph/client/internal"
	"github.com/nspcc-dev/neofs-node/pkg/morph/client/netmap"
)

// Wrapper is a wrapper over netmap contract
// client which implements:
//  * network map storage;
//  * tool for peer state updating.
//
// Working wrapper must be created via constructor New.
// Using the Wrapper that has been created with new(Wrapper)
// expression (or just declaring a Wrapper variable) is unsafe
// and can lead to panic.
type Wrapper struct {
	internal.StaticClient

	client *netmap.Client
}

// Option allows to set an optional
// parameter of Wrapper.
type Option func(*opts)

type opts []client.StaticClientOption

func defaultOpts() *opts {
	return new(opts)
}

// NewFromMorph returns the wrapper instance from the raw morph client.
func NewFromMorph(cli *client.Client, contract util.Uint160, fee fixedn.Fixed8, opts ...Option) (*Wrapper, error) {
	o := defaultOpts()

	for i := range opts {
		opts[i](o)
	}

	staticClient, err := client.NewStatic(cli, contract, fee, ([]client.StaticClientOption)(*o)...)
	if err != nil {
		return nil, fmt.Errorf("can't create netmap static client: %w", err)
	}

	return &Wrapper{
		StaticClient: staticClient,
		client:       netmap.New(staticClient),
	}, nil
}

// Morph returns raw morph client.
func (w Wrapper) Morph() *client.Client {
	return w.client.Morph()
}

// TryNotary returns option to enable
// notary invocation tries.
func TryNotary() Option {
	return func(o *opts) {
		*o = append(*o, client.TryNotary())
	}
}

// AsAlphabet returns option to sign main TX
// of notary requests with client's private
// key.
//
// Considered to be used by IR nodes only.
func AsAlphabet() Option {
	return func(o *opts) {
		*o = append(*o, client.AsAlphabet())
	}
}
