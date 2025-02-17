package audit

import (
	"fmt"

	"github.com/nspcc-dev/neo-go/pkg/vm/stackitem"
	"github.com/nspcc-dev/neofs-node/pkg/morph/client"
)

// ListResultsArgs groups the arguments
// of "list all audit result IDs" test invoke call.
type ListResultsArgs struct{}

// ListResultsByEpochArgs groups the arguments
// of "list audit result IDs by epoch" test invoke call.
type ListResultsByEpochArgs struct {
	epoch int64
}

// ListResultsByCIDArgs groups the arguments
// of "list audit result IDs by epoch and CID" test invoke call.
type ListResultsByCIDArgs struct {
	ListResultsByEpochArgs

	cid []byte
}

// ListResultsByNodeArgs groups the arguments
// of "list audit result IDs by epoch, CID, and node key" test invoke call.
type ListResultsByNodeArgs struct {
	ListResultsByCIDArgs

	nodeKey []byte
}

// ListResultsValues groups the stack parameters
// returned by "list audit results" test invoke.
type ListResultsValues struct {
	rawResults [][]byte // audit results in a binary format
}

// RawResults returns list of audit result IDs
// in a binary format.
func (v *ListResultsValues) RawResults() [][]byte {
	return v.rawResults
}

// SetEpoch sets epoch of listing audit results.
func (v *ListResultsByEpochArgs) SetEpoch(epoch int64) {
	v.epoch = epoch
}

// SetCID sets container ID of listing audit results.
func (v *ListResultsByCIDArgs) SetCID(cid []byte) {
	v.cid = cid
}

// SetNodeKey sets public key of node that produced listing audit results.
func (v *ListResultsByNodeArgs) SetNodeKey(key []byte) {
	v.nodeKey = key
}

// ListAuditResults performs the test invoke of "list all audit result IDs"
// method of NeoFS Audit contract.
func (c *Client) ListAuditResults(args ListResultsArgs) (*ListResultsValues, error) {
	invokePrm := client.TestInvokePrm{}

	invokePrm.SetMethod(listResultsMethod)

	items, err := c.client.TestInvoke(invokePrm)
	if err != nil {
		return nil, fmt.Errorf("could not perform test invocation (%s): %w", listResultsMethod, err)
	}

	return parseAuditResults(items, listResultsMethod)
}

// ListAuditResultsByEpoch performs the test invoke of "list audit result IDs
// by epoch" method of NeoFS Audit contract.
func (c *Client) ListAuditResultsByEpoch(args ListResultsByEpochArgs) (*ListResultsValues, error) {
	invokePrm := client.TestInvokePrm{}

	invokePrm.SetMethod(listByEpochResultsMethod)
	invokePrm.SetArgs(args.epoch)

	items, err := c.client.TestInvoke(invokePrm)
	if err != nil {
		return nil, fmt.Errorf("could not perform test invocation (%s): %w", listByEpochResultsMethod, err)
	}

	return parseAuditResults(items, listByEpochResultsMethod)
}

// ListAuditResultsByCID performs the test invoke of "list audit result IDs
// by epoch and CID" method of NeoFS Audit contract.
func (c *Client) ListAuditResultsByCID(args ListResultsByCIDArgs) (*ListResultsValues, error) {
	invokePrm := client.TestInvokePrm{}

	invokePrm.SetMethod(listByCIDResultsMethod)
	invokePrm.SetArgs(args.epoch, args.cid)

	items, err := c.client.TestInvoke(invokePrm)
	if err != nil {
		return nil, fmt.Errorf("could not perform test invocation (%s): %w", listByCIDResultsMethod, err)
	}

	return parseAuditResults(items, listByCIDResultsMethod)
}

// ListAuditResultsByNode performs the test invoke of "list audit result IDs
// by epoch, CID, and node key" method of NeoFS Audit contract.
func (c *Client) ListAuditResultsByNode(args ListResultsByNodeArgs) (*ListResultsValues, error) {
	invokePrm := client.TestInvokePrm{}

	invokePrm.SetMethod(listByNodeResultsMethod)
	invokePrm.SetArgs(args.epoch, args.cid, args.nodeKey)

	items, err := c.client.TestInvoke(invokePrm)
	if err != nil {
		return nil, fmt.Errorf("could not perform test invocation (%s): %w", listByNodeResultsMethod, err)
	}

	return parseAuditResults(items, listByNodeResultsMethod)
}

func parseAuditResults(items []stackitem.Item, method string) (*ListResultsValues, error) {
	if ln := len(items); ln != 1 {
		return nil, fmt.Errorf("unexpected stack item count (%s): %d", method, ln)
	}

	items, err := client.ArrayFromStackItem(items[0])
	if err != nil {
		return nil, fmt.Errorf("could not get stack item array from stack item (%s): %w", method, err)
	}

	res := &ListResultsValues{
		rawResults: make([][]byte, 0, len(items)),
	}

	for i := range items {
		rawRes, err := client.BytesFromStackItem(items[i])
		if err != nil {
			return nil, fmt.Errorf("could not get byte array from stack item (%s): %w", method, err)
		}

		res.rawResults = append(res.rawResults, rawRes)
	}

	return res, nil
}
