#!/bin/bash
source env.sh

timestamp=$(date +"%Y%m%d_%H%M%S")
echo $timestamp

./lineage_benchmark.sh ledgerdb $timestamp > ${logdir}/lineage_benchmark_ledgerdb_${timestamp}.log 2>&1
python3 plot_lineage_benchmark.py ${resdir}/results_ledgerdb lineage_benchmark_${timestamp}

./lineage_benchmark.sh sqlledger $timestamp > ${logdir}/lineage_benchmark_sqlledger_${timestamp}.log 2>&1
python3 plot_lineage_benchmark.py ${resdir}/results_sqlledger lineage_benchmark_${timestamp}

./lineage_benchmark.sh qldb $timestamp > ${logdir}/lineage_benchmark_qldb_${timestamp}.log 2>&1
python3 plot_lineage_benchmark.py ${resdir}/results_qldb lineage_benchmark_${timestamp}


