FROM mysql:8.0
COPY ./deployments/database/user/conf/my.cnf /etc/mysql/conf.d/my.cnf

# Add the contents of the sql-scripts/ directory to your image.
# All scripts in docker-entrypoint-initdb.d/ are automatically
# executed during container startup.
COPY ./deployments/database/user/sql-scripts/production-user.sql /docker-entrypoint-initdb.d/production-user.sql

ENV MYSQL_DATABASE user
EXPOSE 3306