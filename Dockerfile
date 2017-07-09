# build stage
FROM golang:alpine AS build-env
RUN apk add --no-cache git
RUN apk add --no-cache sqlite
RUN apk add --no-cache g++
ADD . /src
RUN cd /src && go get -d -v ./... && go build -o goapp

# final stage
FROM alpine
ENV SQLITE_PATH /data/saved.sqlite
VOLUME /data
WORKDIR /app
COPY --from=build-env /src/goapp /app/
COPY --from=build-env /src/saved.sqlite /data/
ENTRYPOINT ./goapp
