package config

type ConfigOpt func(*config)

type config struct {
	// AppName is name of current application.
	AppName string `env:"APP_NAME" envDefault:"product-service"`
	// NatsURL is url of nats server.
	NatsURL string `env:"NATS_URL" envDefault:"nats://localhost:4222"`
	// PostgresDSN is connection string of postgres database.
	PostgresDSN string `env:"POSTGRES_DSN" envDefault:"postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"`
	// NameProductStream is name of product stream.
	NameProductStream string `env:"PRODUCT_STREAM_NAME" envDefault:"PRODUCT"`
	NameErrorStream   string `env:"ERROR_STREAM_NAME" envDefault:"ERROR"`
	// NatsHeaderAuthUserID is header name of auth user id.
	NatsHeaderAuthUserID string `env:"PRODUCT_STREAM_HEADER_AUTH_USER_ID" envDefault:"Auth-User-Id"`
	// NatsHeaderOccuredAt is header name of occured at.
	NatsHeaderOccuredAt string `env:"PRODUCT_STREAM_HEADER_OCCURED_AT" envDefault:"Msg-Occured-At"`
	// NatsHeaderMsgID is header name of msg id.
	NatsHeaderMsgID string `env:"PRODUCT_STREAM_HEADER_OCCURED_AT" envDefault:"Nats-Msg-Id"`
	// Names of subject events of product stream.
	SubjectEventProductCreated string `env:"SUBJECT_EVENT_PRODUCT_CREATED" envDefault:"product.event.product_created"`
	SubjectEventProductUpdated string `env:"SUBJECT_EVENT_PRODUCT_UPDATED" envDefault:"product.event.product_updated"`
	SubjectEventProductDeleted string `env:"SUBJECT_EVENT_PRODUCT_DELETED" envDefault:"product.event.product_deleted"`
	SubjectEventProductFound   string `env:"SUBJECT_EVENT_PRODUCT_FOUND" envDefault:"product.event.product_found"`
	SubjectEventProductsFound  string `env:"SUBJECT_EVENT_PRODUCTS_FOUND" envDefault:"product.event.products_found"`
	SubjectEventProductError   string `env:"SUBJECT_EVENT_PRODUCT_ERROR" envDefault:"error.product"`
	// Names of subject commands of product stream.
	SubjectCommandCreateProduct string `env:"SUBJECT_COMMAND_CREATE_PRODUCT" envDefault:"product.command.create_product"`
	SubjectCommandUpdateProduct string `env:"SUBJECT_COMMAND_UPDATE_PRODUCT" envDefault:"product.command.update_product"`
	SubjectCommandDeleteProduct string `env:"SUBJECT_COMMAND_DELETE_PRODUCT" envDefault:"product.command.delete_product"`
	// Names of subject query of product stream.
	SubjectQueryFindProduct  string `env:"SUBJECT_QUERY_FIND_PRODUCT" envDefault:"product.query.find_product"`
	SubjectQueryFindProducts string `env:"SUBJECT_QUERY_FIND_PRODUCTS" envDefault:"product.query.find_products"`
	// LogLevel is level of log.
	LogLevel string `env:"LOG_LEVEL" envDefault:"debug"`
	// TraceEnabled is a flag to enable tracing.
	TraceEnabled bool `env:"TRACE_ENABLED" envDefault:"true"`
}

func NewConfig(opts ...ConfigOpt) (*config, error) {
	cfg := &config{}
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg, nil
}
