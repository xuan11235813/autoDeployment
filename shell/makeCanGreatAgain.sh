#! /bin/bash
sudo pkill networkchecker
sudo bash -c "echo performance > /sys/devices/system/cpu/cpu0/cpufreq/scaling_governor"
sudo bash -c "echo performance > /sys/devices/system/cpu/cpu1/cpufreq/scaling_governor"
sudo bash -c "echo performance > /sys/devices/system/cpu/cpu2/cpufreq/scaling_governor"
sudo bash -c "echo performance > /sys/devices/system/cpu/cpu3/cpufreq/scaling_governor"
sudo ip link set can0 type can bitrate 500000
sudo ip link set can1 type can bitrate 500000
sudo ip link set can2 type can bitrate 500000
sudo ip link set can3 type can bitrate 500000
sudo ip link set can0 up
sudo ip link set can1 up
sudo ip link set can2 up
sudo ip link set can3 up

