FROM golang:1.24 AS builder

WORKDIR /src

COPY ../../go.mod ./
COPY ../../go.sum ./
RUN go mod download


COPY ../../ ./

RUN cd cmd/accounts_service && go build -v -o accounts_service

FROM golang:1.24

WORKDIR /app
COPY --from=builder /src/cmd/accounts_service ./
COPY --from=builder /src/signature.pem ./
COPY --from=builder /src/signature.pub ./

CMD ["/app/accounts_service"]
