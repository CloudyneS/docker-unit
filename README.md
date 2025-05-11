# docker-unit
Container images with various applications for Nginx Unit

## Base PHP Image
Base PHP image with some added extensions and flags for use with Nginx Unit. This image is used as a base for other images. Additional features compared to the base PHP image:

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

## Base Nginx Unit Image
This is Nginx Unit built on top of the extended base PHP image (alpine). This image is built for use with PHP and Go applications, and includes NJS.
The Bookworm image is basically the same as the official unit image, but with the PHP packages above added.

### PHP Bedrock
Roots Bedrock installed in /app on the PHP Unit image. Also includes an initializer for some common on-start tasks (DB/File import from URL or path, setting the theme, changing ownership of files, conversion to webp (requires cloudyne-extras plugin))

### PHP Itop
PHP image for Combodo Itop based on the PHP Unit Container installed with composer support. Additional extensions can be installed by creating a new image and installing with composer.

### PHP Bookstack
Bookstack based on the PHP Unit container
