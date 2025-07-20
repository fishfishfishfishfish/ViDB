#!/bin/bash
source env.sh

db_name=$1
test_name=$2
echo "db_name: $db_name, test_name=$test_name"

# define experiment parameters
load_account=(1000000)
batch_sizes=(500 1000 2000 3000 4000 5000)
value_sizes=(256 512 1024 2048)
num_transaction_version=20
load_batch_size=10000
key_size=64

data_path="${datadir}/data/"
index_path="${datadir}/index"
result_path="${resdir}/results_${db_name}/micro_benchmark_${test_name}"
echo "data_path: $data_path"
echo "index_path: $index_path"
echo "result_path: $result_path"

mkdir -p ${data_path}
mkdir -p ${index_path}
mkdir -p ${result_path}
rm -rf ${data_path}/*
rm -rf ${index_path}/*
rm -rf ${result_path}/*

for n_acc in "${load_account[@]}"; do
    for batch_size in "${batch_sizes[@]}"; do
        for value_size in "${value_sizes[@]}"; do
            set -x
            # clean
            rm -rf $data_path/*
            rm -rf $index_path/*            
            result_file="${result_path}/e${n_acc}b${batch_size}v${value_size}.csv"

            # 运行测试并提取结果
            ${builddir}/build_release_letus/bin/microBenchmark -a $n_acc -b $load_batch_size -t $num_transaction_version -z $batch_size -k $key_size -v $value_size -d $data_path -i $index_path -r $result_file
            sleep 5
            set +x

            # get file size
            if [ ! -d "$data_path" ]; then
                echo "$data_path do not exist."
                exit 1
            fi
            if [ ! -d "$index_path" ]; then
                echo "$index_path do not exist."
                exit 1
            fi
            data_folder_size=$(du -sk "$data_path" | cut -f1)
            index_folder_size=$(du -sk "$index_path" | cut -f1)
            echo "${n_acc},${batch_size},${value_size},${key_size},${data_folder_size},${index_folder_size},$(${data_folder_size}+${index_folder_size})" >> "${result_dir}/size"  
        done
    done
done