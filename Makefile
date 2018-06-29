SRCS = main.go
VERSION_FLAGS=-ldflags "-X main.version=`cat VERSION` -X main.date=`date -u +%Y/%m/%d-%H:%M:%S`"
ARGS ?= -h ; printf "\n***\n* Set ARGS variable in make invocation\n***"
SCALE ?= 0.5

all: build

build:
	go build $(VERSION_FLAGS) $(SRCS)

clean:
	go clean

run:
	go run $(VERSION_FLAGS) $(SRCS) $(ARGS)

run_flat:
	go run $(VERSION_FLAGS) $(SRCS) -w world-50m.txt -s $(SCALE) -p flat | less -S

run_miller37:
	go run $(VERSION_FLAGS) $(SRCS) -w world-50m.txt -s $(SCALE) -p miller37 | less -S

run_miller43:
	go run $(VERSION_FLAGS) $(SRCS) -w world-50m.txt -s $(SCALE) -p miller43 | less -S

run_miller50:
	go run $(VERSION_FLAGS) $(SRCS) -w world-50m.txt -s $(SCALE) -p miller50 | less -S

update:
	dep ensure -v -update

format:
	find . -path ./vendor -prune -o -name '*.go' -exec gofmt -s -w {} \;
