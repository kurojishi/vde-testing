#!/bin/bash
tunctl -u kurojishi -t tap0
ip link set tap0 up
ip addr add 192.168.4.1/24 dev tap0
