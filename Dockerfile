# syntax=docker/dockerfile:1.7-labs

#[=======================================================================[
# Description : Docker image containing the
#]=======================================================================]

# Docker image repository to use. Use `` for public images.
ARG IMAGE_BASE_REGISTRY

#### ---- Build ---- ####
FROM ${IMAGE_BASE_REGISTRY}golang:1.22.2-alpine3.19 as devenv

LABEL maintainer=arash.idelchi

USER root

# Basic good practices
SHELL ["/bin/ash", "-o", "pipefail", "-c"]

# Utilities to debug & run tests
RUN apk add --no-cache \
    curl==8.5.0-r0 \
    tzdata==2024a-r0 \
    build-base==0.5-r3

# Create User (Alpine)
ARG USER=user
RUN addgroup -S -g 1001 ${USER} && \
    adduser -S -u 1001 -G ${USER} -h /home/${USER} -s /bin/ash ${USER}

USER ${USER}
WORKDIR /home/${USER}

WORKDIR /tmp/go

COPY --chown=${USER}:${USER} go.mod go.sum /tmp/go/
RUN go mod download

COPY --parents --chown=${USER}:${USER} **/*.go /tmp/go/
ARG GO_WSLINT_VERSION="unofficial & built by unknown"
ARG GO_WSLINT_COMPILE_OPTIONS=""
ARG GOCACHE=/home/${USER}/.cache/go-build
RUN --mount=type=cache,target=${GOCACHE},uid=1001,gid=1001 \
    CGO_ENABLED=1 go install ${GO_WSLINT_COMPILE_OPTIONS} -tags musl -ldflags="-extldflags -static -s -w -X 'main.version=${GO_WSLINT_VERSION}'" ./... && \
    rm -rf /tmp/go/

COPY --link --chown=${USER}:${USER} utils/ /home/${USER}/utils/

# Clear the base image entrypoint
ENTRYPOINT [""]
CMD ["/bin/ash"]

# Timezone
ENV TZ=Europe/Zurich

# Add some convenient aliases
ENV ENV="/home/${USER}/.ashrc"
COPY --link --chown=${USER}:${USER} .ashrc /home/${USER}/.ashrc

USER ${USER}
WORKDIR /home/${USER}

#### ---- wslint ---- ####
FROM scratch as wslint

LABEL maintainer=arash.idelchi

ARG USER=user

# Copy artifacts from the devenv stage
COPY --from=devenv /etc/passwd /etc/passwd
COPY --from=devenv /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=devenv /go/bin/wslint /wslint
COPY --from=devenv --chown=${USER}:${USER} /home/user/utils /home/${USER}/utils

USER ${USER}
WORKDIR /home/${USER}

# Clear the base image entrypoint
ENTRYPOINT ["/wslint"]
CMD [""]

# Timezone
ENV TZ=Europe/Zurich
