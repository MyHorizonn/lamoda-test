FROM golang:alpine
RUN apk update && apk add --no-cache git && apk add --no-cach bash && apk add build-base
RUN mkdir /lamoda-test
WORKDIR /lamoda-test
COPY . .
RUN go get -d -v ./...
RUN go install -v ./...
RUN go build ./cmd/goods
EXPOSE 8000