IMAGE=ezenki/nats-auth-redirector
VERSION=${TAG:-latest}
COMMIT=$(git rev-parse --short HEAD)
DATE=$(date -u +%Y-%m-%dT%H:%M:%SZ)

docker buildx build \
  --platform linux/amd64,linux/arm64 \
  --build-arg VERSION=$VERSION \
  --build-arg COMMIT=$COMMIT \
  --build-arg DATE=$DATE \
  -t $IMAGE:$VERSION \
  --push \
  .