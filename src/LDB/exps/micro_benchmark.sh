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
load_batch_size=5000
key_size=32

# init dirs
data_path="${datadir}"
result_path="${resdir}/results_${db_name}/micro_benchmark_${test_name}"
echo "data_path: $data_path"
echo "result_path: $result_path"


mkdir -p $data_path
rm -rf ${data_path}/*
mkdir -p ${result_path}
rm -rf ${result_path}/*
echo "entry_count, batch_size, value_size, key_size, folder_size" > "${result_path}/size"    

# 运行测试
for n_acc in "${load_account[@]}"; do
    for batch_size in "${batch_sizes[@]}"; do
        for value_size in "${value_sizes[@]}"; do
            set -x
            # clean data
            rm -rf ${data_path}/*            
            result_file="${result_path}/e${n_acc}b${batch_size}v${value_size}.csv"

            # run
            ${builddir}/build_release_${db_name}/bin/microBenchmark -a $n_acc -b $load_batch_size -t $num_transaction_version -z $batch_size -k $key_size -v $value_size -d $data_path -r $result_file

            # get file sizes
            if [ ! -d "$data_path" ]; then
                echo "$data_path do not exist."
                exit 1
            fi
            folder_size=$(du -sk "$data_path" | cut -f1)
            echo "${n_acc}, ${batch_size}, ${value_size}, ${key_size}, ${folder_size}" >> "${result_path}/size"  

            sleep 5
            set +x
        done
    done
done