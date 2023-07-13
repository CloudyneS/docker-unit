ARG SOURCE_IMAGE
FROM ghcr.io/cloudynes/${SOURCE_IMAGE}

USER root
ARG DEBUG_PACKAGES
RUN apk add --no-cache ${DEBUG_PACKAGES}
USER unit

