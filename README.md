# micro

Production-ready template for Go + React / Redux microservice applications. Sample website running micro: [https://www.davidado.com/](https://www.davidado.com/)

## Why?

This is meant to be a production-ready template for more complex applications so it may seem needlessly complicated just for a blog. I went the microservices route primarily because in my current monoliths, database schema changes and large updates to big tables (> 50 GB) cause slave lag that affects all other unrelated services. Separating domains into their own databases and codebases was a way to reduce side effects to other services so that a disruption in one does not affect the other (e.g. if the user service goes down, you can still view the blog posts but you just won't see the author).

Shared dependencies are located in `/pkg` and are currently a weak point in this project but versioning should help mitigate some of that risk (maybe Go modules when I get around to it).

One peculiarity you might notice is the call to two (or more) APIs when one can do like a call to `/api/v1/blog/1` and then `/api/v1/blog/1/blog_1_20190911180259.json`. I've placed requests for content that does not change very often under the `.json` extension to enable caching for large text data. You can configure your server to set a long max-age header for `.json` responses which tells a user's browser to cache that post for future visits. Better yet, you can place your application behind a reverse proxy like Cloudflare and set a Page Rule like `www.yoursite.com/*.json` => `Cache Level: Cache Everything` and `Edge Cache TTL: a month` so that all requests for that resource are delivered by their CDN. Doing this for one of my popular applications reduced my cloud egress bandwidth costs by 60%.

## Installation

### Database Setup

**Docker:** In `/deployments/database`, copy `docker-compose-sample.yml` to `docker-compose.yml` and configure settings according to your environment then run `docker-compose up -d`. Note that multiple ports are used depending on the database (e.g. 3400 for user, 3401 for blog).

**Manual setup:** In the blog and user directories in `/deployments/database`, find the setup SQL scripts in the `sql-scripts` directory and execute those to create the `blog` and `user` databases. You can reference `/deployments/database/blog/conf/my.cnf` and `/deployments/database/user/conf/my.cnf` for configuring MySql.

### Cache Setup

**Docker:** In `/deployments/cache`, copy `docker-compose-sample.yml` to `docker-compose.yml` and configure settings according to your environment then run `docker-compose up -d`.

**Manual setup:** Simply install Redis on your platform. You can reference `/deployments/cache/redis.conf` for this installation.

### Go Setup

Copy `deployments/bin/blog/.sample.env.production` to `deployments/bin/blog/.env.production` and configure settings according to your environment (do the same for `deployments/bin/user/.sample.env.production`).

Copy `deploy-services-sample.sh` to `deploy-services.sh` and configure settings according to your environment. Running this script builds the necessary binaries, uploads `/deployments/bin` to your server, and restarts each respective service.

Debian/Ubuntu setup: To install the binaries as a service to be managed by `systemd`, copy `deployments/bin/blog/blog.service-sample` to `deployments/bin/blog/blog.service` and configure it for your environment - mainly change `youruser` to your Linux user account name (do the same for `deployments/bin/user/user.service-sample`). Place each respective `*.service` file in `/lib/systemd/system` then run `chmod 755 blog.service`, `systemctl enable blog.service`, and `systemctl start blog.service` (do the same for the user service).

### React Setup

Copy `deploy-web-sample.sh` to `deploy-web.sh` and configure settings according to your environment. Running this scripts creates a production-ready build of the React application (runs `yarn build` in `/web`) and uploads it to your server.

## Development

`/cmd` is the entrypoint for each service so you can run `go run cmd/blog/main.go` and `go run cmd/user/main.go` during development.

`/pkg` contains all internal packages that can be imported (i.e. you're going to be writing most of your Go code in here).

`/web` contains frontend code like React components and public files.
