# Clients docker
FROM golang:1.20 as build

ENV BIN_FILE /opt/client/client-app
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} ./cmd/client

FROM alpine:latest

ENV BIN_FILE "/opt/client/client-app"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

CMD ${BIN_FILE}