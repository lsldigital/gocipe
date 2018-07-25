VERSION := 1.0.0
COMMIT := `git rev-parse HEAD`
DATE := `date +%FT%T%z`

.PHONY: build generate install

build:
	@echo ":: Build"
	@echo "   - Cleaning"
	@go clean
	@echo "   - Building"
	@go build -ldflags "-X=main.appVersion=$(VERSION) -X=main.appCommit=$(COMMIT) -X=main.appBuilt=$(DATE)"
	@echo -e "   ✓ Done\n"

generate:
	@echo ":: Generate"
	@echo "   - Running go:generate"
	@go generate
	@echo -e "   ✓ Done\n"

install:
	@echo ":: Install"
	@echo "   - Cleaning"
	@go clean
	@echo "   - Building and installing"
	@go install -ldflags "-X=main.appVersion=$(VERSION) -X=main.appCommit=$(COMMIT) -X=main.appBuilt=$(DATE)"
	@echo -e "   ✓ Done\n"
