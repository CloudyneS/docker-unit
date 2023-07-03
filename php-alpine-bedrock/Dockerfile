ARG SOURCE_VERSION
# FROM ghcr.io/cloudynes/php-unit:${SOURCE_VERSION} AS BEDROCK
FROM golang:1.18-buster as INIT

LABEL Version="1.0"
LABEL Maintainer="Cloudyne Systems"
LABEL org.opencontainers.image.source="https://github.com/cloudynes/docker-apps"
LABEL Description="Container for Roots Bedrock based on extended PHP image with Nginx Unit"
LABEL org.opencontainers.image.description="Container for Roots Bedrock based on extended PHP image with Nginx Unit"
LABEL org.opencontainers.image.licenses="MIT"

# Set the working directory and installation user
# USER root

# WORKDIR /app

# Install GIT and Composer
# RUN apt-get update && apt-get -y install git && apt-get clean \ 
#     && chown -R unit:unit /app \
#     && curl -sS https://getcomposer.org/installer | php -- --install-dir=/usr/local/bin --filename=composer \
#     && composer create-project --no-interaction --no-dev roots/bedrock .


USER root
WORKDIR /init-go
ADD init-go /init-go

RUN go build -o init-go main.go

FROM ghcr.io/cloudynes/php-alpine-unit:${SOURCE_VERSION}

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


RUN apk add --no-cache git unzip mariadb-client && \
    mkdir -p /etc/composer /etc/wpcli/packages && \
    curl -sS https://getcomposer.org/installer | php -- --install-dir=/usr/local/bin --filename=composer && \
    curl -o /usr/local/bin/wp https://raw.githubusercontent.com/wp-cli/builds/gh-pages/phar/wp-cli.phar && \
    chmod +x /usr/local/bin/wp && \
    /usr/local/bin/wp --allow-root package install aaemnnosttv/wp-cli-dotenv-command && \
    composer create-project roots/bedrock --no-interaction --no-dev . && \
    chown -R unit:unit /etc/wpcli /etc/composer /app && \
    rm -rf /var/lib/apt/lists/* /tmp/* /usr/share/doc/*

USER unit

ENTRYPOINT ["/usr/local/bin/docker-entrypoint.sh"]

CMD [ "unitd", "--no-daemon", "--control", "unix:/var/run/control.unit.sock" ]

EXPOSE 8080