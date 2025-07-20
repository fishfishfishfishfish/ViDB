#!/bin/bash
source env.sh

db_name=$1
test_name=$2
echo "db_name: $db_name, test_name=$test_name"

# define experiment parameters
load_account=(1000000 10000000)
load_batch_size=10000
value_sizes=(1024)
key_size=32
query_versions="1,2,4,10,20,40"
get_proof=""
num_transaction_account=20
num_transaction_version=40

# init dirs
data_path="${datadir}"
result_path="${resdir}/results_${db_name}/lineage_benchmark_${test_name}"
echo "data_path: $data_path"
echo "result_path: $result_path"


mkdir -p $data_path
rm -rf ${data_path}/*
mkdir -p ${result_path}
rm -rf ${result_path}/*


for n_acc in "${load_account[@]}"; do
    for value_size in "${value_sizes[@]}"; do
        set -x
        # clean data dir
        rm -rf $data_path/*
        result_file="${result_path}/e${n_acc}v${value_size}.csv"

        # run
        ${builddir}/build_release_${db_name}/bin/lineageBenchmark -a $n_acc -b $load_batch_size -t $num_transaction_version -z $num_transaction_account $get_proof -l $query_versions -k $key_size -v $value_size -d $data_path -r $result_file
        
        sleep 5
        set +x
    done

done