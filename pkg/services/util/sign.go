package util

import (
	"context"
	"crypto/ecdsa"

	"github.com/nspcc-dev/neofs-api-go/v2/signature"
	"github.com/pkg/errors"
)

type UnaryHandler func(context.Context, interface{}) (interface{}, error)

type UnarySignService struct {
	key *ecdsa.PrivateKey
}

func NewUnarySignService(key *ecdsa.PrivateKey) *UnarySignService {
	return &UnarySignService{
		key: key,
	}
}

func (s *UnarySignService) HandleServerStreamRequest(ctx context.Context, req interface{}, handler UnaryHandler) (interface{}, error) {
	return s.verifyAndProc(ctx, req, handler)
}

func (s *UnarySignService) HandleUnaryRequest(ctx context.Context, req interface{}, handler UnaryHandler) (interface{}, error) {
	// verify and process request
	resp, err := s.verifyAndProc(ctx, req, handler)
	if err != nil {
		return nil, err
	}

	// sign the response
	if err := signature.SignServiceMessage(s.key, resp); err != nil {
		return nil, errors.Wrap(err, "could not sign response")
	}

	return resp, nil
}

func (s *UnarySignService) verifyAndProc(ctx context.Context, req interface{}, handler UnaryHandler) (interface{}, error) {
	// verify request signatures
	if err := signature.VerifyServiceMessage(req); err != nil {
		return nil, errors.Wrap(err, "could not verify request")
	}

	// process request
	resp, err := handler(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "could not handle request")
	}

	return resp, nil
}
