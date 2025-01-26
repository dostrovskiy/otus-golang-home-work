# build/storer/Dockerfile
FROM golang:1.22 as build

ENV BIN_FILE /opt/storer/storer-app
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

# Кэшируем слои с модулями
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

# Собираем статический бинарник Go
ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} cmd/storer/*

# На выходе тонкий образ
FROM alpine:3.9

LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="storer"
LABEL MAINTAINERS="dostrovskiy@diasoft.ru"

ENV BIN_FILE "/opt/storer/storer-app"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ENV CONFIG_FILE /etc/storer/config.yml
COPY ./configs/storer-config.yml ${CONFIG_FILE}

CMD ${BIN_FILE}
