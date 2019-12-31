FROM mysql:8.0

ENV MYSQL_DATABASE blog

# Percona package to to install
# ENV PACKAGE percona-xtrabackup-24

# Dependencies
# RUN apt-get update && apt-get install -y wget

# Install Percona apt repository and Percona Xtrabackup
# RUN wget https://repo.percona.com/apt/percona-release_0.1-4.stretch_all.deb && \
#   dpkg -i percona-release_0.1-4.stretch_all.deb && \
#   apt-get update && \
#   apt-get install -y $PACKAGE

# Create the backup destination
# RUN mkdir -p /backups

# Allow mountable backup path
# VOLUME ["/backups"]

# Add the contents of the sql-scripts/ directory to your image.
# All scripts in docker-entrypoint-initdb.d/ are automatically
# executed during container startup.
# COPY ./sql-scripts/ /docker-entrypoint-initdb.d/

# Build this image.
# docker build -t dockerhub_username/micro-blog:0.1 .

# Start the MySql container.
# -v option mounts a host volume to the container
# where [host_dir]:[container_dir]
# docker run -d -p 3306:3306 --name micro-blog -e MYSQL_ROOT_PASSWORD=yourrootpassword -v /micro-blog/data:/var/lib/mysql -v /micro-blog/backups:/backups -v /micro-blog/conf:/etc/mysql/conf.d micro-blog:0.1

# Verify database using:
# docker exec -it micro-blog bash
# mysql > show databases;

EXPOSE 3306