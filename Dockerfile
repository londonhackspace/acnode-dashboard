# First container - does the build
FROM golang:1.16-alpine AS builder

RUN apk update && apk upgrade && \
    apk add --no-cache git

WORKDIR /build/acnode-dashboard

# copy sources into container build environment
COPY acnode acnode
COPY acserver_api acserver_api
COPY acserverwatcher acserverwatcher
COPY api api
COPY apitypes apitypes
COPY auth auth
COPY config config

COPY *.go ./
COPY go.mod .
COPY go.sum .
COPY .git .git

RUN go build
RUN git rev-list -1 HEAD > version

# Second container - this one actually runs the code
FROM alpine:latest

WORKDIR /opt/acnode-dashboard

COPY --from=builder /build/acnode-dashboard/acnode-dashboard acnode-dashboard
COPY --from=builder /build/acnode-dashboard/version version
COPY static static
COPY templates templates

RUN adduser -S acnodedashboard

USER acnodedashboard

CMD ["./acnode-dashboard"]