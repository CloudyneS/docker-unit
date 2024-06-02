ARG SOURCE_VERSION=8.1
FROM ghcr.io/cloudynes/php-extended:${SOURCE_VERSION}-alpine as BUILDER

LABEL Version="1.0"
LABEL Maintainer="Cloudyne Systems"
LABEL org.opencontainers.image.source="https://github.com/cloudynes/docker-apps"
LABEL Description="PHP and Nginx Unit container for Kubernetes deployment"
LABEL org.opencontainers.image.description="PHP and Nginx Unit container for Kubernetes deployment"
LABEL org.opencontainers.image.licenses="MIT"

RUN set -ex \
    && apk add --no-cache \
        mercurial \
        alpine-sdk \
        openssl-dev \
        pcre2-dev \
        dpkg dpkg-dev \
    && mkdir -p /usr/lib/unit/modules /usr/lib/unit/debug-modules \
    && hg clone https://hg.nginx.org/unit \
    && cd unit \
    && hg up 1.32.0-1 \
    && NCPU="$(getconf _NPROCESSORS_ONLN)" \
    && DEB_HOST_MULTIARCH="$(dpkg-architecture -q DEB_HOST_MULTIARCH)" \
    && CC_OPT="$(DEB_BUILD_MAINT_OPTIONS="hardening=+all,-pie" DEB_CFLAGS_MAINT_APPEND="-Wp,-D_FORTIFY_SOURCE=2 -fPIC" dpkg-buildflags --get CFLAGS)" \
    && LD_OPT="$(DEB_BUILD_MAINT_OPTIONS="hardening=+all,-pie" DEB_LDFLAGS_MAINT_APPEND="-Wl,--as-needed -pie" dpkg-buildflags --get LDFLAGS)" \
    && CONFIGURE_ARGS_MODULES="--prefix=/usr \
                --statedir=/var/lib/unit \
                --control=unix:/var/run/control.unit.sock \
                --runstatedir=/var/run \
                --pid=/var/run/unit.pid \
                --logdir=/var/log \
                --log=/var/log/unit.log \
                --tmpdir=/var/tmp \
                --user=unit \
                --group=unit \
                --openssl \
                --libdir=/usr/lib/$DEB_HOST_MULTIARCH" \
    && CONFIGURE_ARGS="$CONFIGURE_ARGS_MODULES \
                --njs" \
    && make -j $NCPU -C pkg/contrib .njs \
    && export PKG_CONFIG_PATH=$(pwd)/pkg/contrib/njs/build \
    && ./configure $CONFIGURE_ARGS --cc-opt="$CC_OPT" --ld-opt="$LD_OPT" --modules=/usr/lib/unit/debug-modules --debug \
    && make -j $NCPU unitd \
    && install -pm755 build/sbin/unitd /usr/sbin/unitd-debug \
    && make clean \
    && ./configure $CONFIGURE_ARGS --cc-opt="$CC_OPT" --ld-opt="$LD_OPT" --modules=/usr/lib/unit/modules \
    && make -j $NCPU unitd \
    && install -pm755 build/sbin/unitd /usr/sbin/unitd \
    && make clean \
    && ./configure $CONFIGURE_ARGS_MODULES --cc-opt="$CC_OPT" --modules=/usr/lib/unit/debug-modules --debug \
    && ./configure php \
    && make -j $NCPU php-install \
    && make clean \
    && ./configure $CONFIGURE_ARGS_MODULES --cc-opt="$CC_OPT" --modules=/usr/lib/unit/modules \
    && ./configure php \
    && make -j $NCPU php-install \
    && cd \
    && rm -rf /usr/src/unit \
    && ldd /usr/sbin/unitd | awk '/=>/{print $(NF-1)}' | while read n; do apk -q info --who-owns $n; done | sed 's/.*owned by //' | sort | uniq > /requirements.apk

FROM ghcr.io/cloudynes/php-extended:${SOURCE_VERSION}-alpine
COPY docker-entrypoint.sh /usr/local/bin/
COPY --from=BUILDER /usr/sbin/unitd /usr/sbin/unitd
COPY --from=BUILDER /usr/sbin/unitd-debug /usr/sbin/unitd-debug
COPY --from=BUILDER /usr/lib/unit/ /usr/lib/unit/
COPY --from=BUILDER /requirements.apk /requirements.apk

RUN ldconfig /

RUN set -x \
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
    && ln -sf /dev/stdout /var/log/unit.log && \
    mkdir -p /var/lib/unit && chown -R unit:unit /var/run /run /var/lib/unit

USER unit

STOPSIGNAL SIGTERM

ENTRYPOINT ["/usr/local/bin/docker-entrypoint.sh"]

CMD ["unitd", "--no-daemon", "--control", "unix:/var/run/control.unit.sock"]
