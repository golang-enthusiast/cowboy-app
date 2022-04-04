package main

import (
	"context"
	"encoding/json"
	"fmt"

	pkgCowboy "cowboy-app/internal/cowboy"
	pkgDomain "cowboy-app/internal/domain"
	pkgDynamodb "cowboy-app/internal/dynamodb"
	errors "cowboy-app/internal/error"
	pkgHelpers "cowboy-app/internal/helpers"
	pkgQueue "cowboy-app/internal/sqsqueue"
	pkgWorker "cowboy-app/internal/worker"

	"github.com/aws/aws-sdk-go/service/sqs"
	kitzapadapter "github.com/go-kit/kit/log/zap"
	"github.com/go-kit/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {

	// Load application config.
	//
	config, err := pkgHelpers.LoadAppConfig()
	if err != nil {
		panic(err)
	}

	// Create a single logger, which we'll use and give to other components.
	//
	zapLogger, _ := zap.NewProduction()
	defer func() {
		_ = zapLogger.Sync()
	}()

	var logger log.Logger
	logger = kitzapadapter.NewZapSugarLogger(zapLogger, zapcore.InfoLevel)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)
	logger = log.With(logger, "name", config.CowboyName)

	// Setup AWS services.
	//
	awsSession := pkgHelpers.GetAwsSession()
	sqsClient := sqs.New(awsSession)

	// Setup repository.
	//
	cowboyRepository := pkgDynamodb.NewCowboyRepository(awsSession, config.CowboyTableName)

	// Setup service layer.
	//
	queueService := pkgQueue.NewQueueService(sqsClient, logger)
	cowboyRandService := pkgCowboy.NewCowboyRandomGenerator()
	cowboyService := pkgCowboy.NewService(config.CowboyName, cowboyRepository, queueService, cowboyRandService, logger)

	// Setup worker.
	//
	cowboyWorker := pkgWorker.NewCowboyWorker(
		&pkgWorker.Props{
			WorkerName:          config.CowboyName,
			QueueName:           config.CowboyName,
			MaxNumberOfMessages: 10,
		},
		sqsClient,
		cowboyService,
		logger,
	)

	cowboyWorker.Start(context.Background(), pkgDomain.QueueHandlerFunc(func(ctx context.Context, msg *sqs.Message) error {
		// Define message type.
		messageTypeValue, ok := msg.MessageAttributes[pkgDomain.MessageTypeAttributeKey]
		if !ok {
			return errors.NewErrNotFound("Message type not found in attributes")
		}
		if messageTypeValue == nil {
			return errors.NewErrNotFound("Message type required")
		}
		messageType := pkgDomain.MessageType(*messageTypeValue.StringValue)

		// Switch message types.
		switch messageType {
		case pkgDomain.PrepareGunsMessageType:
			payload := &pkgDomain.PrepareGunsMessage{}
			return cowboyWorker.HandlePrepareGunsMessage(ctx, payload)
		case pkgDomain.ShootMessageType:
			payload := &pkgDomain.ShootMessage{}
			err := json.Unmarshal([]byte(*msg.Body), payload)
			if err != nil {
				return err
			}
			return cowboyWorker.HandleShootMessage(ctx, payload)
		case pkgDomain.WinnerMessageType:
			payload := &pkgDomain.WinnerMessage{}
			err := json.Unmarshal([]byte(*msg.Body), payload)
			if err != nil {
				return err
			}
			return cowboyWorker.HandleWinnerMessage(ctx, payload)
		default:
			return errors.NewErrNotFound(fmt.Sprintf("Unsupported message type %v", messageType))
		}
	}))
}
