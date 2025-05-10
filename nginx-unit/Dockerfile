ARG SOURCE_VERSION=8.2
FROM golang:1.24.3-alpine AS GOGET
# FROM ghcr.io/cloudynes/php-extended:${SOURCE_VERSION}-alpine as BUILDER
FROM test/php-extended:${SOURCE_VERSION}-alpine as BUILDER

LABEL Version="1.0"
LABEL Maintainer="Cloudyne Systems"
LABEL org.opencontainers.image.source="https://github.com/cloudynes/docker-apps"
LABEL Description="PHP and Nginx Unit container for Kubernetes deployment"
LABEL org.opencontainers.image.description="PHP and Nginx Unit container for Kubernetes deployment"
LABEL org.opencontainers.image.licenses="MIT"

COPY --from=GOGET /usr/local/go /usr/local/go
COPY --from=GOGET /go /go

ENV GOTOOLCHAIN=local
ENV GOBIN=/usr/local/go/bin
ENV GOPATH=/go
ENV GOLANG_VERSION=1.24.3
ENV PATH=$PATH:/usr/local/go/bin:/usr/lib/unit:/usr/src/unit/cargo/bin:

# musl/gnu
ENV RUSTUP_VARIANT="musl"
ENV RUST_VERSION=1.83.0
ENV RUSTUP_HOME=/usr/src/unit/rustup
ENV CARGO_HOME=/usr/src/unit/cargo

RUN set -ex && \
    apk add --no-cache \
        ca-certificates \
        git \
        alpine-sdk \
        openssl-dev \
        pcre2-dev \
        dpkg \
        dpkg-dev \
        curl \
        cmake \
    && \
    dpkgArch="$(dpkg --print-architecture)" \
    && case "${dpkgArch##*-}" in \
         amd64) rustArch="x86_64-unknown-linux-${RUSTUP_VARIANT}"; rustupSha256="6aeece6993e902708983b209d04c0d1dbb14ebb405ddb87def578d41f920f56d" ;; \
         arm64) rustArch="aarch64-unknown-linux-${RUSTUP_VARIANT}"; rustupSha256="1cffbf51e63e634c746f741de50649bbbcbd9dbe1de363c9ecef64e278dba2b2" ;; \
         *) echo >&2 "unsupported architecture: ${dpkgArch}"; exit 1 ;; \
       esac \
    && url="https://static.rust-lang.org/rustup/archive/1.27.1/${rustArch}/rustup-init" \
    && curl -L -O "$url" \
    && chmod +x ./rustup-init \
    && ./rustup-init -y --no-modify-path --profile minimal --default-toolchain $RUST_VERSION --default-host ${rustArch} \
    && rm rustup-init \
    && rustup --version \
    && cargo --version \
    && rustc --version \
    && mkdir -p /usr/lib/unit/modules /usr/lib/unit/debug-modules \
    && mkdir -p /usr/src/unit \
    && cd /usr/src/unit \
    && git clone --depth 1 -b 1.34.0-1 https://github.com/nginx/unit \
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
                --njs \
                --otel" \
    && make -j $NCPU -C pkg/contrib .njs \
    && export PKG_CONFIG_PATH=$(pwd)/pkg/contrib/njs/build \
    && ./configure $CONFIGURE_ARGS --cc-opt="$CC_OPT" --ld-opt="$LD_OPT" --modulesdir=/usr/lib/unit/debug-modules --debug \
    && make -j $NCPU unitd \
    && install -pm755 build/sbin/unitd /usr/sbin/unitd-debug \
    && make clean \
    && ./configure $CONFIGURE_ARGS --cc-opt="$CC_OPT" --ld-opt="$LD_OPT" --modulesdir=/usr/lib/unit/modules \
    && make -j $NCPU unitd \
    && install -pm755 build/sbin/unitd /usr/sbin/unitd \
    && make clean \
    && /bin/true \
    && ./configure $CONFIGURE_ARGS_MODULES --cc-opt="$CC_OPT" --modulesdir=/usr/lib/unit/debug-modules --debug \
    && ./configure go --go-path=${GOPATH} \
    && ./configure php \
    && make -j $NCPU go-install-src libunit-install php-install \
    && make clean \
    && ./configure $CONFIGURE_ARGS_MODULES --cc-opt="$CC_OPT" --modulesdir=/usr/lib/unit/modules \
    && ./configure go --go-path=${GOPATH} \
    && ./configure php \
    && make -j $NCPU go-install-src libunit-install php-install \
    && cd \
    && rm -rf /usr/src/unit \
    && for f in /usr/sbin/unitd /usr/lib/unit/modules/*.unit.so; do \
        ldd $f | awk '/=>/{print $(NF-1)}' | while read n; do apk -q info --who-owns $n; done | sed 's/.*owned by //' | sort | uniq >> /requirements.apk \
       done

# FROM ghcr.io/cloudynes/php-extended:${SOURCE_VERSION}-alpine as BUILDER
FROM test/php-extended:${SOURCE_VERSION}-alpine as BUILDER
COPY docker-entrypoint.sh /usr/local/bin/
COPY --from=BUILDER /usr/sbin/unitd /usr/sbin/unitd
COPY --from=BUILDER /usr/sbin/unitd-debug /usr/sbin/unitd-debug
COPY --from=BUILDER /usr/lib/unit/ /usr/lib/unit/
COPY --from=BUILDER /requirements.apk /requirements.apk
COPY --from=BUILDER --chown=unit:unit /app/goapp /app/goapp

COPY --from=GOGET /usr/local/go /usr/local/go
COPY --from=GOGET /go /go

ENV GOTOOLCHAIN=local
ENV GOPATH=/go
ENV GOLANG_VERSION=1.23.2
ENV PATH=$PATH:/usr/local/go/bin
ENV GOCACHE=/tmp

RUN ldconfig / \
    && set -x \
    && mkdir -p /var/lib/unit/ \
    && mkdir -p /docker-entrypoint.d/ \
    && addgroup --system --gid 101 unit \
    && adduser \
        --system \
        --ingroup unit \
        --home /app \
        --no-create-home \
        --gecos "unit user" \
        --shell /bin/false \
        --uid 101 \
        unit \
    && cd /etc/apk/keys && curl -JO https://git.cloudyne.io/api/packages/linux/alpine/key \
    && echo 'https://git.cloudyne.io/api/packages/linux/alpine/all/repository' >> /etc/apk/repositories \
    && apk add --no-cache bash vvv $(cat /requirements.apk) \
    && rm -f /requirements.apk \
    && ln -sf /dev/stdout /var/log/unit.log \
    && mkdir -p /var/lib/unit /app/web && chown -R unit:unit /var/run /run /var/lib/unit

USER unit
STOPSIGNAL SIGTERM

ENTRYPOINT ["/usr/local/bin/docker-entrypoint.sh"]

CMD ["unitd", "--no-daemon", "--control", "unix:/var/run/control.unit.sock"]