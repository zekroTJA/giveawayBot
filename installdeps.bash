#!/bin/bash

read -p "Have you set your \$GOPATH to current work directory? [yN] " yn

if [ "$yn" = "y" ] || [ "$yn" = "Y" ]
then
    go get github.com/bwmarrin/discordgo
    go get github.com/go-yaml/yaml
    exit 0
fi

echo Canceled.