FROM golang:alpine AS build
WORKDIR /build/
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-w -s -v' -v -o ./dist/go-microservice ./cmd/microservice/*.go

FROM alpine
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=build /build/dist/go-microservice /app/
COPY --from=build /build/config.toml /app/
ENTRYPOINT ["./go-microservice"]
CMD ["serve"]
EXPOSE 3035