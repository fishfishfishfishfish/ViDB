#!/bin/bash
source env.sh

timestamp=$(date +"%Y%m%d_%H%M%S")
./rollback_benchmark.sh vidb "${timestamp}" > ${logdir}/rollback_benchmark_vidb_${timestamp}.log 2>&1
# python3 parse_log_rollback.py ${resdir}/results_vidb rollback_benchmark_${timestamp}