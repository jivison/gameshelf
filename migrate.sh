#!/bin/zsh

if [[ "$1" == "reverse" ]]; then
    cd ./db; rambler -c "./rambler.hjson" reverse -all
else
    cd ./db; rambler -c "./rambler.hjson" apply -all
fi
