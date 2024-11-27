# docker-unit
Container images with various applications for Nginx Unit

## Base Images
There are two types of images, alpine and debian. Alpine is the most maintained image, since there are not a lot of viable alternatives to this.

## PHP Versions
Everything from 8.0 and up is actively maintained up until deprecation, after this it's on an as-needed basis.

## Images
### PHP Extended
PHP recompiled with the allow_embed flag and some additional extensions for web:
- PD
- Intl
- PDO
- Mysqli
- Opcache
- Zip
- Bcmath
- Kerberos
- Imap
- Webp

Dynamically generated based on the official PHP Alpine image, with some added flags to compilation, thus restricted to the same version combinations of PHP/Alpine that are used there.

### PHP Unit
Nginx Unit on top of the PHP Extended image, based on the official unit container with some additions (e.g. NJS)

### PHP Bedrock
Roots Bedrock installed in /app on the PHP Unit image. Also includes an initializer for some common on-start tasks (DB/File import from URL or path, setting the theme, changing ownership of files, conversion to webp (requires cloudyne-extras plugin))

### PHP Itop
PHP image for Combodo Itop based on the PHP Unit Container installed with composer support. Additional extensions can be installed by creating a new image and installing with composer.

### PHP Bookstack
Bookstack based on the PHP Unit container
x
