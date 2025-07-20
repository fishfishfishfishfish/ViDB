#!/bin/bash
source env.sh

timestamp=$(date +"%Y%m%d_%H%M%S")
echo $timestamp

./range_benchmark.sh ledgerdb $timestamp > ${logdir}/range_benchmark_ledgerdb_${timestamp}.log 2>&1
python3 plot_range_benchmark.py ${resdir}/results_ledgerdb range_benchmark_${timestamp}
./range_benchmark.sh sqlledger $timestamp > ${logdir}/range_benchmark_sqlledger_${timestamp}.log 2>&1
python3 plot_range_benchmark.py ${resdir}/results_sqlledger range_benchmark_${timestamp}
./range_benchmark.sh qldb $timestamp > ${logdir}/range_benchmark_qldb_${timestamp}.log 2>&1
python3 plot_range_benchmark.py ${resdir}/results_qldb range_benchmark_${timestamp}
