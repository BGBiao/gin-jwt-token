.PHONY: list vet default fmt fmt-check clean
all: list default fmt docker_build docker_tag docker_push clean
# BINARY即为app的名称
BINARY="ali-warn-api"
VERSION=0.0.1
BUILD=`date +%F`
IMAGE="harbor.soulapp-inc.cn/soul-ops/${BINARY}"

PACKAGES=`go list ./... | grep -v /vendor/`
VETPACKAGES=`go list ./... | grep -v /vendor/ | grep -v /examples/`
GOFILES=`find . -name "*.go" -type f -not -path "./vendor/*"`

default:
	@echo "build the ${BINARY}"
	@#这个--tags=jsoniter是个什么鬼啊
	@GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ${BINARY}  -tags=jsoniter
	@echo "build done."

list:
	@echo ${PACKAGES}
	@echo ${VETPACKAGES}
	@echo ${GOFILES}

fmt:
	@echo "fmt the project"
	@gofmt -s -w ${GOFILES}

fmt-check:
	@diff=$$(gofmt -s -d $(GOFILES)); \
    if [ -n "$$diff" ]; then \
    echo "Please run 'make fmt' and commit the result:"; \
    echo "$${diff}"; \
    exit 1; \
    fi;

install:
	@govendor sync -v

test:
	@go test -cpu=1,2,4 -v -tags integration ./...

vet:
	@echo "check the project codes."
	@go vet $(VETPACKAGES)
	@echo "check done."

docker_build:
	@docker build --build-arg=app=${BUILD} -t ${IMAGE}:${BUILD} .

docker_tag:
	@docker tag  ${IMAGE}:${BUILD} ${IMAGE}

docker_push:
	@docker push ${IMAGE}:${BUILD}

clean:
	@if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
