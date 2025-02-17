package deletesvc

import (
	objectV2 "github.com/nspcc-dev/neofs-api-go/v2/object"
	deletesvc "github.com/nspcc-dev/neofs-node/pkg/services/object/delete"
	"github.com/nspcc-dev/neofs-node/pkg/services/object/util"
	"github.com/nspcc-dev/neofs-sdk-go/object"
)

type tombstoneBodyWriter struct {
	body *objectV2.DeleteResponseBody
}

func (s *Service) toPrm(req *objectV2.DeleteRequest, respBody *objectV2.DeleteResponseBody) (*deletesvc.Prm, error) {
	commonPrm, err := util.CommonPrmFromV2(req)
	if err != nil {
		return nil, err
	}

	p := new(deletesvc.Prm)
	p.SetCommonParameters(commonPrm)

	body := req.GetBody()
	p.WithAddress(object.NewAddressFromV2(body.GetAddress()))
	p.WithTombstoneAddressTarget(&tombstoneBodyWriter{
		body: respBody,
	})

	return p, nil
}

func (w *tombstoneBodyWriter) SetAddress(addr *object.Address) {
	w.body.SetTombstone(addr.ToV2())
}
