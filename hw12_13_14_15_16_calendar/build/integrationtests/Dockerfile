FROM golang:1.22

RUN mkdir -p /opt/integration_tests
WORKDIR /opt/integration_tests

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY ./integrationtests .

ENV CONFIG_FILE /etc/integrationtests/config.yml
COPY ./configs/integrationtests-config.yml ${CONFIG_FILE}

CMD ["go", "test"]
