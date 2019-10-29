FROM redis:5
LABEL Name=redis-micro
RUN mkdir -p /var/lib/redis/6379
COPY redis.conf /usr/local/etc/redis/redis.conf
CMD [ "redis-server", "/usr/local/etc/redis/redis.conf" ]
EXPOSE 6379