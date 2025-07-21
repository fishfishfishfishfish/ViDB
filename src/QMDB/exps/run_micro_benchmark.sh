#!/bin/bash
source env.sh

timestamp=$(date +"%Y%m%d_%H%M%S")
./micro_benchmark.sh qmdb "${timestamp}" > ${logdir}/micro_benchmark_qmdb_${timestamp}.log 2>&1
python3 plot_micro_benchmark.py ${resdir}/results_qmdb micro_benchmark_${timestamp}