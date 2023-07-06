#!/bin/bash

rm -rf ./php-alpine-original && \
git clone https://github.com/docker-library/php.git ./php-alpine-original && \
sed -ie 's/sqlite-dev \\/sqlite-dev \\ \n      zlib-dev jpeg-dev freetype-dev libwebp-dev icu-dev libpng-dev libressl-dev libzip-dev mariadb-dev krb5-dev \\/' php-alpine-original/8.0/alpine3.16/cli/Dockerfile && \
sed -ie 's/sqlite-dev \\/sqlite-dev \\ \n      zlib-dev jpeg-dev freetype-dev libwebp-dev icu-dev libpng-dev libressl-dev libzip-dev mariadb-dev krb5-dev \\/' php-alpine-original/8.1/alpine3.18/cli/Dockerfile && \
sed -ie 's/--enable-option-checking=fatal \\/--enable-option-checking=fatal \\ \n    --enable-embed \\ \n    --enable-gd --enable-intl --with-pdo-mysql --with-mysqli --enable-opcache --with-zip --enable-bcmath --with-kerberos --with-imap-ssl --with-freetype --with-jpeg --with-webp \\ /' php-alpine-original/8.0/alpine3.16/cli/Dockerfile && \
sed -ie 's/--enable-option-checking=fatal \\/--enable-option-checking=fatal \\ \n    --enable-embed \\ \n    --enable-gd --enable-intl --with-pdo-mysql --with-mysqli --enable-opcache --with-zip --enable-bcmath --with-kerberos --with-imap-ssl --with-freetype --with-jpeg --with-webp \\ /' php-alpine-original/8.1/alpine3.18/cli/Dockerfile

