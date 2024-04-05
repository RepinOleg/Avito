BUILD=${CURDIR}/build
PACKAGE=github.com/RepinOleg/Banner_service

.PHONY: build run

build:
	go build -o ${BUILD}/server ${PACKAGE}/cmd/

run: db-up build
	${BUILD}/server

stop: db-down

db-up:
	docker-compose up -d
	@sleep 2

db-down:
	docker-compose down
