FROM golang:1.24

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

EXPOSE 8080

CMD ["sh", "-c", "swag init -d ./cmd,./internal -g main.go --parseInternal --parseDependency --parseDepth 2 -o ./docs && go run ./cmd"]