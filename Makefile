GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
BINARY_NAME=timecop

TIMECOP_ROOT=cmd/timecop/*.go

all:	clean test build
build:	timecop

timecop:
		$(GOBUILD) -o $(BINARY_NAME) $(TIMECOP_ROOT)

clean:
	if [ -f ${BINARY_NAME} ]; then rm ${BINARY_NAME}; fi
test:
	$(GOTEST) -v ./...

.PHONY: build timecop clean test
