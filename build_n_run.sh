#!/bin/sh

go build -o main .

sudo chmod o+rw /dev/ttyACM0

./main >> logfile