VERSION := 1.0.0
COMMIT := `git rev-parse HEAD`
DATE := `date +%FT%T%z`

build:
	go clean && go build -ldflags "-X=main.appVersion=$(VERSION) -X=main.appCommit=$(COMMIT) -X=main.appBuilt=$(DATE)"

install:
	go clean && go install -ldflags "-X=main.appVersion=$(VERSION) -X=main.appCommit=$(COMMIT) -X=main.appBuilt=$(DATE)"
