#!/bin/bash

sudo ip link set br0 down
sudo brctl delbr br0
sudo ip link del veth1
sudo ip link del veth3
