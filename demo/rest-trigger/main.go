package main

import (
	context2 "context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/protobuf/proto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/konstellation-io/kai-sdk/go-sdk/runner"
	"github.com/konstellation-io/kai-sdk/go-sdk/runner/trigger"
	"github.com/konstellation-io/kai-sdk/go-sdk/sdk"
)

func main() {
	runner.
		NewRunner().
		TriggerRunner().
		WithInitializer(initializer).
		WithRunner(restServerRunner).
		WithFinalizer(func(kaiSDK sdk.KaiSDK) {
			kaiSDK.Logger.Info("Finalizer")
		}).
		Run()
}

func initializer(kaiSDK sdk.KaiSDK) {
	kaiSDK.Logger.Info("Writing test value to the ephemeral store", "value", "testValue")
	err := kaiSDK.Storage.Ephemeral.Save("test", []byte("testValue"))
	if err != nil {
		kaiSDK.Logger.Error(err, "Error saving object")
	}

	kaiSDK.Logger.Info("Writing test value to the centralized config",
		"value", "testConfigValue")
	err = kaiSDK.CentralizedConfig.SetConfig("test", "testConfigValue")
	if err != nil {
		kaiSDK.Logger.Error(err, "Error setting config")
	}

	kaiSDK.Logger.V(1).Info("Metadata",
		"process", kaiSDK.Metadata.GetProcess(),
		"product", kaiSDK.Metadata.GetProduct(),
		"workflow", kaiSDK.Metadata.GetWorkflow(),
		"version", kaiSDK.Metadata.GetVersion(),
		"kv_product", kaiSDK.Metadata.GetProductCentralizedConfigurationName(),
		"kv_workflow", kaiSDK.Metadata.GetWorkflowCentralizedConfigurationName(),
		"kv_process", kaiSDK.Metadata.GetProcessCentralizedConfigurationName(),
		"ephemeral_store", kaiSDK.Metadata.GetEphemeralStorageName(),
	)

	kaiSDK.Logger.V(1).Info("PathUtils",
		"getBasePath", kaiSDK.PathUtils.GetBasePath(),
		"composeBasePath", kaiSDK.PathUtils.ComposePath("test"))
}

func restServerRunner(tr *trigger.Runner, kaiSDK sdk.KaiSDK) {
	kaiSDK.Logger.Info("Starting http server", "port", 8080)

	bgCtx, stop := signal.NotifyContext(context2.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	r := gin.Default()
	r.GET("/hello", responseHandler(kaiSDK, tr.GetResponseChannel))
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			kaiSDK.Logger.Error(err, "Error running http server")
		}
	}()

	<-bgCtx.Done()
	stop()
	kaiSDK.Logger.Info("Shutting down server...")

	// The sdk is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	bgCtx, cancel := context2.WithTimeout(context2.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(bgCtx); err != nil {
		kaiSDK.Logger.Error(err, "Error shutting down server")
	}

	kaiSDK.Logger.Info("Server stopped")
}

func responseHandler(kaiSDK sdk.KaiSDK, getResponseChannel func(requestID string) <-chan *anypb.Any) func(c *gin.Context) {
	return func(c *gin.Context) {
		nameRequest := c.Query("name")

		stringb := fmt.Sprintf("Hello %s!", nameRequest)

		anyValue := &wrappers.StringValue{
			Value: stringb,
		}

		reqID := uuid.New().String()

		kaiSDK.Logger.Info("Sending message to nats",
			"execution ID", reqID, "message", anyValue)

		responseChannel := getResponseChannel(reqID)
		err := kaiSDK.Messaging.SendOutputWithRequestID(anyValue, reqID)
		if err != nil {
			kaiSDK.Logger.Error(err, "Error sending message to nats")
			return
		}

		response := <-responseChannel

		stringValue := &wrappers.StringValue{}

		// Unmarshall response to StringValue
		err = proto.Unmarshal(response.GetValue(), stringValue)

		kaiSDK.Logger.Info("Response received from nats", "response", stringValue.GetValue())

		c.JSON(http.StatusOK, gin.H{
			"message": strings.Split(stringValue.GetValue(), ","),
		})
	}
}
