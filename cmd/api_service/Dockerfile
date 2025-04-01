FROM golang:1.24 AS builder

WORKDIR /src

COPY ../../go.mod ./
COPY ../../go.sum ./
RUN go mod download


COPY ../../ ./

RUN cd cmd/api_service && go build -v -o api_service

FROM golang:1.24

WORKDIR /app
COPY --from=builder /src/cmd/api_service ./

CMD ["/app/api_service"]
