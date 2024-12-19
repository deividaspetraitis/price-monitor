# syntax=docker/dockerfile:1

ARG GO_VERSION="1.23"
ARG RUNNER_IMAGE="ubuntu"

# --------------------------------------------------------
# Builder
# --------------------------------------------------------

FROM golang:1.23-alpine as builder

ARG GIT_VERSION
ARG GIT_COMMIT

WORKDIR /app

COPY go.mod ./
COPY . .

RUN set -eux; apk add --no-cache ca-certificates build-base linux-headers

RUN GOWORK=off go build -mod=readonly \
    -ldflags \
    -v -o /app/build/monitord /app/cmd/monitord/main.go 

# --------------------------------------------------------
# Runner
# --------------------------------------------------------

FROM ${RUNNER_IMAGE}
COPY --from=builder /app/build/monitord /bin/monitord
ENV HOME /app
WORKDIR $HOME
RUN apt-get update && \
    apt-get install curl vim nano -y

# Use JSON array format for ENTRYPOINT
# If array is not used, the command arguments to docker run are ignored.
ENTRYPOINT ["/bin/monitord"]

# Default CMD
CMD []
