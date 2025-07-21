#!/bin/bash
source env.sh

timestamp=$(date +"%Y%m%d_%H%M%S")
./range_benchmark.sh qmdb "${timestamp}" > ${logdir}/range_benchmark_qmdb_${timestamp}.log  2>&1
python3 plot_range_benchmark.py ${resdir}/results_qmdb range_benchmark_${timestamp}
