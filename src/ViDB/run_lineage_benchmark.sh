#!/bin/bash
source env.sh

timestamp=$(date +"%Y%m%d_%H%M%S")
./lineage_benchmark.sh vidb "${timestamp}" > ${logdir}/lineage_benchmark_vidb_${timestamp}.log 2>&1
# python3 parse_log_lineage.py ${resdir}/results_vidb lineage_benchmark_${timestamp}