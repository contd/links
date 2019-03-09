# Builder image
FROM golang AS builder

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o goapp

# Production image
FROM scratch
ENV SQLITE_PATH=/data/saved.sqlite
VOLUME /data
COPY --from=builder /app/goapp /app/
COPY --from=builder /app/web/favicon.svg /app/web/
COPY --from=builder /app/web/links.html /app/web/
COPY --from=builder /app/data/saved.sqlite /data/

EXPOSE 5555
ENTRYPOINT [ "/app/goapp" ]