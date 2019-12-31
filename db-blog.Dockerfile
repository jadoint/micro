FROM mysql:8.0
COPY ./deployments/database/blog/conf/my.cnf /etc/mysql/conf.d/my.cnf

# Add the contents of the sql-scripts/ directory to your image.
# All scripts in docker-entrypoint-initdb.d/ are automatically
# executed during container startup.
COPY ./deployments/database/blog/sql-scripts/ /docker-entrypoint-initdb.d/

ENV MYSQL_DATABASE blog
EXPOSE 3306