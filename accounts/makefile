BINARY=accounts
ENV= -e APP_DB_HOST="127.0.0.1"
DEV_CONTAINER_PORT= -p 7020:80
DOCKER_REGISTRY_HOST=sipsynergy
VERSION=test

build:
	docker run --rm -v ${PWD}:/go/src/app -w /go/src/app -e GOOS=linux -e GOARCH=386 sipsynergy/go-builder /bin/sh -c "godep get && godep go build" && \
    docker build -t ${DOCKER_REGISTRY_HOST}/${BINARY}:${VERSION} .

start:
	docker run -d ${ENV} ${DEV_CONTAINER_PORT} --name ${BINARY} ${DOCKER_REGISTRY_HOST}/${BINARY}:${VERSION}

stop:
	docker stop ${BINARY} && docker rm ${BINARY}

exec:
	docker exec -it ${BINARY} /bin/sh

clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

.PHONY: clean
