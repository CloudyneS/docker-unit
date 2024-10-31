ARG SOURCE_VERSION=8.1
FROM golang:1.23.2-alpine AS GOGET
FROM ghcr.io/cloudynes/php-extended:${SOURCE_VERSION}-alpine as BUILDER

LABEL Version="1.0"
LABEL Maintainer="Cloudyne Systems"
LABEL org.opencontainers.image.source="https://github.com/cloudynes/docker-apps"
LABEL Description="PHP and Nginx Unit container for Kubernetes deployment"
LABEL org.opencontainers.image.description="PHP and Nginx Unit container for Kubernetes deployment"
LABEL org.opencontainers.image.licenses="MIT"

COPY --from=GOGET /usr/local/go /usr/local/go
COPY --from=GOGET /go /go
COPY ./goapp /src/goapp

ENV GOTOOLCHAIN=local
ENV GOBIN=/usr/local/go/bin
ENV GOPATH=/go
ENV GOLANG_VERSION=1.23.2

ENV PATH=$PATH:/usr/local/go/bin:/usr/lib/unit

RUN set -ex \
    && export PATH=$PATH:/usr/local/go/bin \
    && apk add --no-cache \
        git \
        alpine-sdk \
        openssl-dev \
        pcre2-dev \
        dpkg dpkg-dev \
    && mkdir -p /usr/lib/unit/modules /usr/lib/unit/debug-modules \
    && mkdir -p /usr/src/unit \
    && cd /usr/src/unit \
    && git clone --depth 1 -b 1.33.0-1 https://github.com/nginx/unit \
    && cd unit \
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
    && ./configure go --go-path=${GOPATH} \
    && ./configure php \
    && make -j $NCPU go-install-src libunit-install php-install \
    && make clean \
    && ./configure $CONFIGURE_ARGS_MODULES --cc-opt="$CC_OPT" --modules=/usr/lib/unit/modules \
    && ./configure go --go-path=${GOPATH} \
    && ./configure php \
    && make -j $NCPU go-install-src libunit-install php-install \
    && rm -rf /usr/src/unit \
    && ldd /usr/sbin/unitd | awk '/=>/{print $(NF-1)}' | while read n; do apk -q info --who-owns $n; done | sed 's/.*owned by //' | sort | uniq > /requirements.apk \
    && cp /usr/lib/x86_64-linux-musl/libunit.a /usr/lib/libunit.a

