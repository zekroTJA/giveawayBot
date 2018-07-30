#!/bin/bash

if [ -z $1 ]
then
    echo "Please enter the version as argument you want to install!"
    exit -1
fi

VERSION=$1

if ! [ $(command -v curl) ]
then
    echo "'curl' installation is required for this installer."
    exit -1
fi

mkdir -p giveawayBot/lang
cd giveawayBot

curl https://github.com/zekroTJA/giveawayBot/releases/download/$VERSION/build >> build
curl https://raw.githubusercontent.com/zekroTJA/giveawayBot/master/config_example.yaml >> config.yaml
curl https://raw.githubusercontent.com/zekroTJA/giveawayBot/master/lang/en-US.yaml >> lang/en-US.yaml
curl https://raw.githubusercontent.com/zekroTJA/giveawayBot/master/lang/de-DE.yaml >> lang/de-DE.yaml

chmod +x build

echo "Files downloaded and installed."
echo "Open 'config.yaml', enter your preferences and start the bot with the 'build' binary."