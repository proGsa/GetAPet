FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN chmod +x tests/run-integration-tests.sh tests/run-unit-tests.sh

# CMD ["tests/run-integration-tests.sh"]
