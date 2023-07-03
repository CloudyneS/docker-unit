ARG SOURCE_VERSION=8.2
FROM php:${SOURCE_VERSION}-cli

LABEL Version="1.0"
LABEL Maintainer="Cloudyne Systems"
LABEL org.opencontainers.image.source="https://github.com/cloudynes/docker-apps"
LABEL Description="PHP Base container extended with imap,gd,intl,pdo,mysqli,opcache,imap,zip,bcmath"
LABEL org.opencontainers.image.description="PHP Base container extended with imap,gd,intl,pdo,mysqli,opcache,imap,zip,bcmath"
LABEL org.opencontainers.image.licenses="MIT"

ENV PHPIZE_DEPS \
    zlib1g-dev libfreetype6-dev \
    curl libjpeg-dev libwebp-dev \
    libicu-dev libpng-dev \
    libssl-dev libzip-dev \
    default-libmysqlclient-dev \
    libc-client-dev libkrb5-dev \
	gcc g++ libc-dev make
        
RUN set -eux; \
	\
	savedAptMark="$(apt-mark showmanual)"; \
	apt-get update; \
	apt-get install -y --no-install-recommends \
		$PHPIZE_DEPS \
	; \
	\
    docker-php-ext-configure imap --with-kerberos --with-imap-ssl; \
    docker-php-ext-configure gd --with-freetype --with-jpeg --with-webp; \
	docker-php-ext-install -j "$(nproc)" \
        gd \
        intl \
        pdo_mysql \
        mysqli \
        opcache \
        imap \
        zip \
        bcmath \
        ; \
    # reset apt-mark's "manual" list so that "purge --auto-remove" will remove all build dependencies
	apt-mark auto '.*' > /dev/null; \
	[ -z "$savedAptMark" ] || apt-mark manual $savedAptMark; \
	find /usr/local -type f -executable -exec ldd '{}' ';' \
		| awk '/=>/ { print $(NF-1) }' \
		| sort -u \
		| xargs -r dpkg-query --search \
		| cut -d: -f1 \
		| sort -u \
		| xargs -r apt-mark manual \
	; \
	apt-get purge gcc -y --auto-remove -o APT::AutoRemove::RecommendsImportant=false; \
	rm -rf /var/lib/apt/lists/*; \
	\
# update pecl channel definitions https://github.com/docker-library/php/issues/443
	pecl update-channels; \
	rm -rf /tmp/pear ~/.pearrc; \
# Set higher memory limit as default
    echo 'memory_limit = 1024M' >> /usr/local/etc/php/conf.d/docker-php-memlimit.ini; \
	\
# Choose the php.ini-production as default
	cp "$PHP_INI_DIR/php.ini-production" "$PHP_INI_DIR/php.ini"; \
# smoke test
	php --version
