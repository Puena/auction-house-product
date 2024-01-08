package nats

import (
	"context"
	"fmt"
	"time"

	logger "github.com/Puena/auction-house-logger"
	"github.com/Puena/auction-house-product/internal/core/dto"
	"github.com/Puena/auction-house-product/internal/core/ports"
	"github.com/Puena/auction-messages-golang"
	"github.com/nats-io/nats.go/jetstream"
	"google.golang.org/grpc/codes"
)

const (
	natsConsumerMaxDeliver = 2
	natsNakDelay           = 3 * time.Second
	natsConsumerAckWait    = 4 * time.Second
)

const (
	headerAuthUserIDNotPresented = "failed while getting user id from headers"
	decodingMsgDataFailed        = "failed while decoding message data: %w"
	inProgressFailed             = "failed while setting message to in progress: %w"
	doubleAckFailed              = "failed while double ack message: %w"
	ackFailed                    = "failed while acknowledgment: %w"
	termFailed                   = "failed while term message: %w"
	publishEventErrorFailed      = "danger! failed while publishing event error"
)

var (
	errHeaderAuthUserIDNotPresented = fmt.Errorf(headerAuthUserIDNotPresented)
)

// ConsumeCommandRepositoryConfig represent the config for the consume command repository.
type ConsumeCommandRepositoryConfig struct {
	AppName                            string
	ProductStreamName                  string
	NatsHeaderAuthUserID               string
	NatsHeaderOccuredAt                string
	NatsHeaderMsgID                    string
	SubjectProductCommandCreateProduct string
	SubjectProductCommandUpdateProduct string
	SubjectProductCommandDeleteProduct string
	SubjectProductQueryFindProduct     string
	SubjectProductQueryFindProducts    string
}

// Validate validates the config.
func (c *ConsumeCommandRepositoryConfig) Validate() error {
	if c.AppName == "" {
		return fmt.Errorf("AppName can't be empty")
	}
	if c.ProductStreamName == "" {
		return fmt.Errorf("ProductStreamName can't be empty")
	}
	if c.NatsHeaderAuthUserID == "" {
		return fmt.Errorf("ProductStreamHeaderAuthUserID can't be empty")
	}
	if c.NatsHeaderOccuredAt == "" {
		return fmt.Errorf("ProductStreamHeaderOccuredAt can't be empty")
	}
	if c.SubjectProductCommandCreateProduct == "" {
		return fmt.Errorf("SubjectProductCommandCreateProduct can't be empty")
	}
	if c.SubjectProductCommandUpdateProduct == "" {
		return fmt.Errorf("SubjectProductCommandUpdateProduct can't be empty")
	}
	if c.SubjectProductCommandDeleteProduct == "" {
		return fmt.Errorf("SubjectProductCommandDeleteProduct can't be empty")
	}
	if c.SubjectProductQueryFindProduct == "" {
		return fmt.Errorf("SubjectProductQueryFindProduct can't be empty")
	}
	if c.SubjectProductQueryFindProducts == "" {
		return fmt.Errorf("SubjectProductQueryFindProducts can't be empty")
	}

	return nil
}

type productStreamConsumer struct {
	config  ConsumeCommandRepositoryConfig
	broker  jetstream.JetStream
	service ports.Service
}

// NewProductStreamConsumer creates a new consume command repository.
func NewProductStreamConsumer(nats jetstream.JetStream, service ports.Service, config ConsumeCommandRepositoryConfig) *productStreamConsumer {
	if err := config.Validate(); err != nil {
		logger.Fatal().Err(err).Msg("failed while validating config")
	}

	return &productStreamConsumer{
		config:  config,
		broker:  nats,
		service: service,
	}
}

// ConsumeCreateProductCommand consume and handle create product command.
func (r *productStreamConsumer) ConsumeCreateProductCommand(ctx context.Context) error {
	// configure jetstream consumer
	name := fmt.Sprintf("%s-consumer-command-create-product", r.config.AppName)
	filterSubject := r.config.SubjectProductCommandCreateProduct
	consumerConfig := createNewConsumerConfig(name, filterSubject)
	c, err := r.broker.CreateOrUpdateConsumer(ctx, r.config.ProductStreamName, consumerConfig)
	if err != nil {
		return err
	}

	// start consuming messages from subject
	c.Consume(func(m jetstream.Msg) {
		r.consumerHandler(ctx, consumerConfig.Name, m, r.handleCreateProductCommand)
	})

	return nil
}

