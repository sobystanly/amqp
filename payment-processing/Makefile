BINARY_NAME_PAYMENT=payment-processing

## build: Build binaries
build:
	@echo "Building..."
	@env CGO_ENABLED=0 go build -ldflags="-s -w" -o ${BINARY_NAME_PAYMENT}
	@echo "Built!"

## run: builds and runs the applications
run: build
	@echo "Starting..."
	@env ./${BINARY_NAME_PAYMENT}
	@echo "Started!"

## clean: runs go clean and deletes binaries
clean:
	@echo "Cleaning..."
	@go clean
	@rm ${BINARY_NAME_PAYMENT}
	@echo "Cleaned!"

## start: an alias to run
start: run

## stop: stops the running applications
stop:
	@echo "Stopping..."
	@-pkill -SIGTERM -f "./${BINARY_NAME_PAYMENT}"
	@echo "Stopped!"

## restart: stops and starts the applications
restart: stop start

## test: runs all tests
test:
	go test -v ./... # Run tests for the payment-processing service
