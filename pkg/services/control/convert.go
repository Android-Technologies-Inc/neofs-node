package control

import (
	"github.com/nspcc-dev/neofs-api-go/v2/rpc/grpc"
	"github.com/nspcc-dev/neofs-api-go/v2/rpc/message"
)

type requestWrapper struct {
	message.Message
	m grpc.Message
}

func (w *requestWrapper) ToGRPCMessage() grpc.Message {
	return w.m
}

type healthCheckResponseWrapper struct {
	m *HealthCheckResponse
}

func (w *healthCheckResponseWrapper) ToGRPCMessage() grpc.Message {
	return w.m
}

func (w *healthCheckResponseWrapper) FromGRPCMessage(m grpc.Message) error {
	var ok bool

	w.m, ok = m.(*HealthCheckResponse)
	if !ok {
		return message.NewUnexpectedMessageType(m, w.m)
	}

	return nil
}

type netmapSnapshotResponseWrapper struct {
	message.Message
	m *NetmapSnapshotResponse
}

func (w *netmapSnapshotResponseWrapper) ToGRPCMessage() grpc.Message {
	return w.m
}

func (w *netmapSnapshotResponseWrapper) FromGRPCMessage(m grpc.Message) error {
	var ok bool

	w.m, ok = m.(*NetmapSnapshotResponse)
	if !ok {
		return message.NewUnexpectedMessageType(m, w.m)
	}

	return nil
}

type setNetmapStatusResponseWrapper struct {
	message.Message
	m *SetNetmapStatusResponse
}

func (w *setNetmapStatusResponseWrapper) ToGRPCMessage() grpc.Message {
	return w.m
}

func (w *setNetmapStatusResponseWrapper) FromGRPCMessage(m grpc.Message) error {
	var ok bool

	w.m, ok = m.(*SetNetmapStatusResponse)
	if !ok {
		return message.NewUnexpectedMessageType(m, w.m)
	}

	return nil
}

type dropObjectsResponseWrapper struct {
	message.Message
	m *DropObjectsResponse
}

func (w *dropObjectsResponseWrapper) ToGRPCMessage() grpc.Message {
	return w.m
}

func (w *dropObjectsResponseWrapper) FromGRPCMessage(m grpc.Message) error {
	var ok bool

	w.m, ok = m.(*DropObjectsResponse)
	if !ok {
		return message.NewUnexpectedMessageType(m, w.m)
	}

	return nil
}

type listShardsResponseWrapper struct {
	m *ListShardsResponse
}

func (w *listShardsResponseWrapper) ToGRPCMessage() grpc.Message {
	return w.m
}

func (w *listShardsResponseWrapper) FromGRPCMessage(m grpc.Message) error {
	var ok bool

	w.m, ok = m.(*ListShardsResponse)
	if !ok {
		return message.NewUnexpectedMessageType(m, w.m)
	}

	return nil
}

type setShardModeResponseWrapper struct {
	m *SetShardModeResponse
}

func (w *setShardModeResponseWrapper) ToGRPCMessage() grpc.Message {
	return w.m
}

func (w *setShardModeResponseWrapper) FromGRPCMessage(m grpc.Message) error {
	var ok bool

	w.m, ok = m.(*SetShardModeResponse)
	if !ok {
		return message.NewUnexpectedMessageType(m, w.m)
	}

	return nil
}

type dumpShardResponseWrapper struct {
	*DumpShardResponse
}

func (w *dumpShardResponseWrapper) ToGRPCMessage() grpc.Message {
	return w.DumpShardResponse
}

func (w *dumpShardResponseWrapper) FromGRPCMessage(m grpc.Message) error {
	r, ok := m.(*DumpShardResponse)
	if !ok {
		return message.NewUnexpectedMessageType(m, (*DumpShardResponse)(nil))
	}

	w.DumpShardResponse = r
	return nil
}

type restoreShardResponseWrapper struct {
	*RestoreShardResponse
}

func (w *restoreShardResponseWrapper) ToGRPCMessage() grpc.Message {
	return w.RestoreShardResponse
}

func (w *restoreShardResponseWrapper) FromGRPCMessage(m grpc.Message) error {
	r, ok := m.(*RestoreShardResponse)
	if !ok {
		return message.NewUnexpectedMessageType(m, (*RestoreShardResponse)(nil))
	}

	w.RestoreShardResponse = r
	return nil
}
