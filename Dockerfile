FROM golang as builder

ADD . .

RUN go build .

FROM alpine



