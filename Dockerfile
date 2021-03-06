# build stage
FROM golang:alpine AS build-env

WORKDIR /go/src/github.com/vdv/nginx-clickhouse

RUN apk update && apk add make g++ git curl

COPY . /go/src/github.com/vdv/nginx-clickhouse
RUN cd /go/src/github.com/vdv/nginx-clickhouse && go get -v .
RUN cd /go/src/github.com/vdv/nginx-clickhouse && make build

# final stage
FROM scratch

COPY --from=build-env /go/src/github.com/vdv/nginx-clickhouse/nginx-clickhouse /
CMD [ "/nginx-clickhouse" ]
