######
# Jellyfin-Newsletter Go-engine entrypoint builder
######

FROM golang:1.26-alpine AS entrypoint-builder
WORKDIR /app

COPY engine-go/entrypoint/main.go .

RUN go mod init entrypoint && \
    CGO_ENABLED=0 GOOS=linux go build -o entrypoint main.go

######
# Jellyfin-Newsletter Go-engine application builder
######

FROM golang:1.26-alpine AS app-builder
WORKDIR /app
ARG VERSION="dev"


COPY engine-go/go.mod engine-go/go.sum ./
RUN go mod download

COPY engine-go/internal/ /app/internal/
COPY engine-go/main.go /app/main.go

RUN mkdir /app/config
COPY config/config-example.yml /app/config/

RUN mkdir /app/previews

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.version=${VERSION}"  -o /app/jellyfin-newsletter .

FROM gcr.io/distroless/static AS runtime
COPY --from=app-builder /app/jellyfin-newsletter /app/jellyfin-newsletter
COPY --from=app-builder /app/config /app/config
COPY --from=entrypoint-builder /app/entrypoint /app/entrypoint

ENTRYPOINT ["/app/entrypoint", "-config", "/app/config/config.yml"]
