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
COPY bootstrapper bootstrapper
COPY logwatcher logwatcher
COPY usagelogs usagelogs

COPY *.go ./
COPY go.mod .
COPY go.sum .
COPY .git .git

RUN go build
RUN cd bootstrapper && go build
RUN cd logwatcher && go build
RUN git rev-list -1 HEAD > version

# Second container - build frontend
FROM node:16.0-alpine as nodebuilder

WORKDIR /build/acnode-dashboard/frontend

COPY frontend /build/acnode-dashboard/frontend

RUN npm install
RUN npm build && cd dist && tar -cf ../bundle.tar static index.html

# Third container - this one actually runs the code
FROM alpine:latest

WORKDIR /opt/acnode-dashboard

COPY --from=builder /build/acnode-dashboard/acnode-dashboard acnode-dashboard
COPY --from=builder /build/acnode-dashboard/bootstrapper/bootstrapper bootstrapper
COPY --from=builder /build/acnode-dashboard/version version
COPY --from=nodebuilder /build/acnode-dashboard/frontend/bundle.tar frontend.tar
COPY static static
COPY templates templates

RUN mkdir static/newfrontend && cd static/newfrontend && tar -xf ../../frontend.tar && rm -f ../../frontend.tar

RUN adduser -S acnodedashboard

USER acnodedashboard

CMD ["./acnode-dashboard"]