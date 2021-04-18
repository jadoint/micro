FROM nginx:mainline-alpine
WORKDIR /home/app

COPY ./deployments/nginx/default.dev.conf /etc/nginx/conf.d/default.conf
COPY ./private/ssl/site.pem /etc/nginx/ssl/site.pem
COPY ./private/ssl/site.key /etc/nginx/ssl/site.key

# Create self-signed certs below if not copying SSL certs
# Note: req.conf defines alternate_names which allows for
# alternative domains and subdomains for Chrome.
# RUN mkdir -p /etc/nginx/ssl
# RUN apk add --update openssl && \
#   rm -rf /var/cache/apk/*
# RUN openssl req -config ./private/ssl/req.conf -new -x509 -sha256 -newkey rsa:2048 -nodes -keyout ./private/ssl/site.key -days 3650 -out ./private/ssl/site.pem -subj "/C=US/ST=WA/L=Seattle/O=David Ado/OU=IT/CN=David Ado"