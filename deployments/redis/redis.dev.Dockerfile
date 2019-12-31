FROM redis:5
RUN mkdir -p /var/lib/redis/6379
COPY redis.dev.conf /usr/local/etc/redis/redis.conf
CMD [ "redis-server", "/usr/local/etc/redis/redis.conf" ]