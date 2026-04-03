FROM golang:1.26-alpine AS builder
WORKDIR /app
ARG VERSION="dev"


COPY engine-go/go.mod engine-go/go.sum ./
RUN go mod download

COPY engine-go/internal/ /app/internal/
COPY engine-go/main.go /app/main.go

RUN mkdir /app/config
COPY config/config-example.yml /app/config/
RUN chown 1001:1001 /app/config

RUN mkdir /app/previews
RUN chown 1001:1001 /app/previews

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.version=${VERSION}"  -o /app/jellyfin-newsletter .


FROM gcr.io/distroless/static AS runtime
COPY --from=builder /app/jellyfin-newsletter /app//jellyfin-newsletter
COPY --from=builder /app/config /app/config
USER 1001:1001
ENTRYPOINT ["/app/jellyfin-newsletter", "-config", "/app/config/config.yml"]
