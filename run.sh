#!/bin/sh
clear && printf '\e[3J' && go build -o .out/main cmd/main.go && .out/main
