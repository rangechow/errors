PKGS := github.com/rangechow/errors
SRCDIRS := $(shell go list -f '{{.Dir}}' $(PKGS))
GO := go

check: test vet gofmt

test: 
	$(GO) test $(PKGS)

vet: | test
	$(GO) vet $(PKGS)

gofmt:  
	@echo Checking code is gofmted
	@test -z "$(shell gofmt -s -l -d -e $(SRCDIRS) | tee /dev/stderr)"