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
		WithFinalizer(func(sdk sdk.KaiSDK) {
			sdk.Logger.Info("Finalizer")
		}).
		Run()
}

func initializer(sdk sdk.KaiSDK) {
	sdk.Logger.Info("Writing test value to the object store", "value", "testValue")
	err := sdk.Storage.Ephemeral.Save("test", []byte("testValue"))
	if err != nil {
		sdk.Logger.Error(err, "Error saving object")
	}

	sdk.Logger.Info("Writing test value to the centralized config",
		"value", "testConfigValue")
	err = sdk.CentralizedConfig.SetConfig("test", "testConfigValue")
	if err != nil {
		sdk.Logger.Error(err, "Error setting config")
	}

	sdk.Logger.V(1).Info("Metadata",
		"process", sdk.Metadata.GetProcess(),
		"product", sdk.Metadata.GetProduct(),
		"workflow", sdk.Metadata.GetWorkflow(),
		"version", sdk.Metadata.GetVersion(),
		"kv_product", sdk.Metadata.GetKeyValueStoreProductName(),
		"kv_workflow", sdk.Metadata.GetKeyValueStoreWorkflowName(),
		"kv_process", sdk.Metadata.GetKeyValueStoreProcessName(),
		"object_store", sdk.Metadata.GetEphemeralStorageName(),
	)

	sdk.Logger.V(1).Info("PathUtils",
		"getBasePath", sdk.PathUtils.GetBasePath(),
		"composeBasePath", sdk.PathUtils.ComposePath("test"))
}

func restServerRunner(tr *trigger.Runner, sdk sdk.KaiSDK) {
	sdk.Logger.Info("Starting http server", "port", 8080)

	bgCtx, stop := signal.NotifyContext(context2.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	r := gin.Default()
	r.GET("/hello", responseHandler(sdk, tr.GetResponseChannel))
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			sdk.Logger.Error(err, "Error running http server")
		}
	}()

	<-bgCtx.Done()
	stop()
	sdk.Logger.Info("Shutting down server...")

	// The sdk is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	bgCtx, cancel := context2.WithTimeout(context2.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(bgCtx); err != nil {
		sdk.Logger.Error(err, "Error shutting down server")
	}

	sdk.Logger.Info("Server stopped")
}

func responseHandler(sdk sdk.KaiSDK, getResponseChannel func(requestID string) <-chan *anypb.Any) func(c *gin.Context) {
	return func(c *gin.Context) {
		nameRequest := c.Query("name")

		stringb := fmt.Sprintf("Hello %s!", nameRequest)

		anyValue := &wrappers.StringValue{
			Value: stringb,
		}

		reqID := uuid.New().String()

		sdk.Logger.Info("Sending message to nats",
			"execution ID", reqID, "message", anyValue)

		responseChannel := getResponseChannel(reqID)
		err := sdk.Messaging.SendOutputWithRequestID(anyValue, reqID)
		if err != nil {
			sdk.Logger.Error(err, "Error sending message to nats")
			return
		}

		response := <-responseChannel

		stringValue := &wrappers.StringValue{}

		// Unmarshall response to StringValue
		err = proto.Unmarshal(response.GetValue(), stringValue)

		sdk.Logger.Info("Response received from nats", "response", stringValue.GetValue())

		c.JSON(http.StatusOK, gin.H{
			"message": strings.Split(stringValue.GetValue(), ","),
		})
	}
}
