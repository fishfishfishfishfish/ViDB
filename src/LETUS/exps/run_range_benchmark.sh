#!/bin/bash
source env.sh

timestamp=$(date +"%Y%m%d_%H%M%S")
./range_benchmark.sh letus "${timestamp}" > ${logdir}/range_benchmark_letus_${timestamp}.log  2>&1
python3 plot_range_benchmark.py ${resdir}/results_letus range_benchmark_${timestamp}
