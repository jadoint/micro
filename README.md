# Micro

Production-ready template for Go + React / Redux microservice applications. Sample website running micro: [https://www.davidado.com/](https://www.davidado.com/)

## Description

Micro is a production-ready microservices template.

`/cmd` is the entrypoint for each service so you can run `go run cmd/blog/main.go` and `go run cmd/user/main.go` during development.

`/pkg` contains all importable packages:

- `/blog` blog service
- `/clean` strips user-submitted content of XSS vulnerabilities
- `/conn` struct containing database and cache connection clients
- `/contextkey` keys for use with context
- `/cookie` cookie helper
- `/db` database and query helper
- `/env` .env configuration loader
- `/errutil` error helper
- `/fmtdate` date helper
- `/hash` argon2 hash helper
- `/logger` custom error logger
- `/msg` sends json messages to the end user
- `/onesignal` onesignal integration for push messaging
- `/paginate` pagination helper
- `/token` JWT helper
- `/user` user service
- `/useragent` user-agent helper
- `/validate` struct validator
- `/visitor` creates and contains details about a visitor
- `/words` word count and word censor

`/web` contains frontend code like React components and public files.

## Installation

### Development

Copy `docker-compose-sample.yml` to `docker-compose.yml` and configure settings according to your environment then run `docker-compose up -d`.

### Production

Refer to `docker-swarm-init-deploy-sample.sh` for the initial server setup and `docker-swarm-deploy-sample.sh` for regular deployment.

### Edge Server (CDN) Caching

Many API calls are .json requests such as `/api/v1/blog/1/blog_1_20190911180259.json` to enable caching for large text data. You can configure your server to set a long max-age header for `.json` responses which tells a user's browser to cache that post for future visits. Better yet, you can place your application behind a reverse proxy like Cloudflare and set a Page Rule like `www.yoursite.com/*.json` => `Cache Level: Cache Everything` and `Edge Cache TTL: a month` so that all requests for that resource are delivered by their CDN. Doing this for one of my popular applications reduced my cloud egress bandwidth costs by 60%.

## Addendum

I went the microservices route primarily because in my current monoliths, database schema changes and large updates to big tables (> 50 GB) cause slave lag that affects all other unrelated services. Separating domains into their own databases and codebases was a way to reduce side effects to other services so that a disruption in one does not affect the other (i.e. if the user service goes down, you can still view the blog posts even if you won't see the author).
