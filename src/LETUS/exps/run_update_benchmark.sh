#!/bin/bash
source env.sh

timestamp=$(date +"%Y%m%d_%H%M%S")
./update_benchmark.sh letus "${timestamp}" > ${logdir}/update_benchmark_letus_${timestamp}.log  2>&1
python3 plot_update_benchmark.py ${resdir}/results_letus update_benchmark_${timestamp}