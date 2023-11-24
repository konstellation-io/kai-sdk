package sdk

import (
	"context"
	"os"

	centralizedconfiguration "github.com/konstellation-io/kai-sdk/go-sdk/sdk/centralized-configuration"
	objectstore "github.com/konstellation-io/kai-sdk/go-sdk/sdk/ephemeral-storage"
	modelregistry "github.com/konstellation-io/kai-sdk/go-sdk/sdk/model-registry"
	persistentstorage "github.com/konstellation-io/kai-sdk/go-sdk/sdk/persistent-storage"

	"github.com/go-logr/logr"
	kai "github.com/konstellation-io/kai-sdk/go-sdk/protos"
	msg "github.com/konstellation-io/kai-sdk/go-sdk/sdk/messaging"
	meta "github.com/konstellation-io/kai-sdk/go-sdk/sdk/metadata"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

const (
	LoggerRequestID = "request_id"
)

//go:generate mockery --name messaging --output ../mocks --filename messaging_mock.go --structname MessagingMock
type messaging interface {
	SendOutput(response proto.Message, channelOpt ...string) error
	SendOutputWithRequestID(response proto.Message, requestID string, channelOpt ...string) error
	SendAny(response *anypb.Any, channelOpt ...string)
	SendEarlyReply(response proto.Message, channelOpt ...string) error
	SendEarlyExit(response proto.Message, channelOpt ...string) error
	GetErrorMessage() string
	GetRequestID(msg *nats.Msg) (string, error)

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
	GetEphemeralStorageName() string
	GetGlobalCentralizedConfigurationName() string
	GetProductCentralizedConfigurationName() string
	GetWorkflowCentralizedConfigurationName() string
	GetProcessCentralizedConfigurationName() string
}

type Storage struct {
	Ephemeral  ephemeralStorage
	Persistent persistentStorage
}

//go:generate mockery --name ephemeralStorage --output ../mocks --filename ephemeral_storage_mock.go --structname EphemeralStorageMock
type ephemeralStorage interface {
	Save(key string, value []byte, overwrite ...bool) error
	Get(key string) ([]byte, error)
	List(regexp ...string) ([]string, error)
	Delete(key string) error
	Purge(regexp ...string) error
}

//go:generate mockery --name persistentStorage --output ../mocks --filename persistent_storage_mock.go --structname PersistentStorageMock
type persistentStorage interface {
	Save(key string, value []byte, ttlDays ...int) (*persistentstorage.ObjectInfo, error)
	Get(key string, version ...string) (*persistentstorage.Object, error)
	List() ([]*persistentstorage.ObjectInfo, error)
	ListVersions(key string) ([]*persistentstorage.ObjectInfo, error)
	Delete(key string, version ...string) error
}

//go:generate mockery --name centralizedConfig --output ../mocks --filename centralized_config_mock.go --structname CentralizedConfigMock
type centralizedConfig interface {
	GetConfig(key string, scope ...msg.Scope) (string, error)
	SetConfig(key, value string, scope ...msg.Scope) error
	DeleteConfig(key string, scope msg.Scope) error
}

//nolint:godox // Task to be done.
// TODO add metrics interface.

//go:generate mockery --name measurements --output ../mocks --filename measurements_mock.go --structname MeasurementsMock
type measurements interface{}

//go:generate mockery --name modelRegistry --output ../mocks --filename model_registry_mock.go --structname ModelRegistryMock
type modelRegistry interface {
	RegisterModel(model []byte, name, version, description, modelFormat string) error
	GetModel(name string, version ...string) (*modelregistry.Model, error)
	ListModels() ([]*modelregistry.ModelInfo, error)
	ListModelVersions(name string) ([]*modelregistry.ModelInfo, error)
	DeleteModel(name string) error
}

type KaiSDK struct {
	// Metadata
	ctx context.Context

	// Needed deps
	nats           *nats.Conn
	jetstream      nats.JetStreamContext
	requestMessage *kai.KaiNatsMessage

	// Main methods
	Logger            logr.Logger
	Metadata          metadata
	Messaging         messaging
	Storage           Storage
	ModelRegistry     modelRegistry
	CentralizedConfig centralizedConfig
	Measurements      measurements
}

func NewKaiSDK(logger logr.Logger, natsCli *nats.Conn, jetstreamCli nats.JetStreamContext) KaiSDK {
	metadata := meta.New()

	centralizedConfigInst, err := centralizedconfiguration.New(logger, jetstreamCli)
	if err != nil {
		logger.WithName("[CENTRALIZED CONFIGURATION]").
			Error(err, "Error initializing Centralized Configuration")
		os.Exit(1)
	}

	ephemeralStg, err := objectstore.New(logger, jetstreamCli)
	if err != nil {
		logger.WithName("[EPHEMERAL STORAGE]").Error(err, "Error initializing ephemeral storage")
		os.Exit(1)
	}

	persistentStg, err := persistentstorage.New(logger)
	if err != nil {
		logger.WithName("[PERSISTENT STORAGE]").Error(err, "Error initializing persistent storage")
		os.Exit(1)
	}

	storageManager := Storage{
		Ephemeral:  ephemeralStg,
		Persistent: persistentStg,
	}

	messagingInst := msg.New(logger, natsCli, jetstreamCli, nil)

	modelRegistryInst, err := modelregistry.New(logger)
	if err != nil {
		logger.WithName("[MODEL REGISTRY]").Error(err, "Error initializing model registry")
		os.Exit(1)
	}

	sdk := KaiSDK{
		ctx:               context.Background(),
		nats:              natsCli,
		jetstream:         jetstreamCli,
		Logger:            logger,
		Metadata:          metadata,
		Messaging:         messagingInst,
		Storage:           storageManager,
		ModelRegistry:     modelRegistryInst,
		CentralizedConfig: centralizedConfigInst,
		Measurements:      nil,
	}

	return sdk
}

func (sdk *KaiSDK) GetRequestID() string {
	if sdk.requestMessage == nil {
		return ""
	}

	return sdk.requestMessage.GetRequestId()
}

func ShallowCopyWithRequest(sdk *KaiSDK, requestMsg *kai.KaiNatsMessage) KaiSDK {
	hSdk := *sdk
	hSdk.requestMessage = requestMsg
	hSdk.Logger = sdk.Logger.WithValues(LoggerRequestID, requestMsg.GetRequestId())
	hSdk.Messaging = msg.New(hSdk.Logger, sdk.nats, sdk.jetstream, requestMsg)

	return hSdk
}
