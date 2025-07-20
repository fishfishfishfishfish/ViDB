#!/bin/bash
source env.sh

db_name=$1
test_name=$2
echo "db_name: $db_name, test_name=$test_name"

# define experiment parameters
load_account=(500 1000 2000 3000 4000 5000)
value_sizes=(1024)
key_size=64
update_count=100

data_path="${datadir}/data/"
index_path="${datadir}/index"
result_path="${resdir}/results_${db_name}/update_benchmark_${test_name}"
echo "data_path: $data_path"
echo "index_path: $index_path"
echo "result_path: $result_path"

mkdir -p ${data_path}
mkdir -p ${index_path}
mkdir -p ${result_path}
rm -rf ${data_path}/*
rm -rf ${index_path}/*
rm -rf ${result_path}/*

# 运行测试
for n_acc in "${load_account[@]}"; do
    for value_size in "${value_sizes[@]}"; do
        set -x
        # clean
        rm -rf $data_path/*
        rm -rf $index_path/*
        
        result_file="${result_path}/e${n_acc}u${update_count}v${value_size}.csv"
        
        # run
        ${builddir}/build_release_${db_name}/bin/updateBenchmark -a $n_acc -t $update_count -k $key_size -v $value_size -d $data_path -i $index_path -r $result_file
        sleep 5
        set +x
    done
done