#!/bin/bash

echo start downloading... $1

version=$1
mkdir helm-$version
cd helm-$version && \
    wget https://get.helm.sh/helm-$version-darwin-amd64.tar.gz && \
    wget https://get.helm.sh/helm-$version-windows-amd64.zip && \
    wget https://get.helm.sh/helm-$version-linux-386.tar.gz && \
    wget https://get.helm.sh/helm-$version-linux-amd64.tar.gz && \
    wget https://get.helm.sh/helm-$version-linux-arm.tar.gz && \
    wget https://get.helm.sh/helm-$version-linux-arm64.tar.gz

echo $1 download success...
