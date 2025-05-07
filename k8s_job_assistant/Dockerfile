# --- Go Builder Stage ---
FROM golang:1.24-alpine AS backend-builder
ENV CGO_ENABLED=0
WORKDIR /backend

COPY goapp/go.* .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

COPY goapp/. .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -trimpath -ldflags="-s -w" -o bin/service

# --- React Builder Stage ---
FROM --platform=$BUILDPLATFORM node:21.6-alpine3.18 AS frontend-builder
WORKDIR /ui

COPY reactapp/package*.json ./
RUN --mount=type=cache,target=/usr/src/app/.npm \
    npm set cache /usr/src/app/.npm && \
    npm ci

COPY reactapp ./
RUN npm run build

# --- Final Minimal Distroless Stage ---
FROM gcr.io/distroless/static:nonroot
ARG VERSION
ENV VERSION=$VERSION

LABEL org.opencontainers.image.title="K8sJobAssistant" \
    org.opencontainers.image.description="A tool to Run K8s Job from a simple UI"

WORKDIR /
COPY --from=backend-builder /backend/bin/service /service
COPY --from=frontend-builder /ui/build /ui

USER nonroot
ENTRYPOINT ["/service"]
