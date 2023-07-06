ARG SOURCE_VERSION
FROM ghcr.io/cloudynes/php-unit:${SOURCE_VERSION}-alpine AS BUILDER

WORKDIR /app

USER root
RUN apk add --no-cache git unzip \
    && mkdir -p /etc/composer \
    && chown -R unit:unit /etc/composer \
    && curl -sS https://getcomposer.org/installer | php -- --install-dir=/usr/local/bin --filename=composer

USER unit

RUN git clone https://github.com/BookStackApp/BookStack.git /app --branch release --single-branch \
    && composer install --no-dev --no-interaction --prefer-dist --optimize-autoloader --no-progress \
    && chmod ug+w -R /app/storage /app/bootstrap/cache /app/public/uploads \
    && rm -rf /app/.git

ARG SOURCE_VERSION
FROM ghcr.io/cloudynes/php-unit:${SOURCE_VERSION}-alpine

USER root
RUN mv "$PHP_INI_DIR/php.ini-production" "$PHP_INI_DIR/php.ini"
USER unit

WORKDIR /app

COPY --chown=unit:unit --from=BUILDER /app /app
COPY --chown=root:root --from=BUILDER "$PHP_INI_DIR/php.ini-production" "$PHP_INI_DIR/php.ini"
COPY --chown=unit:unit ./unit.config.json /docker-entrypoint.d/unit.json
COPY --chown=unit:unit ./docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh

ENV LOG_CHANNEL='stdout'
