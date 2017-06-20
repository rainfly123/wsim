#!/bin/sh
while [ 1 ]
do
wsim=$(ps -ef | grep wsim | grep -v "grep")
if [ -z "$wsim" ];then
    nohup ./wsim &
fi
sleep 3
done
