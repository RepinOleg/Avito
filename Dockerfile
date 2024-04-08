FROM golang:1.22

RUN go version
ENV GOPATH=/

COPY ./ ./

# install psql
RUN apt-get update
RUN apt-get -y install postgresql-client

RUN chmod +x wait-for-postgres.sh

# build go app
RUN go mod download
RUN go build -o banner-app ./cmd/main.go

CMD ["./banner-app"]