#!/bin/bash
source env.sh
db_name=${1:-vidb}
test_name=${2:-test}
echo "db_name: $db_name, test_name=$test_name"

# detect platform & arch
os=$(uname -s | tr '[:upper:]' '[:lower:]')
arch=$(uname -m)

case "$arch" in
    x86_64) arch="amd64" ;;
    aarch64) arch="arm64" ;;
    arm64) arch="arm64" ;;
esac

letus_vidb_binary="${builddir}/build_release_${db_name}/vidb"

echo "Using binary: $letus_vidb_binary"

# define experiment parameters
treeCapacities=("1M" "500K" "300K" "250K" "200K" "100K" "50K") # n1
operationCounts=("1M") # n2
batchSizes=(5000)      # n3
valueSizes=(1024)
keySize=32
bloomCap=5000
cacheCost=0

# directory
dataPath="${datadir}"
outputDir="${resdir}/results_${db_name}/point_brk_test_${test_name}"
mkdir -p ${dataPath}
mkdir -p ${outputDir}
rm -rf ${dataPath}/*
rm -rf ${outputDir}/*
timestamp=$(date +%s)

total_start_time=$(date +%s)

# convert human-readable sizes to numbers
human_to_number() {
    local human_size=$1
    local num=${human_size%[KMG]}
    local unit=${human_size##*[0-9]}
    case $unit in
        K) echo $((num * 1000)) ;;
        M) echo $((num * 1000000)) ;;
        G) echo $((num * 1000000000)) ;;
        *) echo $num ;;
    esac
}

# run with error checking
run_benchmark() {
    local cmd=$1
    local logFile=$2
    local description=$3
    echo "=== Starting $description ===" | tee -a "$logFile"
    echo "Command: $cmd" | tee -a "$logFile"
    echo "Timestamp: $(date)" | tee -a "$logFile"

    if eval "$cmd" >> "$logFile" 2>&1; then
        echo "=== Completed successfully ===" | tee -a "$logFile"
    else
        echo "=== FAILED with exit code $? ===" | tee -a "$logFile"
        exit 1
    fi
    echo -e "\n" | tee -a "$logFile"
}

# main loop
for human_op in "${operationCounts[@]}"; do
    operationCount=$(human_to_number "$human_op")
    for batchSize in "${batchSizes[@]}"; do
        for valueSize in "${valueSizes[@]}"; do
            for human_treeCap in "${treeCapacities[@]}"; do
                rm -rf ${dataPath}
                mkdir -p ${dataPath}
                treeCap=$(human_to_number "$human_treeCap")
                Log="${outputDir}/point_brk_${human_op}_${human_treeCap}_${valueSize}_${timestamp}.log"
                Cmd="${letus_vidb_binary} LatencyBreakdownSingleQuery \
                    --n_1=$treeCap \
                    --n_2=$operationCount \
                    --cacheCost=$cacheCost \
                    --valueSize=$valueSize \
                    --keySize=$keySize \
                    --dataPath=$dataPath"
                run_benchmark "$Cmd" "$Log" "zipf test (ops=${human_op}, batch=${batchSize})"
                sleep 5
            done
        done
    done
done

# for human_op in "${operationCounts[@]}"; do
#     operationCount=$(human_to_number "$human_op")
#     for batchSize in "${batchSizes[@]}"; do
#         for valueSize in "${valueSizes[@]}"; do
#             for treeh in "${tree_height[@]}"; do
#                 rm -rf ${dataPath}
#                 mkdir -p ${dataPath}
#                 Log="${outputDir}/point_brk_${human_op}_${human_treeCap}_${valueSize}_${timestamp}.log"
#                 Cmd="${letus_vidb_binary} LatencyBreakdownSingleQuery \
#                     --treeHeight=$treeh \
#                     --n_2=$operationCount \
#                     --cacheCost=$cacheCost \
#                     --bloomCap=$bloomCap \
#                     --valueSize=$valueSize \
#                     --keySize=$keySize \
#                     --dataPath=$dataPath"
#                 run_benchmark "$Cmd" "$Log" "point breakdown test (ops=${human_op}, batch=${batchSize}, treeHeight=${treeh}, bloomCap=${bloomCap})"
#                 sleep 5
#             done
#         done
#     done
# done

total_end_time=$(date +%s)
total_execution_time=$((total_end_time - total_start_time))
echo "All benchmarks completed successfully in $total_execution_time seconds. Results saved in $outputDir/"