#!/bin/sh

cd /home/sim_pc/Sim_Reader_go_server/go-server

sudo chmod o+rw /dev/ttyACM0

sudo nmcli conn up netplan-enx00e04c364564
sudo systemctl restart isc-dhcp-server.service

#go build -o main .

./main