# syntax=docker/dockerfile:1.7

FROM --platform=$BUILDPLATFORM golang:1.25.5 AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
ARG VERSION=dev
ARG COMMIT=none
ARG DATE=unknown

ARG TARGETOS
ARG TARGETARCH

RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
  go build -trimpath \
  -ldflags="-s -w \
    -X 'main.version=${VERSION}' \
    -X 'main.commit=${COMMIT}' \
    -X 'main.date=${DATE}'" \
  -o /out/runner ./

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/runner /runner
ENTRYPOINT ["/runner"]
