package dynamodb

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"cowboy-app/internal/domain"

	aws "github.com/aws/aws-sdk-go/aws"
	awsSession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	testSession          *awsSession.Session
	testCowboyRepository domain.CowboyRepository
	testCowboyTableName  = "test.cowboys"
)

func TestMain(m *testing.M) {

	// Setup Database
	ctx := context.Background()
	exposedPort := "4566"
	req := testcontainers.ContainerRequest{
		Image:        "localstack/localstack:latest",
		ExposedPorts: []string{exposedPort},
		WaitingFor:   wait.ForListeningPort(nat.Port(exposedPort)),
		Env: map[string]string{
			"DEBUG":    "1",
			"SERVICES": "dynamodb",
		},
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("Failed to start container %v", err)
	}

	defer func() { _ = container.Terminate(ctx) }()

	// Setup test dependencies.
	ip, err := container.Host(ctx)
	if err != nil {
		log.Fatalf("Failed to start container host %v", err)
	}
	port, err := container.MappedPort(ctx, nat.Port(exposedPort))
	if err != nil {
		log.Fatalf("Failed to start container port %v", err)
	}
	mockServerAddress := fmt.Sprintf("http://%s:%s", ip, port.Port())
	fmt.Printf("Mock server address: %s\n", mockServerAddress)

	// Init AWS Session.
	testSession, err = awsSession.NewSession(&aws.Config{
		Endpoint: aws.String(mockServerAddress),
	})
	if err != nil {
		log.Fatalf("Failed to start elasticsearch client %v", err)
	}

	testCowboyRepository = NewCowboyRepository(testSession, testCowboyTableName)

	exitVal := m.Run()
	os.Exit(exitVal)
}

func runSetup(t *testing.T, setup func() error) {
	err := setup()
	if err != nil {
		t.Errorf("setup faild with error: %s", err)
		t.FailNow()
	}
}
