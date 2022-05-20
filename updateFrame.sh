#!/bin/bash

#/usr/bin/raspistill -w 1024 -h 768 -n -q 80 -ex auto -hf -vf -t 1 -o /tmp/frame.jpg
/usr/bin/raspistill -w 2588 -h 1920 -n -q 80 -br 60 -co 20 -ex auto -hf -vf -t 1 -o /tmp/frame.jpg