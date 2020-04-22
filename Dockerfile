FROM golang:1.14 AS build
ENV CGO_ENABLED=0
WORKDIR /go/src
COPY . ./

RUN go test -v -cover ./...
RUN go test -c -tags integration -o ./bin/server.test ./server

RUN go build -a -o ./bin ./...

FROM alpine AS runtime
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates
COPY --from=build /go/src/bin/server* ./
ENTRYPOINT ["./server"]
