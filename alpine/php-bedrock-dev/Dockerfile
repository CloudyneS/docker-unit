ARG SOURCE_VERSION=8.1
FROM ghcr.io/cloudynes/php-unit:${SOURCE_VERSION}-alpine

LABEL Version="1.0"
LABEL Maintainer="Cloudyne Systems"
LABEL org.opencontainers.image.source="https://github.com/cloudynes/docker-apps"
LABEL Description="PHP and Nginx Unit container for Kubernetes deployment"
LABEL org.opencontainers.image.description="PHP and Nginx Unit container for Kubernetes deployment"
LABEL org.opencontainers.image.licenses="MIT"

USER root

WORKDIR /app

ENV WP_CLI_CACHE_DIR="/tmp/wpcli/cache"             \
    WP_CLI_CONFIG_PATH="/etc/wpcli/wpcli.conf"      \
    WP_CLI_PACKAGES_DIR="/etc/wpcli/packages"       \
    CD_CONFIG="/init-go/config.json"                \
    COMPOSER_HOME="/etc/composer"

COPY ./composer-config.json /etc/composer/config.json
COPY --chown=vscode:unit ./config.json /docker-entrypoint.d/unit.json
COPY docker-entrypoint.sh /usr/local/bin/

RUN set -x \
    && apk add --no-cache \
        curl \
        grep \
        zlib-dev \
        curl-dev \
        expat-dev \
        asciidoc \
        xmlto \
        perl-error \
        perl-dev \
        tcl \
        tk \
        g++ \
        python3-dev \
        pcre2-dev \
        git \
        zip \
        unzip \
        mariadb-client \
        ghostscript \
        less \
        nodejs \
        npm \
        autoconf \
        gcc \
        libc-dev \
        linux-headers \
        make \
        bash \
        libstdc++ \
        sudo \
        openssh \
    && pecl install xdebug \
    && echo "vscode ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers \
    && echo -e "[XDebug]\nxdebug.remote_enable = 1\nxdebug.remote_autostart = 1" > $PHP_INI_DIR/conf.d/xdebug.ini \
    && adduser \
        --ingroup unit \
        --home /app \
        --shell /bin/bash \
        --no-create-home \
        --uid 1000 \
        --disabled-password \
        vscode \
    && chown -R vscode:unit /app /var/run /run /var/lib/unit \
    && mv /usr/local/etc/php/php.ini-development /usr/local/etc/php/php.ini \
    && echo -e "memory_limit = 1G \n" >> /usr/local/etc/php/conf.d/docker-php-memory_limit.ini \
    && sed -i 's/memory_limit = .*/memory_limit = 1G/' /usr/local/etc/php/php.ini \
    && mkdir -p /etc/composer /etc/wpcli/packages \
    && curl -sS https://getcomposer.org/installer | php -- --install-dir=/usr/local/bin --filename=composer \
    && curl -o /usr/local/bin/wp https://raw.githubusercontent.com/wp-cli/builds/gh-pages/phar/wp-cli.phar \
    && chmod +x /usr/local/bin/wp \
    && /usr/local/bin/wp --allow-root package install aaemnnosttv/wp-cli-dotenv-command \
    && echo -e 'clear_env = no\r\nvariables_order = "EGPCS"\r\n' >> /usr/local/etc/php/conf.d/additional_vars.ini \
    && chown -R vscode:unit /etc/wpcli /etc/composer /app /var/lib/unit /var/run \
    && chmod +x /usr/local/bin/docker-entrypoint.sh \
    && rm -rf /var/lib/apt/lists/* /tmp/* /usr/share/doc/*

USER vscode

RUN composer create-project roots/bedrock --no-interaction . && \
    cp /app/.env.example /app/.env && \
    mkdir /app/.vscode

COPY --chown=vscode:unit vscode-launch.json /app/.vscode/launch.json


ENV ADMIN_USER=admin \
    ADMIN_PASSWORD=password \
    ADMIN_EMAIL="admin@host.local" \
    SITE_TITLE="WordPress Development Site" \
    SITE_URL="http://localhost:8080" \
    DB_PREFIX="wpdev_" \
    DATABASE_URL="mysql://wp_user:wp_password@db:3306/wordpress"



ENTRYPOINT ["/usr/local/bin/docker-entrypoint.sh"]

CMD [ "unitd", "--no-daemon", "--control", "unix:/var/run/control.unit.sock" ]

EXPOSE 8080
