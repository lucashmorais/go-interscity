#!/bin/bash

for i in $(seq 0 7);
do
  if [ $i != 0 ]
	  then
		  echo 1 > /sys/devices/system/cpu/cpu"$i"/online
  fi	

  echo performance > /sys/devices/system/cpu/cpu"$i"/cpufreq/scaling_governor
  echo performance > /sys/devices/system/cpu/cpu"$i"/cpufreq/energy_performance_preference
  echo "$1000" > /sys/devices/system/cpu/cpu"$i"/cpufreq/scaling_max_freq
  echo "$1000" > /sys/devices/system/cpu/cpu"$i"/cpufreq/scaling_min_freq
done
