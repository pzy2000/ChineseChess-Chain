#!/bin/bash

BROWSER_BIN="chainmaker-browser.bin"

go build -o ${BROWSER_BIN} ../src
CONFIG_PATH="../configs/"
nohup ./${BROWSER_BIN} -config ${CONFIG_PATH} --env dev >output 2>&1 &