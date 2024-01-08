package nats

import (
	"fmt"

	auction "github.com/Puena/auction-messages-golang"
	"github.com/nats-io/nats.go/jetstream"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func decodeDataCommandCreateProduct(data []byte) (*auction.CommandCreateProduct, error) {
	var commandCreateProduct auction.CommandCreateProduct
	if err := proto.Unmarshal(data, &commandCreateProduct); err != nil {
		return nil, fmt.Errorf("failed while unmarshalling data: %w", err)
	}

	return &commandCreateProduct, nil
}

func decodeMsgData[T any, PT interface {
	protoreflect.ProtoMessage
	*T
}](jsMsg jetstream.Msg) (*T, error) {
	var msg T
	if err := proto.Unmarshal(jsMsg.Data(), PT(&msg)); err != nil {
		return nil, fmt.Errorf("failed while unmarshalling data: %w", err)
	}

	return &msg, nil
}
