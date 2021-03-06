FROM golang:1.17 AS build

# install protobuf from source
RUN apt-get update && \
    apt-get -y install git unzip build-essential autoconf libtool
RUN git clone https://github.com/google/protobuf.git
RUN cd protobuf && ./autogen.sh
RUN cd protobuf && ./configure
RUN cd protobuf && make
RUN cd protobuf && make install
RUN cd protobuf && ldconfig
RUN cd protobuf && make clean
RUN rm -r protobuf

WORKDIR /go/src/github.com/auth_service/

COPY ./cmd ./cmd
COPY ./api ./api
COPY ./pkg ./pkg
COPY ./service ./service
COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum
COPY ./models ./models
COPY ./internal ./internal
COPY ./middlewhare ./middlewhare
COPY ./tools ./tools
RUN go test ./...
RUN mkdir bin
RUN go build  -o bin/auth_bin cmd/main.go

FROM postgres
COPY migrations.sql /docker-entrypoint-initdb.d/
COPY --from=build /go/src/github.com/auth_service/bin/auth_bin .
EXPOSE 8080