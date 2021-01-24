FROM alpine:3.13

# Consistency across Lambda containers
ENV LANG en_US.UTF-8

RUN apk --no-cache add bash freshclam clamav-libunrar

# Same group and downloader user as clamav's dockerfile so that clamav can have read-only
# access to the DB file.
RUN addgroup -g 667 -S lambda_clamav && \
    adduser -S -s /sbin/nologin -G lambda_clamav -u 668 -H -D lambda_clamav_downloader && \
    mkdir -p /data/clamdb && chown lambda_clamav_downloader:lambda_clamav /data/clamdb && \
    chmod 750 /data/clamdb

USER lambda_clamav_downloader

CMD ["/usr/bin/env", "bash", "-c", "/usr/bin/freshclam --stdout --checks=4 --log=/dev/stdout --daemon-notify=/dev/null --datadir=/data/clamdb -d --pid=/tmp/freshclam.pid && while [ -e /proc/$(cat /tmp/freshclam.pid) ]; do sleep 60; done"]
