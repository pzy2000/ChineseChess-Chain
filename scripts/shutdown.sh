#!/bin/bash

BROWSER_BIN="chainmaker-browser.bin"
pid=`ps -ef | grep ${BROWSER_BIN} | grep -v grep | awk '{print $2}'`
if [ ! -z ${pid} ];then
    kill -9 $pid
fi