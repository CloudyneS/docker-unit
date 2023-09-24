ARG SOURCE_VERSION
FROM golang:1.18-buster as INIT

LABEL Version="1.0"
LABEL Maintainer="Cloudyne Systems"
LABEL org.opencontainers.image.source="https://github.com/cloudynes/docker-apps"
LABEL Description="Container for Roots Bedrock based on extended PHP image with Nginx Unit"
LABEL org.opencontainers.image.description="Container for Roots Bedrock based on extended PHP image with Nginx Unit"
LABEL org.opencontainers.image.licenses="MIT"

USER root
WORKDIR /init-go
ADD init-go /init-go

RUN go build -o init-go main.go

FROM ghcr.io/cloudynes/php-unit:${SOURCE_VERSION}-debian

USER root
WORKDIR /app

ENV WP_CLI_CACHE_DIR="/tmp/wpcli/cache"             \
    WP_CLI_CONFIG_PATH="/etc/wpcli/wpcli.conf"      \
    WP_CLI_PACKAGES_DIR="/etc/wpcli/packages"       \
    CD_CONFIG="/init-go/config.json"                \
    COMPOSER_HOME="/etc/composer"

COPY --from=INIT --chown=unit:unit /init-go/init-go /init-go/init-go
COPY --chown=unit:unit ./config.json /docker-entrypoint.d/unit.json
COPY ./init-go/config-sample.json /init-go/config.json
COPY ./composer-config.json /etc/composer/config.json

RUN apt-get update && apt-get -y install git zip unzip mariadb-client && apt-get clean \
    && mv /usr/local/etc/php/php.ini-production /usr/local/etc/php/php.ini \
    && echo -e "memory_limit = 1GB\n" >> /usr/local/etc/php/conf.d/docker-php-memory_limit.ini \
    && sed -i 's/memory_limit = .*/memory_limit = 1G/' /usr/local/etc/php/php.ini \
    && mkdir -p /etc/composer /etc/wpcli/packages \
    && curl -sS https://getcomposer.org/installer | php -- --install-dir=/usr/local/bin --filename=composer \
    && curl -o /usr/local/bin/wp https://raw.githubusercontent.com/wp-cli/builds/gh-pages/phar/wp-cli.phar \
    && chmod +x /usr/local/bin/wp \
    && /usr/local/bin/wp --allow-root package install aaemnnosttv/wp-cli-dotenv-command \
    && composer create-project roots/bedrock --no-interaction --no-dev . \
    && chown -R unit:unit /etc/wpcli /etc/composer /app \
    && rm -rf /var/lib/apt/lists/* /tmp/* /usr/share/doc/* 

USER unit

ENTRYPOINT ["/usr/local/bin/docker-entrypoint.sh"]

CMD [ "unitd", "--no-daemon", "--control", "unix:/var/run/control.unit.sock" ]

EXPOSE 8080