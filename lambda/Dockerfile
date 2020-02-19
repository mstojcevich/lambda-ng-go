FROM golang:1.13-alpine AS build

# Consistency across Lambda containers
ENV LANG en_US.utf8

RUN apk --no-cache add \
    git \
    gcc libc-dev pkgconfig \
    vips vips-dev

RUN addgroup -g 666 -S lambda_build && \
    adduser -S -s /sbin/nologin -G lambda_build -u 666 -D lambda_build
USER lambda_build

COPY --chown=lambda_build:lambda_build . /go/src/github.com/mstojcevich/lambda-ng-go
WORKDIR /go/src/github.com/mstojcevich/lambda-ng-go

# RUN go build -ldflags "-linkmode external -extldflags -static" lambda.go
RUN go build lambda.go

# TODO if we can statically link vips (or switch to a go lib) then we can use "FROM scratch" to save space
FROM alpine:3.11
RUN apk --no-cache add vips
RUN addgroup -g 666 -S lambda && \
    adduser -S -s /sbin/nologin -G lambda -u 666 -H -D lambda
USER lambda
COPY --from=build /go/src/github.com/mstojcevich/lambda-ng-go/lambda /lambda
COPY --chown=lambda ./html /html
COPY --chown=lambda ./static /static
CMD ["/lambda"]