# First container - does the build
FROM golang:1.16-alpine AS builder

RUN apk update && apk upgrade && \
    apk add --no-cache git

WORKDIR /build/acnode-dashboard

# copy sources into container build environment
COPY . /build/acnode-dashboard/

RUN git status
RUN git rev-list -1 HEAD > version
RUN cat version
RUN go build
RUN cd bootstrapper && go build
RUN cd logwatcher && go build
RUN cd cleanuptool && go build

# Second container - build frontend
FROM node:16.0-alpine as nodebuilder
RUN apk add --no-cache git

WORKDIR /build/acnode-dashboard/frontend

COPY . /build/acnode-dashboard/

RUN npm install
RUN npm run build && cd dist && tar -cf ../bundle.tar static index.html

# Third container - this one actually runs the code
FROM alpine:latest

WORKDIR /opt/acnode-dashboard

COPY --from=builder /build/acnode-dashboard/acnode-dashboard acnode-dashboard
COPY --from=builder /build/acnode-dashboard/bootstrapper/bootstrapper bootstrapper
COPY --from=builder /build/acnode-dashboard/cleanuptool/cleanuptool cleanuptool
COPY --from=builder /build/acnode-dashboard/version version
COPY --from=nodebuilder /build/acnode-dashboard/frontend/bundle.tar frontend.tar
COPY static static
COPY templates templates

RUN cd static && tar -xf ../frontend.tar && rm -f ../frontend.tar

RUN adduser -S acnodedashboard

USER acnodedashboard

CMD ["./acnode-dashboard"]