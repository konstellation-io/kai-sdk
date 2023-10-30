package sdk

import (
	"context"
	"os"

	kaiCommon "github.com/konstellation-io/kai-sdk/go-sdk/internal/common"

	persistentstorage "github.com/konstellation-io/kai-sdk/go-sdk/sdk/persistent-storage"

	centralizedConfiguration "github.com/konstellation-io/kai-sdk/go-sdk/sdk/centralized-configuration"
	pathutils "github.com/konstellation-io/kai-sdk/go-sdk/sdk/path-utils"

	objectstore "github.com/konstellation-io/kai-sdk/go-sdk/sdk/ephemeral-storage"

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
	GetErrorMessage() string

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
	Save(key string, value []byte) (string, error)
	Get(key string, version ...string) ([]byte, error)
	List() ([]string, error)
	ListVersions(key string) ([]string, error)
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

type KaiSDK struct {
	// Metadata
	ctx context.Context

	// Needed deps
	nats           *nats.Conn
	jetstream      nats.JetStreamContext
	requestMessage *kai.KaiNatsMessage

	// Main methods
	Logger            logr.Logger
	PathUtils         pathUtils
	Metadata          metadata
	Messaging         messaging
	CentralizedConfig centralizedConfig
	Measurements      measurements
	Storage           Storage
}

func NewKaiSDK(logger logr.Logger, natsCli *nats.Conn, jetstreamCli nats.JetStreamContext) KaiSDK {
	metadata := meta.NewMetadata()

	logger = logger.WithValues(
		kaiCommon.LoggerProductID, metadata.GetProduct(),
		kaiCommon.LoggerVersionID, metadata.GetVersion(),
		kaiCommon.LoggerWorkflowID, metadata.GetWorkflow(),
		kaiCommon.LoggerProcessID, metadata.GetProcess(),
	)

	centralizedConfigInst, err := centralizedConfiguration.NewCentralizedConfiguration(logger, jetstreamCli)
	if err != nil {
		logger.WithName("[CENTRALIZED CONFIGURATION]").
			Error(err, "Error initializing Centralized Configuration")
		os.Exit(1)
	}

	ephemeralStorage, err := objectstore.NewEphemeralStorage(logger, jetstreamCli)
	if err != nil {
		logger.WithName("[EPHEMERAL STORAGE]").Error(err, "Error initializing ephemeral storage")
		os.Exit(1)
	}

	persistentStorage, err := persistentstorage.NewPersistentStorage(logger.WithName("[PERSISTENT STORAGE]"))
	if err != nil {
		logger.WithName("[PERSISTENT STORAGE]").Error(err, "Error initializing persistent storage")
		os.Exit(1)
	}

	storageManager := Storage{
		Ephemeral:  ephemeralStorage,
		Persistent: persistentStorage,
	}

	messagingInst := msg.NewMessaging(logger.WithName("[MESSAGING]"), natsCli, jetstreamCli, nil)

	sdk := KaiSDK{
		ctx:               context.Background(),
		nats:              natsCli,
		jetstream:         jetstreamCli,
		Logger:            logger,
		PathUtils:         pathutils.NewPathUtils(logger),
		Metadata:          metadata,
		Messaging:         messagingInst,
		CentralizedConfig: centralizedConfigInst,
		Measurements:      nil,
		Storage:           storageManager,
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
