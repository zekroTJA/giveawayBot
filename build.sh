#!/bin/bash


#!/bin/bash

TAG=$(git describe --tags)

if [ ! -d builds ]; then
    mkdir builds
fi

BUILDS=( \
    'linux;arm' \
    'linux;amd64' \
    'windows;amd64' \
    'darwin;amd64' \
)

for BUILD in ${BUILDS[*]}; do

    IFS=';' read -ra SPLIT <<< "$BUILD"
    OS=${SPLIT[0]}
    ARCH=${SPLIT[1]}

    echo "Building ${OS}_$ARCH..."
    (env GOOS=$OS GOARCH=$ARCH \
        go build -o builds/giveawayBot_${OS}_$ARCH \
        -ldflags "-X main.appVersion=$TAG")

    if [ "$OS" = "windows" ]; then
        mv builds/giveawayBot_windows_$ARCH builds/giveawayBot_windows_${ARCH}.exe
    fi

done


wait