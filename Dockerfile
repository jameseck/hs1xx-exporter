FROM golang:alpine AS build-env

RUN apk add git
WORKDIR /src
COPY src/go.mod .
COPY src/go.sum .
RUN go mod download

COPY src/ .
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o hs1xx-exporter .

FROM busybox

RUN adduser hs1xx-exporter -D -h /app

COPY --from=build-env /src/hs1xx-exporter /

ENV HOME=/app

ENTRYPOINT ["./hs1xx-exporter"]

