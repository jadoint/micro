FROM node:13-alpine AS builder
WORKDIR /build
COPY ./web/package.json .
RUN yarn install
RUN yarn build

# Production build
FROM nginx:mainline-alpine
WORKDIR /home/app
COPY ./deployments/nginx/default.conf.template /etc/nginx/conf.d/default.conf.template
COPY --from=builder /build/build .
