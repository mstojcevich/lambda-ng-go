FROM alpine:3.11

# Consistency across Lambda containers
ENV LANG en_US.UTF-8

# clamav-libunrar is for freshclam. Maybe FROM alpine:3.11 after fetch is done to get rid of it
RUN apk --no-cache add \
    clamav-daemon \
    clamav-libunrar

# lambda_clamav and lambda_clamav_downloader are different users so that lambda_clamav doesn't need to have
# write access to the database files.
# Make sure to make any changes to freshclam's dockerfile too!
RUN addgroup -g 667 -S lambda_clamav && \
    adduser -S -s /sbin/nologin -G lambda_clamav -u 667 -H -D lambda_clamav && \
    adduser -S -s /sbin/nologin -G lambda_clamav -u 668 -H -D lambda_clamav_downloader && \
    mkdir -p /data/clamdb && chown lambda_clamav_downloader:lambda_clamav /data/clamdb && \
    chmod 750 /data/clamdb

COPY --chown=lambda_clamav:lambda_clamav ./clamd.conf /clamconf/clamd.conf
RUN chmod 400 /clamconf/clamd.conf

# We don't get a freshclam daemon since this is a container. Download the files manually.
# TODO Maybe run freshclam in another container? Then send SIGUSR2 somehow to make it reload the sig databases
USER lambda_clamav_downloader
RUN freshclam --stdout --log=/dev/stdout --daemon-notify=/dev/null --datadir=/data/clamdb && chmod 640 /data/clamdb/*

USER lambda_clamav

CMD ["/usr/sbin/clamd", "--config-file=/clamconf/clamd.conf"]