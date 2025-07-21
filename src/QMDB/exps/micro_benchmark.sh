#!/bin/bash
source env.sh

db_name=${1:-qmdb}
test_name=${2:-test}
echo "db_name: $db_name, test_name=$test_name"

# define experiment parameters
entries_counts=(1000000)
batch_sizes=(500 1000 2000 4000 5000)
value_sizes=(256 512 1024 2048)
tps_blocks=20
key_size=32

randsrc_file="${builddir}/build_release_${db_name}/randsrc.dat"
data_path="${datadir}"
result_path="${resdir}/results_${db_name}/micro_benchmark_${test_name}"

mkdir -p ${data_path}
mkdir -p ${result_path}
rm -rf ${data_path}/*
rm -rf ${result_path}/*

for n_acc in "${entries_counts[@]}"; do
    for ops_per_block in "${batch_sizes[@]}"; do
        for value_size in "${value_sizes[@]}"; do
            set -x
            # clean data
            rm -rf ${data_path}/*
            result_file="${result_path}/e${n_acc}b${ops_per_block}v${value_size}.csv"

            # run
            ${builddir}/build_release_${db_name}/release/micro_benchmark --randsrc-filename ${randsrc_file} --db-dir ${data_path} --tps-blocks ${tps_blocks} --entry-count ${n_acc} --ops-per-block ${ops_per_block} --key-size ${key_size} --val-size ${value_size} --output-filename ${result_file}
            
            # get file sizes
            if [ ! -d "$data_path" ]; then
                echo "$data_path do not exist."
                exit 1
            fi
            folder_size=$(du -sk "$data_path" | cut -f1)
            echo "${n_acc}, ${ops_per_block}, ${value_size}, ${key_size}, ${folder_size}" >> "${result_path}/size"  
            sleep 5
            set +x
        done
    done
done