FROM golang:alpine3.18
WORKDIR /app
COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum
RUN  go mod download
COPY ./main.go ./main.go
RUN go build main.go
CMD [ "./main" ]
