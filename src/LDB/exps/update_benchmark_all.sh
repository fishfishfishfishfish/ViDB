#!/bin/bash
source env.sh

timestamp=$(date +"%Y%m%d_%H%M%S")
echo $timestamp

./update_benchmark.sh qldb $timestamp > ${logdir}/update_benchmark_qldb_${timestamp}.log 2>&1
python plot_update_benchmark.py ${resdir}/results_qldb update_benchmark_qldb_${timestamp}
./update_benchmark.sh sqlledger $timestamp > ${logdir}/update_benchmark_sqlledger_${timestamp}.log 2>&1
python plot_update_benchmark.py ${resdir}/results_sqlledger update_benchmark_sqlledger_${timestamp}
./update_benchmark.sh ledgerdb $timestamp > ${logdir}/update_benchmark_ledgerdb_${timestamp}.log 2>&1
python plot_update_benchmark.py ${resdir}/results_ledgerdb update_benchmark_ledgerdb_${timestamp}
