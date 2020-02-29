FROM golang:1.13-alpine AS build

# Consistency across Lambda containers
ENV LANG en_US.UTF-8

RUN addgroup -g 666 -S lambda_build && \
    adduser -S -s /sbin/nologin -G lambda_build -u 666 -D lambda_build
USER lambda_build

COPY --chown=lambda_build:lambda_build . /go/src/github.com/mstojcevich/lambda-b2-uploader
WORKDIR /go/src/github.com/mstojcevich/lambda-b2-uploader

RUN go build b2_uploader.go

FROM alpine:3.11
RUN addgroup -g 666 -S lambda && \
    adduser -S -s /sbin/nologin -G lambda -u 669 -H -D lambda_file_backuper
USER lambda_file_backuper
COPY --from=build /go/src/github.com/mstojcevich/lambda-b2-uploader/b2_uploader /b2_uploader
CMD ["/b2_uploader"]