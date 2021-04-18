FROM golang:alpine
WORKDIR /home/app

RUN mkdir -p /home/ssl
COPY private/ssl/site.pem /home/ssl/site.pem
COPY private/ssl/site.key /home/ssl/site.key

# build-base includes gcc which kafka needs for builds
RUN apk add build-base

# Fetch dependencies first; they are less susceptible to change on every build
# and will therefore be cached for speeding up the next build
COPY ./go.mod ./go.sum ./
RUN go mod download

# "-tags musl" needed for kafka when doing alpine builds
RUN go get -u -tags musl github.com/cosmtrek/air
ENTRYPOINT air -c .air.user.toml
EXPOSE 9000