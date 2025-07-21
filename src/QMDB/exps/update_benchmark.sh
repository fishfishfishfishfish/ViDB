#!/bin/bash
source env.sh

db_name=${1:-qmdb}
test_name=${2:-test}
echo "db_name: $db_name, test_name=$test_name"

# define experiment parameters
load_account=(500 1000 2000 3000 4000 5000)
value_sizes=(1024)
key_size=32
update_count=100

randsrc_file="${builddir}/build_release_${db_name}/randsrc.dat"
data_path="${datadir}"
result_path="${resdir}/results_${db_name}/update_benchmark_${test_name}"

mkdir -p ${data_path}
mkdir -p ${result_path}
rm -rf ${data_path}/*
rm -rf ${result_path}/*

for n_acc in "${load_account[@]}"; do
    for value_size in "${value_sizes[@]}"; do
        set -x
        # clean data
        rm -rf ${data_path}/*        
        result_file="${result_path}/e${n_acc}u${update_count}v${value_size}.csv"

        # run
        ${builddir}/build_release_${db_name}/release/update_benchmark --randsrc-filename ${randsrc_file} --db-dir ${data_path} --entry-count $n_acc --tps-blocks $update_count --key-size $key_size --val-size $value_size --output-filename $result_file
        sleep 5
        set +x
    done
done