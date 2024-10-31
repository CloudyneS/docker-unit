ARG SOURCE_VERSION=8.2
FROM golang:1.23-alpine AS GOBUILD

WORKDIR /build

ADD install-itop /build

RUN CGO_ENABLED=0 GOOS=linux go build -o install-itop .

# ARG SOURCE_VERSION=8.2
FROM ghcr.io/cloudynes/php-unit:${SOURCE_VERSION}-alpine AS BUILDER

WORKDIR /app

USER root

COPY --chown=unit:unit ./unit.config.json /docker-entrypoint.d/unit.json
COPY ./docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh

RUN apk update \
    && apk add --no-cache git unzip libxml2-dev openldap-dev graphviz \
    && docker-php-ext-install soap \
    && docker-php-ext-install ldap \
    && curl -sS https://getcomposer.org/installer | php -- --install-dir=/usr/local/bin --filename=composer \
    && mv "$PHP_INI_DIR/php.ini-production" "$PHP_INI_DIR/php.ini" \
    && chmod +x /usr/local/bin/docker-entrypoint.sh && \
    chown -R unit:unit /app

USER unit

# RUN composer create-project combodo/itop:"^3.2.0" /app --no-interaction --no-dev --prefer-dist --no-progress \
#     && git clone https://github.com/combodo/itop-toolkit-community.git /app/toolkit \
#     && chmod ug+w -R /app/conf /app/data \
#     && rm -rf /app/.git \
#     && mkdir -p /app/log \
#     && ln -sf /dev/stdout /app/log/itop.log \
#     && sed -i 's/GLOB_BRACE/0/' /app/setup/setuputils.class.inc.php \
    

COPY --from=GOBUILD --chown=unit:unit /build/install-itop /docker-entrypoint.d/install-itop.sh

RUN git clone https://github.com/combodo/itop.git /app --branch 3.2.0 --single-branch --depth 1 \
    && composer install --no-dev --no-interaction --prefer-dist --optimize-autoloader --no-progress \
    && git clone https://github.com/combodo/itop-toolkit-community.git /app/toolkit \
    && sed -i -E "s/(<\?php)/\1\r\nrequire_once\(\"toolkit_auth.php\"\);/" /app/toolkit/*.php \
    && chmod ug+w -R /app/conf /app/data \
    && rm -rf /app/.git \
    && mkdir -p /app/log \
    && ln -sf /dev/stdout /app/log/itop.log \
    && sed -i 's/GLOB_BRACE/0/' /app/setup/setuputils.class.inc.php \
    && chmod +x /docker-entrypoint.d/install-itop.sh

COPY ./toolkit_auth.php /app/toolkit/toolkit_auth.php

WORKDIR /app
# && sed -i -E "s/(<\?php)/\1\r\nrequire_once\(\"toolkit_auth.php\"\);/" /app/toolkit/ajax_toolkit.php \