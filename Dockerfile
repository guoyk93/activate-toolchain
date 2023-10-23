FROM golang:1.21 AS builder
ENV CGO_ENABLED 0
WORKDIR /go/src/app
ADD . .
RUN go build -o ./cmd/activate-toolchain /activate-toolchain

FROM debian:12
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
COPY --from=builder /activate-toolchain /activate-toolchain
CMD ["/activate-toolchain"]
