#!/bin/bash
source env.sh

timestamp=$(date +"%Y%m%d_%H%M%S")
# timestamp="test"
echo $timestamp

./micro_benchmark.sh qldb $timestamp > ${logdir}/micro_benchmark_qldb_${timestamp}.log 2>&1
python3 plot_micro_benchmark.py ${resdir}/results_qldb micro_benchmark_${timestamp}
./micro_benchmark.sh sqlledger $timestamp > ${logdir}/micro_benchmark_sqlledger_${timestamp}.log 2>&1
python3 plot_micro_benchmark.py ${resdir}/results_sqlledger micro_benchmark_${timestamp}
./micro_benchmark.sh ledgerdb $timestamp > ${logdir}/micro_benchmark_ledgerdb_${timestamp}.log 2>&1
python3 plot_micro_benchmark.py ${resdir}/results_ledgerdb micro_benchmark_${timestamp}
