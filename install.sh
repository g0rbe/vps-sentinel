#!/bin/bash

install() {

    if [ $EUID != 0 ]
    then
        echo "Run as root!"
        exit 1
    fi

    cp ./bin/main /usr/bin/vps-sentinel
    chown root:root /usr/bin/vps-sentinel
    chmod 0500 /usr/bin/vps-sentinel

    # Dont overwrite existing config
    if [ ! -e /etc/vps-sentinel.conf ]
    then
        cp ./bin/vps-sentinel.conf /etc/vps-sentinel.conf
        chmod 0600 /etc/vps-sentinel.conf
    fi

    cp ./bin/vps-sentinel.service /etc/systemd/system
    cp ./bin/vps-sentinel.timer /etc/systemd/system

    systemctl daemon-reload
    systemctl enable --now vps-sentinel.timer
}

remove() {

    if [ $EUID != 0 ]
    then
        echo "Run as root!"
        exit 1
    fi

    rm /usr/bin/vps-sentinel
    rm /etc/vps-sentinel.conf

    systemctl disable --now vps-sentinel.timer

    rm /etc/systemd/system/vps-sentinel.*
    systemctl daemon-reload
}

build() {
    go build main.go
    mv ./main ./bin
}

conf() {

    if [ $EUID != 0 ]
    then
        echo "Run as root!"
        exit 1
    fi

    nano /etc/vps-sentinel.conf

    nano /etc/systemd/system/vps-sentinel.*

    systemctl daemon-reload
}

case "$1" in
"install")
    install
    ;;
"remove")
    remove
    ;;
"build")
    build
    ;;
"conf")
    conf
    ;;
*)
    echo "Invalid command: $1"
    exit 1
esac
