#!/bin/bash
set -e

IMAGE_NAME="go-secretspec-protonpass"
CONTAINER_NAME="go-protonpass-app"
PORT="8080"

echo "=== Ensuring pass-cli is authenticated ==="
if ! pass-cli status > /dev/null 2>&1; then
    if [ -n "$PROTON_PASS_PAT" ]; then
        echo "Logging in via PROTON_PASS_PAT..."
        pass-cli login --pat "$PROTON_PASS_PAT"
    else
        echo "Notice: pass-cli not authenticated. Please run 'pass-cli login' or set \$PROTON_PASS_PAT."
    fi
fi

echo "=== Building Podman image ==="
podman build -t $IMAGE_NAME -f Containerfile .

echo "=== Stopping existing container if running ==="
podman stop $CONTAINER_NAME 2>/dev/null || true
podman rm $CONTAINER_NAME 2>/dev/null || true

echo "=== Deploying container with SecretSpec + Proton Pass ==="
podman run -d \
  --name $CONTAINER_NAME \
  -p $PORT:8080 \
  -e PROTON_PASS_PAT="$PROTON_PASS_PAT" \
  --restart always \
  $IMAGE_NAME

echo "=== Deployment complete! ==="
echo "Application running at http://localhost:$PORT"
podman ps -f name=$CONTAINER_NAME
