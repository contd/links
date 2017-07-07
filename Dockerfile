# build stage
FROM golang:alpine AS build-env
ADD . /src
RUN cd /src && go build -o goapp

# final stage
FROM alpine
VOLUME /data
WORKDIR /app
COPY --from=build-env /src/goapp /app/
COPY --from=build-env /src/links.sqlite /data/
ENTRYPOINT ./goapp
