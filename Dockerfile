FROM golang:alpine

WORKDIR /gopher-pay

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o ./bin/api ./cmd/api && go build -o ./bin/migrate ./cmd/migrate && chmod +x ./bin/api ./bin/migrate

CMD ["./bin/api"]
EXPOSE 8080