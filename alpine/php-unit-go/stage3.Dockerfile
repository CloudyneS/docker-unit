FROM ghcr.io/cloudynes/php-unit-go:stage2-8.1-alpine AS BUILDER
FROM ghcr.io/cloudynes/php-extended:8.1-alpine
COPY docker-entrypoint.sh /usr/local/bin/
COPY --from=BUILDER /usr/sbin/unitd /usr/sbin/unitd
COPY --from=BUILDER /usr/sbin/unitd-debug /usr/sbin/unitd-debug
COPY --from=BUILDER /usr/lib/unit/ /usr/lib/unit/
COPY --from=BUILDER /requirements.apk /requirements.apk
COPY --from=BUILDER --chown=unit:unit /app/goapp /app/goapp

RUN ldconfig / && \
    set -x \
    && if [ -f "/tmp/libunit.a" ]; then \
        mv /tmp/libunit.a /usr/lib/$(dpkg-architecture -q DEB_HOST_MULTIARCH)/libunit.a; \
        rm -f /tmp/libunit.a; \
    fi \
    && mkdir -p /var/lib/unit/ \
    && mkdir /docker-entrypoint.d/ \
    && addgroup --system --gid 101 unit \
    && adduser \
        --system \
        --ingroup unit \
        --no-create-home \
        --home /nonexistent \
        --gecos "unit user" \
        --shell /bin/false \
        --uid 101 \
        unit \
    && apk add --no-cache $(cat /requirements.apk) \
    && rm -f /requirements.apk \
    && ln -sf /dev/stdout /var/log/unit.log \
    && mkdir -p /var/lib/unit /app/web \
    && chown -R unit:unit /var/run /run /var/lib/unit /app /app/web 


USER unit

STOPSIGNAL SIGTERM

ENTRYPOINT ["/usr/local/bin/docker-entrypoint.sh"]

CMD ["unitd", "--no-daemon", "--control", "unix:/var/run/control.unit.sock"]
