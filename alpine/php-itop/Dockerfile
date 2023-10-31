ARG SOURCE_VERSION
FROM ghcr.io/cloudynes/php-unit:${SOURCE_VERSION}-alpine AS BUILDER

WORKDIR /app

USER root

COPY --chown=unit:unit ./unit.config.json /docker-entrypoint.d/unit.json
COPY ./docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh

RUN apk update \
    && apk add --no-cache git unzip libxml2-dev openldap-dev \
    && docker-php-ext-install soap \
    && docker-php-ext-install ldap \
    && curl -sS https://getcomposer.org/installer | php -- --install-dir=/usr/local/bin --filename=composer \
    && mv "$PHP_INI_DIR/php.ini-production" "$PHP_INI_DIR/php.ini" \
    && chmod +x /usr/local/bin/docker-entrypoint.sh

USER unit

RUN git clone https://github.com/combodo/itop.git /app --branch 3.1.0-3 --single-branch --depth 1 \
    && composer install --no-dev --no-interaction --prefer-dist --optimize-autoloader --no-progress \
    && chmod ug+w -R /app/conf /app/data \
    && rm -rf /app/.git \
    && mkdir -p /app/log \
    && ln -sf /dev/stdout /app/log/itop.log \
    && sed -i 's/GLOB_BRACE/0/' /app/setup/setuputils.class.inc.php


WORKDIR /app

