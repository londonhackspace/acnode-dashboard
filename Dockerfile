# First container - does the build
FROM golang:1.16 AS builder

WORKDIR /build/acnode-dashboard

# copy sources into container build environment
COPY auth .
COPY *.go .
COPY go.mod .
COPY .git .

RUN go build
RUN 

# Second container - this one actually runs the code
FROM alpine:latest

WORKDIR /opt/acnode-dashboard

COPY --from=builder /build/acnode-dashboard/acnode-dashboard acnode-dashboard
COPY static .
COPY templates .

RUN useradd acnodedashboard

USER acnodedashboard

CMD ["./acnode-dashboard"]