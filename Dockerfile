# build stage
FROM golang:alpine AS build-env

WORKDIR /go/src/nginx-clickhouse

RUN apk update && apk add make g++ git curl

ADD . /go/src/nginx-clickhouse
RUN cd /go/src/nginx-clickhouse && go get -v .
RUN cd /go/src/nginx-clickhouse && make build

# final stage
FROM scratch

COPY --from=build-env /go/src/nginx-clickhouse/nginx-clickhouse /
CMD [ "/nginx-clickhouse" ]
