FROM node:13-alpine AS builder
WORKDIR /build/app
COPY ./web .
# Disable sourcemaps to fix issue:
# FATAL ERROR: Ineffective mark-compacts near heap limit Allocation failed - JavaScript heap out of memory
ENV GENERATE_SOURCEMAP false
RUN yarn install
RUN yarn build

# Production build
FROM nginx:mainline-alpine
WORKDIR /home/app
COPY ./deployments/nginx/default.prod.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /build/app .
