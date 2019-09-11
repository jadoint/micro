FROM redis:5
LABEL Name=micro-redis Version=0.0.1
RUN mkdir -p /var/lib/redis/6379
COPY redis.conf /usr/local/etc/redis/redis.conf
CMD [ "redis-server", "/usr/local/etc/redis/redis.conf" ]
EXPOSE 6379