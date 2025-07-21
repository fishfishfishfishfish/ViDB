#!/bin/bash
source env.sh

timestamp=$(date +"%Y%m%d_%H%M%S")
./update_benchmark.sh vidb "${timestamp}" > ${logdir}/update_benchmark_vidb_${timestamp}.log 2>&1
# python3 parse_log_range.py ${resdir}/results_vidb range_benchmark_${timestamp}