func (r *productStreamConsumer) handleCreateProductCommand(ctx context.Context, consumerName string, m jetstream.Msg) *dto.ProductEventError {

	// get user id from headers, return error if auth header not presented.
	userID := m.Headers().Get(r.config.NatsHeaderAuthUserID)
	if userID == "" {
		return r.newProductEventError(consumerName, m, errHeaderAuthUserIDNotPresented, codes.PermissionDenied)
	}

	// decode message data
	s, err := decodeMsgData[auction.CommandCreateProduct](m)
	if err != nil {
		return r.newProductEventError(consumerName, m, fmt.Errorf(decodingMsgDataFailed, err), codes.InvalidArgument)
	}

	// create product
	p, err := r.service.CreateProduct(ctx, userID, mapCommandCreateProductToDto(s))
	if err != nil {
		return r.composeEventProductError(consumerName, m, err)
	}

	// publish event productCreated
	err = r.service.PublishEventProductCreated(ctx, userID, p)
	if err != nil {
		return r.composeEventProductError(consumerName, m, err)
	}

	return nil
}

// ConsumeUpdateProductCommand consume and handle update product command.
func (r *productStreamConsumer) ConsumeUpdateProductCommand(ctx context.Context) error {
	// configure jetstream consumer
	name := fmt.Sprintf("%s-consumer-command-update-product", r.config.AppName)
	filterSubject := r.config.SubjectProductCommandUpdateProduct
	consumerConfig := createNewConsumerConfig(name, filterSubject)
	c, err := r.broker.CreateOrUpdateConsumer(ctx, r.config.ProductStreamName, consumerConfig)
	if err != nil {
		return err
	}

	// start consuming messages from subject
	c.Consume(func(m jetstream.Msg) {
		r.consumerHandler(ctx, consumerConfig.Name, m, r.handleUpdateProductCommand)
	})

	return nil
}

func (r *productStreamConsumer) handleUpdateProductCommand(ctx context.Context, consumerName string, m jetstream.Msg) *dto.ProductEventError {

	// get user id from headers, return error if auth header not presented.
	userID := m.Headers().Get(r.config.NatsHeaderAuthUserID)
	if userID == "" {
		return r.newProductEventError(consumerName, m, errHeaderAuthUserIDNotPresented, codes.PermissionDenied)
	}

	// decode message data
	s, err := decodeMsgData[auction.CommandUpdateProduct](m)
	if err != nil {
		return r.newProductEventError(consumerName, m, fmt.Errorf(decodingMsgDataFailed, err), codes.InvalidArgument)
	}

	// update product
	p, err := r.service.UpdateProduct(ctx, userID, mapCommandUpdateProductToDto(s))
	if err != nil {
		return r.composeEventProductError(consumerName, m, err)
	}

	// publish event productUpdated
	err = r.service.PublishEventProductUpdated(ctx, userID, p)
	if err != nil {
		return r.composeEventProductError(consumerName, m, err)
	}

	return nil
}

// ConsumeDeleteProductCommand consume and handle delete product command.
func (r *productStreamConsumer) ConsumeDeleteProductCommand(ctx context.Context) error {
	// configure jetstream consumer
	name := fmt.Sprintf("%s-consumer-command-delete-product", r.config.AppName)
	filterSubject := r.config.SubjectProductCommandDeleteProduct
	consumerConfig := createNewConsumerConfig(name, filterSubject)
	c, err := r.broker.CreateOrUpdateConsumer(ctx, r.config.ProductStreamName, consumerConfig)
	if err != nil {
		return err
	}

	// start consuming messages from subject
	c.Consume(func(m jetstream.Msg) {
		r.consumerHandler(ctx, consumerConfig.Name, m, r.handleDeleteProductCommand)
	})

	return nil
}

