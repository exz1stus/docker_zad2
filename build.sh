DOCKER_USER="mbmbpididi"
IMAGE_NAME="weather-app"

docker buildx build \
  --platform linux/amd64,linux/arm64 \
  -t $DOCKER_USER/$IMAGE_NAME:latest \
  --cache-from type=registry,ref=$DOCKER_USER/$IMAGE_NAME:cache \
  --cache-to type=inline \
  --push .