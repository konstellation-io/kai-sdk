package sdk

import (
	"os"

	centralizedConfiguration "github.com/konstellation-io/kai-sdk/go-sdk/v1/sdk/centralized-configuration"
	pathutils "github.com/konstellation-io/kai-sdk/go-sdk/v1/sdk/path-utils"

	objectstore "github.com/konstellation-io/kai-sdk/go-sdk/v1/sdk/object-store"

	"github.com/go-logr/logr"
	kai "github.com/konstellation-io/kai-sdk/go-sdk/v1/protos"
	msg "github.com/konstellation-io/kai-sdk/go-sdk/v1/sdk/messaging"
	meta "github.com/konstellation-io/kai-sdk/go-sdk/v1/sdk/metadata"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

//go:generate mockery --name pathUtils --output ../mocks --filename path_utils_mock.go --structname PathUtilsMock
type pathUtils interface {
	GetBasePath() string
	ComposePath(relativePath ...string) string
}

//go:generate mockery --name messaging --output ../mocks --filename messaging_mock.go --structname MessagingMock
type messaging interface {
	SendOutput(response proto.Message, channelOpt ...string) error
	SendOutputWithRequestID(response proto.Message, requestID string, channelOpt ...string) error
	SendAny(response *anypb.Any, channelOpt ...string)
	SendEarlyReply(response proto.Message, channelOpt ...string) error
	SendEarlyExit(response proto.Message, channelOpt ...string) error

	IsMessageOK() bool
	IsMessageError() bool
	IsMessageEarlyReply() bool
	IsMessageEarlyExit() bool
}

//go:generate mockery --name metadata --output ../mocks --filename metadata_mock.go --structname MetadataMock
type metadata interface {
	GetProcess() string
	GetWorkflow() string
	GetProduct() string
	GetVersion() string
	GetObjectStoreName() string
	GetKeyValueStoreProductName() string
	GetKeyValueStoreWorkflowName() string
	GetKeyValueStoreProcessName() string
}

//go:generate mockery --name objectStore --output ../mocks --filename object_store_mock.go --structname ObjectStoreMock
type objectStore interface {
	List(regexp ...string) ([]string, error)
	Get(key string) ([]byte, error)
	Save(key string, value []byte) error
	Delete(key string) error
	Purge(regexp ...string) error
}

//go:generate mockery --name centralizedConfig --output ../mocks --filename centralized_config_mock.go --structname CentralizedConfigMock
type centralizedConfig interface {
	GetConfig(key string, scope ...msg.Scope) (string, error)
	SetConfig(key, value string, scope ...msg.Scope) error
	DeleteConfig(key string, scope msg.Scope) error
}

// TODO add metrics interface
//
//go:generate mockery --name measurements --output ../mocks --filename measurements_mock.go --structname MeasurementsMock
type measurements interface{}

// TODO add storage interface
//
//go:generate mockery --name storage --output ../mocks --filename storage_mock.go --structname StorageMock
type storage interface{}

type KaiSDK struct {
	// Needed deps
	nats           *nats.Conn
	jetstream      nats.JetStreamContext
	requestMessage *kai.KaiNatsMessage

	// Main methods
	Logger            logr.Logger
	PathUtils         pathUtils
	Metadata          metadata
	Messaging         messaging
	ObjectStore       objectStore
	CentralizedConfig centralizedConfig
	Measurements      measurements
	Storage           storage
}

func NewKaiSDK(logger logr.Logger, natsCli *nats.Conn, jetstreamCli nats.JetStreamContext) KaiSDK {
	logger = logger.WithName("[KAI SDK]")

	centralizedConfigInst, err := centralizedConfiguration.NewCentralizedConfiguration(logger, jetstreamCli)
	if err != nil {
		logger.Error(err, "Error initializing Centralized Configuration")
		os.Exit(1)
	}

	objectStoreInst, err := objectstore.NewObjectStore(logger, jetstreamCli)
	if err != nil {
		logger.Error(err, "Error initializing Object Store")
		os.Exit(1)
	}

	messagingInst := msg.NewMessaging(logger, natsCli, jetstreamCli, nil)

	sdk := KaiSDK{
		nats:              natsCli,
		jetstream:         jetstreamCli,
		Logger:            logger,
		PathUtils:         pathutils.NewPathUtils(logger),
		Metadata:          meta.NewMetadata(logger),
		Messaging:         messagingInst,
		ObjectStore:       objectStoreInst,
		CentralizedConfig: centralizedConfigInst,
		Measurements:      nil,
		Storage:           nil,
	}

	return sdk
}

func (sdk *KaiSDK) GetRequestID() string {
	return sdk.requestMessage.RequestId
}

func ShallowCopyWithRequest(sdk *KaiSDK, requestMsg *kai.KaiNatsMessage) KaiSDK {
	hSdk := *sdk
	hSdk.requestMessage = requestMsg
	hSdk.Messaging = msg.NewMessaging(sdk.Logger, sdk.nats, sdk.jetstream, requestMsg)

	return hSdk
}