func (r *productStreamConsumer) handleDeleteProductCommand(ctx context.Context, consumerName string, m jetstream.Msg) *dto.ProductEventError {

	// get user id from headers, return error if auth header not presented.
	userID := m.Headers().Get(r.config.NatsHeaderAuthUserID)
	if userID == "" {
		return r.newProductEventError(consumerName, m, errHeaderAuthUserIDNotPresented, codes.PermissionDenied)
	}

	// decode message data
	s, err := decodeMsgData[auction.CommandDeleteProduct](m)
	if err != nil {
		return r.newProductEventError(consumerName, m, fmt.Errorf(decodingMsgDataFailed, err), codes.InvalidArgument)
	}

	// delete product
	p, err := r.service.DeleteProduct(ctx, userID, mapCommandDeleteProductToDto(s))
	if err != nil {
		return r.composeEventProductError(consumerName, m, err)
	}

	// publish event productDeleted
	err = r.service.PublishEventProductDeleted(ctx, userID, p)
	if err != nil {
		return r.composeEventProductError(consumerName, m, err)
	}

	return nil
}

// ConsumeFindProductQuery consume and handle find product query.
func (r *productStreamConsumer) ConsumeFindProductQuery(ctx context.Context) error {
	// configure jetstream consumer
	name := fmt.Sprintf("%s-consumer-query-find-product", r.config.AppName)
	filterSubject := r.config.SubjectProductQueryFindProduct
	consumerConfig := createNewConsumerConfig(name, filterSubject)
	c, err := r.broker.CreateOrUpdateConsumer(ctx, r.config.ProductStreamName, consumerConfig)
	if err != nil {
		return err
	}

	// start consuming messages from subject
	c.Consume(func(m jetstream.Msg) {
		r.consumerHandler(ctx, consumerConfig.Name, m, r.handleFindProductQuery)
	})

	return nil
}

func (r *productStreamConsumer) handleFindProductQuery(ctx context.Context, consumerName string, m jetstream.Msg) *dto.ProductEventError {

	// get user id from headers, return error if auth header not presented.
	userID := m.Headers().Get(r.config.NatsHeaderAuthUserID)
	if userID == "" {
		return r.newProductEventError(consumerName, m, errHeaderAuthUserIDNotPresented, codes.PermissionDenied)
	}

	// decode message data
	s, err := decodeMsgData[auction.QueryFindProduct](m)
	if err != nil {
		return r.newProductEventError(consumerName, m, fmt.Errorf(decodingMsgDataFailed, err), codes.InvalidArgument)
	}

	// find product
	p, err := r.service.FindProduct(ctx, userID, mapQueryFindProductToDto(s))
	if err != nil {
		return r.composeEventProductError(consumerName, m, err)
	}

	// publish event productFound
	err = r.service.PublishEventProductFound(ctx, userID, p)
	if err != nil {
		return r.composeEventProductError(consumerName, m, err)
	}

	return nil
}

// ConsumeFindProductsQuery consume and handle find products query.
func (r *productStreamConsumer) ConsumeFindProductsQuery(ctx context.Context) error {
	// configure jetstream consumer
	name := fmt.Sprintf("%s-consumer-query-find-products", r.config.AppName)
	filterSubject := r.config.SubjectProductQueryFindProducts
	consumerConfig := createNewConsumerConfig(name, filterSubject)
	c, err := r.broker.CreateOrUpdateConsumer(ctx, r.config.ProductStreamName, consumerConfig)
	if err != nil {
		return err
	}

	// start consuming messages from subject
	c.Consume(func(m jetstream.Msg) {
		r.consumerHandler(ctx, consumerConfig.Name, m, r.handleFindProductsQuery)
	})

	return nil
}

