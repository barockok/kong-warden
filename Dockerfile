FROM golang:alpine3.18
RUN apk add make
WORKDIR /app
RUN mkdir out
COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum
RUN  go mod download
COPY ./ ./
RUN make