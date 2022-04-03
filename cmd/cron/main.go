package main

import (
	"context"
	"cowboy-app/internal/domain"
	pkgDynamodb "cowboy-app/internal/dynamodb"
	pkgHelpers "cowboy-app/internal/helpers"
	pkgQueue "cowboy-app/internal/sqsqueue"

	"github.com/aws/aws-sdk-go/service/sqs"
	kitzapadapter "github.com/go-kit/kit/log/zap"
	"github.com/go-kit/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {

	// Load application config.
	//
	config, err := pkgHelpers.LoadCronConfig()
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

	// Notify cowboys.
	//
	cowboys, err := cowboyRepository.List(&domain.CowboySearchCriteria{
		HealthGt: 0,
		Limit:    100,
	})
	if err != nil {
		_ = logger.Log("Failed to execute list of cowboys", err.Error())
	}
	for _, cowboy := range cowboys {
		_, err = queueService.SendMessage(context.Background(), cowboy.Name, &domain.PrepareGunsMessage{})
		if err != nil {
			_ = logger.Log("Failed to send message", cowboy.Name, err.Error())
		}
	}
}
