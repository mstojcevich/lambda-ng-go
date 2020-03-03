FROM golang:1.13-alpine AS build

# Consistency across Lambda containers
ENV LANG en_US.UTF-8

RUN addgroup -g 666 -S lambda_build && \
    adduser -S -s /sbin/nologin -G lambda_build -u 666 -D lambda_build
USER lambda_build

COPY --chown=lambda_build:lambda_build . /go/src/github.com/mstojcevich/lambda-db-backupper
WORKDIR /go/src/github.com/mstojcevich/lambda-db-backupper

RUN go build db_backupper.go

FROM alpine:3.11
RUN apk add --no-cache postgresql-client && \
    addgroup -g 670 -S lambda_db_backupper && \
    adduser -S -s /sbin/nologin -G lambda_db_backupper -u 670 -H -D lambda_db_backupper
USER lambda_db_backupper
COPY --from=build /go/src/github.com/mstojcevich/lambda-db-backupper/db_backupper /db_backupper
CMD ["/db_backupper"]