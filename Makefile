build:
	docker-compose build banner-app

run:
	docker-compose up banner-app

export TEST_DB_PASSWORD=test
export TEST_CONTAINER_NAME=test-db
export TEST_DB_USER=test

test.integration:
	docker run --rm -d -p 5432:5432 --name $$TEST_CONTAINER_NAME -e POSTGRES_USER=$$TEST_DB_USER -e POSTGRES_PASSWORD=$$TEST_DB_PASSWORD -e POSTGRES_DB=postgres postgres:latest
	docker cp assets/postgres/init.sql test-db:/docker-entrypoint-initdb.d/1-schema.sql
	@sleep 1
	go test -v ./tests/
	docker stop $$TEST_CONTAINER_NAME

lint:
	golangci-lint run

