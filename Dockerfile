FROM golang:latest AS build

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o rateLimiter ./cmd/main.go

FROM scratch
WORKDIR /app
COPY --from=build /app/rateLimiter .
COPY --from=build /app/.env .
ENTRYPOINT [ "./rateLimiter" ]
