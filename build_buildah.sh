#!/bin/bash 
# Set the required variables
DOCKERFILE="Dockerfile"
REGISTRY="docker.io"
USER="shelltux"
IMAGE_NAME="knx_data_exposer"
IMAGE_TAG="latest"
# Set your manifest name
MANIFEST_NAME="multiarch-knx_data_exposer"

# Create a multi-architecture manifest
buildah manifest create ${MANIFEST_NAME}

for i in amd64 arm64 arm; do
    buildah bud \
        --tag "${REGISTRY}/${USER}/${IMAGE_NAME}:${IMAGE_TAG}" \
        --manifest ${MANIFEST_NAME} \
        --platform linux/${i} \
        .
done

# Push the full manifest, with all CPU Architectures
buildah manifest push ${MANIFEST_NAME} "docker://${REGISTRY}/${USER}/${IMAGE_NAME}:${IMAGE_TAG}" --all
