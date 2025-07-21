#!/bin/bash
source env.sh

timestamp=$(date +"%Y%m%d_%H%M%S")
./tier_benchmark.sh vidb "${timestamp}" > ${logdir}/tier_benchmark_vidb_${timestamp}.log 2>&1
# python3 parse_log_tier.py ${resdir}/results_vidb tier_benchmark_${timestamp}