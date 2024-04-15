FROM golang:latest AS build

WORKDIR /app
COPY . .

RUN go mod download
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o rateLimiter ./cmd/main.go

# ----------------------------

FROM alpine:latest 
WORKDIR /app
COPY --from=build /app/rateLimiter .
COPY --from=build /app/.env .

ENTRYPOINT [ "./rateLimiter" ]