func (r *productStreamConsumer) handleFindProductsQuery(ctx context.Context, consumerName string, m jetstream.Msg) *dto.ProductEventError {

	// get user id from headers, return error if auth header not presented.
	userID := m.Headers().Get(r.config.NatsHeaderAuthUserID)
	if userID == "" {
		return r.newProductEventError(consumerName, m, errHeaderAuthUserIDNotPresented, codes.PermissionDenied)
	}

	// decode message data
	s, err := decodeMsgData[auction.QueryFindProducts](m)
	if err != nil {
		return r.newProductEventError(consumerName, m, fmt.Errorf(decodingMsgDataFailed, err), codes.InvalidArgument)
	}

	// find products
	p, err := r.service.FindProducts(ctx, userID, mapQueryFindProductsToDto(s))
	if err != nil {
		return r.composeEventProductError(consumerName, m, err)
	}

	// publish event productsFound
	err = r.service.PublishEventProductsFound(ctx, userID, p)
	if err != nil {
		return r.composeEventProductError(consumerName, m, err)
	}

	return nil
}

func (r *productStreamConsumer) consumerHandler(ctx context.Context, consumerName string, m jetstream.Msg, handler func(ctx context.Context, consumerName string, m jetstream.Msg) *dto.ProductEventError) {
	// Warning! Don't forget ack message!
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	eventErr := handler(ctx, consumerName, m)
	if eventErr != nil {
		err := r.publishEventProductError(ctx, m.Headers().Get(r.config.NatsHeaderAuthUserID), eventErr)
		if err != nil {
			logger.Error().Err(err).Msg(publishEventErrorFailed)
			return
		}
	}
	err := r.NackOrAckOrTerm(ctx, m, eventErr)
	if err != nil {
		eventErr = r.newProductEventError(consumerName, m, fmt.Errorf(ackFailed, err), codes.Internal)
		err := r.publishEventProductError(ctx, m.Headers().Get(r.config.NatsHeaderAuthUserID), eventErr)
		if err != nil {
			logger.Error().Err(err).Msg(publishEventErrorFailed)
			return
		}
	}
}

func (r *productStreamConsumer) publishEventProductError(ctx context.Context, userID string, event *dto.ProductEventError) error {
	if event == nil {
		logger.Error().Msg("event can not be nil")
		return nil
	}
	logger.Error().Err(event).Msg("error occured")
	return r.service.PublishEventProductError(ctx, userID, *event)
}

func (r *productStreamConsumer) composeEventProductError(consumer string, msg jetstream.Msg, err error) *dto.ProductEventError {
	if r.service.UniqueConstrainError(err) {
		return r.newProductEventError(consumer, msg, err, codes.AlreadyExists)
	}
	if r.service.ValidationError(err) {
		return r.newProductEventError(consumer, msg, err, codes.InvalidArgument)
	}
	if r.service.NotFoundError(err) {
		return r.newProductEventError(consumer, msg, err, codes.NotFound)
	}
	return r.newProductEventError(consumer, msg, err, codes.Internal)
}

func (r *productStreamConsumer) newProductEventError(consumerName string, jsMsg jetstream.Msg, err error, status codes.Code) *dto.ProductEventError {
	return &dto.ProductEventError{
		StreamName:          r.config.ProductStreamName,
		ConsumerName:        consumerName,
		Subject:             jsMsg.Subject(),
		Reference_event_key: jsMsg.Headers().Get(r.config.NatsHeaderMsgID),
		Message:             err.Error(),
		Code:                int(status),
		Data:                jsMsg.Data(),
		Headers:             fmt.Sprintf("%+v", jsMsg.Headers()),
		Time:                time.Now(),
	}
}

func (r *productStreamConsumer) NackOrAckOrTerm(ctx context.Context, msg jetstream.Msg, dtoErr *dto.ProductEventError) error {
	if dtoErr == nil {
		return msg.DoubleAck(ctx)
	}

	switch codes.Code(dtoErr.Code) {
	case codes.OK:
		return msg.DoubleAck(ctx)
	case codes.Internal:
		return msg.NakWithDelay(natsNakDelay)
	default:
		return msg.Ack()
	}
}

func createNewConsumerConfig(name string, subject string) jetstream.ConsumerConfig {
	return jetstream.ConsumerConfig{
		Name:          name,
		FilterSubject: subject,
		AckWait:       natsConsumerAckWait,
		MaxDeliver:    natsConsumerMaxDeliver,
	}
}
