FROM node:13 AS builder
WORKDIR /build/app

COPY ./web/package.json .
RUN yarn add && yarn build

# Production build
FROM nginx:mainline-alpine
WORKDIR /home/app

COPY ./deployments/nginx/default.conf.template /etc/nginx/conf.d/default.conf.template

COPY --from=builder /build/app .