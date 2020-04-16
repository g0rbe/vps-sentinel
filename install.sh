#!/bin/bash
set -e

install() {

    if [ $EUID != 0 ]
    then
        echo "Run as root!"
        exit 1
    fi

    cp ./bin/main /usr/bin/vps-sentinel
    chown root:root /usr/bin/vps-sentinel
    chmod 0500 /usr/bin/vps-sentinel

    
    cp ./bin/vps-sentinel.conf /etc/vps-sentinel.conf
    chmod 0600 /etc/vps-sentinel.conf

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

update() {
    git pull

    cp ./bin/main /usr/bin/vps-sentinel

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

    if [ -e /etc/vps-sentinel.conf ]
    then
        nano /etc/vps-sentinel.conf
    else
        nano ./bin/vps-sentinel.conf
    fi

    if [ -e /etc/systemd/system/vps-sentinel.service ]
    then
        nano /etc/systemd/system/vps-sentinel.service
    else
        nano ./bin/vps-sentinel.service
    fi

    if [ -e /etc/systemd/system/vps-sentinel.timer ]
    then
        nano /etc/systemd/system/vps-sentinel.timer
    else
        nano ./bin/vps-sentinel.timer
    fi

    systemctl daemon-reload
}

case "$1" in
"install")
    install
    ;;
"remove")
    remove
    ;;
"update")
    update
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
