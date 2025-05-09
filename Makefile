BIN_DIR      := bin
GO           := go

.PHONY: client gateway producer calculate consumer payment db insert testReq auth test up down restart logs proto

gateway:
	@echo "🚀 Запуск gateway..."
	@$(GO) build -o $(BIN_DIR)/gateway ./gateway/cmd
	@$(BIN_DIR)/gateway

client:
	@echo "🚀 Запуск client..."
	@$(GO) build -o $(BIN_DIR)/client ./client/presentation/cli
	@$(BIN_DIR)/client $(ARGS)

producer:
	@echo "🚀 Запуск producer..."
	@$(GO) build -o $(BIN_DIR)/producer ./producer/cmd/
	@$(BIN_DIR)/producer

consumer:
	@echo "🚀 Запуск consumer..."
	@$(GO) build -o $(BIN_DIR)/consumer ./consumer/cmd/
	@$(BIN_DIR)/consumer

payment:
	@echo "🚀 Запуск payment..."
	@$(GO) build -o $(BIN_DIR)/payment ./payment/cmd/
	@$(BIN_DIR)/payment

calculate:
	@echo "🚀 Запуск calculator..."
	@$(GO) build -o $(BIN_DIR)/calculate ./calculator/cmd/
	@$(BIN_DIR)/calculate

db:
	@echo "🚀 Запуск database..."
	@$(GO) build -o $(BIN_DIR)/db ./database/cmd/
	@$(BIN_DIR)/db

auth:
	@echo "🚀 Запуск auth..."
	@$(GO) build -o $(BIN_DIR)/auth ./auth/
	@$(BIN_DIR)/auth


test:
	@echo "🧪 Запуск тестов..."
	@$(GO) test -v -race ./...

insert:
	@echo "🧪 inserting..."
	@chmod +x ./scripts/insertPackages.sh
	@./scripts/insertPackages.sh

reqcalculate:
	@echo "🚀 Calculating..."
	@go run ./scripts/calculator/calculator_calcs.go

testReq:
	@chmod +x ./scripts/test_requests.sh
	@./scripts/test_requests.sh

up:
	docker-compose up -d

down:
	docker-compose down

restart:
	docker-compose down && docker-compose up -d

logs:
	docker-compose logs -f

protocalc:
	export PATH="$PATH:$(go env GOPATH)/bin"
	/opt/homebrew/bin/protoc --proto_path=proto \
		--go_out=paths=source_relative:proto \
		--go-grpc_out=paths=source_relative:proto \
		proto/calculator/calculator.proto

protoauth:
	export PATH="$PATH:$(go env GOPATH)/bin"
	/opt/homebrew/bin/protoc --proto_path=proto \
		--go_out=paths=source_relative:proto \
		--go-grpc_out=paths=source_relative:proto \
		proto/auth/auth.proto

protopay:
	export PATH="$PATH:$(go env GOPATH)/bin"
	/opt/homebrew/bin/protoc --proto_path=proto \
		--go_out=paths=source_relative:proto \
		--go-grpc_out=paths=source_relative:proto \
		proto/payment/payment.proto