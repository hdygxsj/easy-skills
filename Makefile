.PHONY: all build build-all test clean dev serve install

# Default target
all: build

# Build Go CLI with embedded web UI
build:
	@echo "Building web UI..."
	cd web && npm install && npm run build
	@echo "Copying web dist to CLI..."
	cp -r web/dist internal/cli/dist
	@echo "Building CLI..."
	go build -o easy-skills ./cmd/easy-skills
	@echo "Done! Binary: ./easy-skills"

# Build CLI only (requires web/dist to exist)
build-cli:
	go build -o easy-skills ./cmd/easy-skills

# Build web UI only
build-web:
	cd web && npm run build
	cp -r web/dist internal/cli/dist

# Clean build artifacts
clean:
	rm -rf easy-skills
	rm -rf web/dist
	rm -rf internal/cli/dist

# Run web dev server
dev:
	cd web && npm run dev

# Run web dev server and serve CLI
dev-all: build-web
	go run ./cmd/easy-skills serve --dev

# Start the GUI server
serve:
	./easy-skills serve

# Start the GUI server in dev mode
serve-dev:
	./easy-skills serve --dev

# Install to system
install:
	go install ./cmd/easy-skills
