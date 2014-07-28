FROM darron/nginx
# FROM nginx_php_base

# Build your new site here.
# RUN echo "Site 123 ABC" > /var/www/index.php
# RUN echo "<?php phpinfo(); ?>" >> /var/www/index.php
ADD index.html /var/www/index.html
ADD index.php /var/www/index.php
ADD unchunked.php /var/www/unchunked.php

EXPOSE 80

CMD php5-fpm -D && nginx