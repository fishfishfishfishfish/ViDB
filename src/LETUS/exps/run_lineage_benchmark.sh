#!/bin/bash
source env.sh

timestamp=$(date +"%Y%m%d_%H%M%S")
./lineage_benchmark.sh letus "${timestamp}" > ${logdir}/lineage_benchmark_letus_${timestamp}.log  2>&1
python3 plot_lineage_benchmark.py ${resdir}/results_letus lineage_benchmark_${timestamp}