# syntax=docker/dockerfile:1

FROM golang:1.19

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY *.go ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-gs-ping

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/engine/reference/builder/#expose
EXPOSE 80

# Run
CMD ["/docker-gs-ping"]

# # syntax=docker/dockerfile:1.4
# FROM --platform=$BUILDPLATFORM golang:1.18-alpine AS builder

# WORKDIR /code

# ENV CGO_ENABLED 0
# ENV GOPATH /go
# ENV GOCACHE /go-build

# COPY go.mod go.sum ./
# RUN --mount=type=cache,target=/go/pkg/mod/cache \
#     go mod download

# COPY . .

# RUN --mount=type=cache,target=/go/pkg/mod/cache \
#     --mount=type=cache,target=/go-build \
#     go build -o bin/backend main.go

# CMD ["/code/bin/backend"]

# FROM builder as dev-envs

# RUN <<EOF
# apk update
# apk add git
# EOF

# RUN <<EOF
# addgroup -S docker
# adduser -S --shell /bin/bash --ingroup docker vscode
# EOF

# # install Docker tools (cli, buildx, compose)
# COPY --from=gloursdocker/docker / /

# CMD ["go", "run", "main.go"]

# FROM scratch
# COPY --from=builder /code/bin/backend /usr/local/bin/backend
# CMD ["/usr/local/bin/backend"]
