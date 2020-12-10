#!/bin/sh

cd /home/sim_pc/Sim_Reader_go_server/go-server

sudo chmod o+rw /dev/ttyACM0

go build -o main .

./main