FROM golang:1.18-alpine AS builder

ARG LDFLAGS

WORKDIR /src

COPY ./ /src

RUN go mod tidy

RUN CGO_ENABLED=0 go build -ldflags="-w -s"

FROM gcr.io/distroless/static-debian11

COPY --from=builder /src/screeps-pusher /bin/screeps-pusher

ENTRYPOINT ["screeps-pusher"]
