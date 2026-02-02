BIN_GO_SERVICE=pppk-json

MIGRATE_CMD=docker compose -f compose.dev.yaml run --rm --entrypoint /bin/sh migrate

ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))

$(eval $(ARGS):;@:)

.PHONY: dev-go-server
dev-go-server:
	@echo "Running Application with args: $(ARGS)"

	mkdir -p backend/web/dist
	touch backend/web/dist/index.html

	cd backend && go mod tidy && gow run ./cmd/server $(ARGS)

prod-go-server: prod-frontend
	@echo "Running Application with args: $(ARGS)"
	cd backend \
	&& go mod tidy \
	&& rm -f ./bin/$(BIN_GO_SERVICE) || true \
	&& go build -o ./bin/$(BIN_GO_SERVICE) ./cmd/server \
	&& ./bin/$(BIN_GO_SERVICE) $(ARGS)

prod-frontend:
	@echo "Building Application for production with args: $(ARGS)"
	cd frontend && npm install && npm run build $(ARGS)

	@echo "Cleaning old assets in backend..."
	rm -rf backend/web/dist

	@echo "Copying new assets to backend/web..."
	mkdir -p backend/web
	cp -r frontend/dist backend/web/

dev-frontend:
	@echo "Running Application with args: $(ARGS)"
	 cd frontend && npm run dev $(ARGS)

db-seed:
	@echo "Running Application with args: $(ARGS)"
	cd backend && go mod tidy && go run ./cmd/seeder $(ARGS)

migrate-create:
	@echo "Creating migration files..."
	# Kita panggil binary 'migrate' secara manual di sini karena entrypoint sudah diganti sh
	$(MIGRATE_CMD) -c 'migrate create -ext sql -dir /migrations -seq $(name)'

migrate-up:
	@echo "Running migrations..."
	# Gunakan -c agar env var $$DATABASE_URL terbaca oleh shell
	$(MIGRATE_CMD) -c 'migrate -path=/migrations -database $$DATABASE_URL up'

migrate-down:
	@echo "Rolling back migrations..."
	$(MIGRATE_CMD) -c 'migrate -path=/migrations -database $$DATABASE_URL down 1'

migrate-force:
	$(MIGRATE_CMD) -c 'migrate -path=/migrations -database $$DATABASE_URL force $(version)'

migrate-reset:
	@echo "⚠️ DESTROYING ALL DATA..."
	$(MIGRATE_CMD) -c 'migrate -path=/migrations -database $$DATABASE_URL down'

swag-gen:
	@echo "Generating swagger documentation..."
	cd backend && swag init -g cmd/server/main.go -o ./docs

compose-up:
	@echo "Starting containers..."
	docker compose -f compose.dev.yaml up -d

compose-down:
	@echo "Stopping containers..."
	docker compose -f compose.dev.yaml down

