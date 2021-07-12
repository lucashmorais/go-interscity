#!/bin/bash

for i in $(seq 0 7);
do
  if [ $i != 0 ]
    then
      chmod a+w /sys/devices/system/cpu/cpu"$i"/online
  fi
      	
  chmod a+w /sys/devices/system/cpu/cpu"$i"/cpufreq/{scaling_governor,scaling_min_freq,scaling_max_freq,energy_performance_preference}
done
