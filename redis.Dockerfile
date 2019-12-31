FROM redis:5-alpine
RUN mkdir -p /var/lib/redis/6379
COPY ./deployments/redis/redis.conf /usr/local/etc/redis/redis.conf
EXPOSE 6379
CMD [ "redis-server", "/usr/local/etc/redis/redis.conf" ]