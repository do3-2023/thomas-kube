## Build the Go API
##
FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY go.mod go.sum .env ./ 

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o api


## Dev image
##
FROM  golang:1.20-alpine AS dev

WORKDIR /app

COPY . .

RUN go mod download && go mod verify


CMD ["go", "run", "main.go"]


## Production image
##
FROM alpine:latest AS prod

RUN addgroup -S myuser && adduser -S myuser -G myuser

USER myuser

WORKDIR /

COPY --from=builder /app/api .
COPY --from=builder /app/migrations ./migrations

CMD [ "/api" ]  