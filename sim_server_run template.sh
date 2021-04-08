#!/bin/sh

cd /home/sim_pc/SIM-go-server

sudo chmod o+rw /dev/ttyACM0

sudo nmcli conn up netplan-<interface name>
sudo systemctl restart isc-dhcp-server.service

#go build -o sim_go_server .

./sim_go_server
