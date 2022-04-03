unittest:
	# Runs all unit-tests.
	AWS_REGION=us-east-1 bash -c 'go test $$(go list ./... | grep -v '/cmd') -v'

testmocks:
	# Generate mocks stubs.
	mockgen \
		-destination=internal/mocks/mock_cowboy_repository.go \
		-package mocks cowboy-app/internal/domain CowboyRepository

	mockgen \
		-destination=internal/mocks/mock_queue_service.go \
		-package mocks cowboy-app/internal/domain QueueService	

	mockgen \
		-destination=internal/mocks/mock_cowboy_random_generator.go \
		-package mocks cowboy-app/internal/domain CowboyRandomGenerator
