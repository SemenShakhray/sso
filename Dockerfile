FROM golang:1.23.4-alpine AS builder

WORKDIR /app

RUN apk --no-cache add bash git make gettext musl-dev

# dependencies
COPY go.mod go.sum ./
RUN go mod download

# build
COPY cmd ./cmd
COPY internal ./internal
COPY pkg ./pkg
RUN CGO_ENABLED=0 go build -o /app/sso ./cmd/sso/main.go

FROM alpine AS runner

WORKDIR /app

COPY --from=builder /app/sso /app/sso
COPY .config /app/.config
COPY /migrations /app/migrations

CMD ["/app/sso"]