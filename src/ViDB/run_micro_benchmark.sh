#!/bin/bash
source env.sh

timestamp=$(date +"%Y%m%d_%H%M%S")
./micro_benchmark.sh vidb "${timestamp}" > ${logdir}/micro_benchmark_vidb_${timestamp}.log 2>&1
python3 parse_log_micro.py ${resdir}/results_vidb micro_benchmark_${timestamp}