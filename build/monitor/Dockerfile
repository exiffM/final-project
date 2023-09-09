# Monitor docker
FROM golang:1.20 as build

ENV BIN_FILE /opt/system.monitor/monitor-app
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

# Кэшируем слои с модулями
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

# Собираем статический бинарник Go (без зависимостей на Си API),
# иначе он не будет работать в alpine образе.
ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} cmd/monitor/*

# На выходе тонкий образ
FROM ubuntu:latest

EXPOSE 50051
RUN apt-get update
RUN apt-get install sysstat -y
RUN apt-get install net-tools
ENV BIN_FILE "/opt/system.monitor/monitor-app"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ENV CONFIG_FILE /etc/system.monitor/config.yml
COPY ./configs/config.yml ${CONFIG_FILE}

CMD ${BIN_FILE} -config ${CONFIG_FILE}