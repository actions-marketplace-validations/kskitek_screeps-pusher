FROM golang:alpine AS builder

ARG LDFLAGS

WORKDIR /src

COPY ./ /src

RUN go mod tidy

RUN CGO_ENABLED=0 go build -ldflags="-w -s"

FROM gcr.io/distroless/static-debian11

COPY --from=builder /src/pusher /bin/pusher

ENTRYPOINT ["pusher"]
