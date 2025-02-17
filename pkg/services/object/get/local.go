package getsvc

import (
	"errors"

	"github.com/nspcc-dev/neofs-node/pkg/core/object"
	objectSDK "github.com/nspcc-dev/neofs-sdk-go/object"
	"go.uber.org/zap"
)

func (exec *execCtx) executeLocal() {
	var err error

	exec.collectedObject, err = exec.svc.localStorage.get(exec)

	var errSplitInfo *objectSDK.SplitInfoError

	switch {
	default:
		exec.status = statusUndefined
		exec.err = err

		exec.log.Debug("local get failed",
			zap.String("error", err.Error()),
		)
	case err == nil:
		exec.status = statusOK
		exec.err = nil
		exec.writeCollectedObject()
	case errors.Is(err, object.ErrAlreadyRemoved):
		exec.status = statusINHUMED
		exec.err = object.ErrAlreadyRemoved
	case errors.As(err, &errSplitInfo):
		exec.status = statusVIRTUAL
		mergeSplitInfo(exec.splitInfo(), errSplitInfo.SplitInfo())
		exec.err = objectSDK.NewSplitInfoError(exec.infoSplit)
	case errors.Is(err, object.ErrRangeOutOfBounds):
		exec.status = statusOutOfRange
		exec.err = object.ErrRangeOutOfBounds
	}
}
