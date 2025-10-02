#!/bin/bash
source env.sh
db_name=${1:-vidb}
test_name=${2:-test}
echo "db_name: $db_name, test_name=$test_name"

# define experiment parameters
# small test
treeCapacities=("500K" "200K" "100K" "50K") # n1
operationCounts=("1M") # Human-readable operation counts, n2
rangeSizes=(5000) # n3, 查询的key的数量
valueSizes=(1024)          
keySize=32

# large test
# treeCapacities=("50M" "20M" "10M" "5M") # n1
# operationCounts=("100M") # Human-readable operation counts, n2
# rangeSizes=(5000) # n3, 查询的key的数量
# valueSizes=(1024)          
# keySize=32

# directory
dataPath="${datadir}"
outputDir="${resdir}/results_${db_name}/range_brk_test_${test_name}"
mkdir -p ${dataPath}
mkdir -p ${outputDir}
rm -rf ${dataPath}/*
rm -rf ${outputDir}/*
# Timestamp for unique filenames
timestamp=$(test_name)

# Start time for total execution
total_start_time=$(date +%s)

# Function to convert human-readable sizes to numbers
human_to_number() {
    local human_size=$1
    local num=${human_size%[KMG]}
    local unit=${human_size##*[0-9]}

    case $unit in
        K) echo $((num * 1000)) ;;
        M) echo $((num * 1000000)) ;;
        G) echo $((num * 1000000000)) ;;
        *) echo $num ;;  # No unit, assume it's already a number
    esac
}

# Function to run command with error checking and logging
run_benchmark() {
    local cmd=$1
    local logFile=$2
    local description=$3

    echo "=== Starting $description ===" | tee -a "$logFile"
    echo "Command: $cmd" | tee -a "$logFile"
    echo "Timestamp: $(date)" | tee -a "$logFile"

    # Run the command and capture output
    if $cmd >> "$logFile" 2>&1; then
        echo "=== Completed successfully ===" | tee -a "$logFile"
    else
        echo "=== FAILED with exit code $? ===" | tee -a "$logFile"
        exit 1
    fi

    echo -e "\n" | tee -a "$logFile"
}

# Main benchmark execution
for human_op in "${operationCounts[@]}"; do
    operationCount=$(human_to_number "$human_op")
    for rangeSize in "${rangeSizes[@]}"; do
        for valueSize in "${valueSizes[@]}"; do
            for human_treeCap in "${treeCapacities[@]}"; do
                mkdir -p ${dataPath}
                rm -rf ${dataPath}/*
                treeCap=$(human_to_number "$human_treeCap")
                Log="${outputDir}/range_brk_${human_op}_${rangeSize}_${valueSize}_${timestamp}.log"
                Cmd="${builddir}/build_release_${db_name}/vidb LatencyBreakdownIteratorQuery --n_1=$treeCap --n_2=$operationCount --n_3=$rangeSize --valueSize=$valueSize --keySize=$keySize --dataPath=$dataPath" 
                run_benchmark "$Cmd" "$Log" "zipf test (ops=${human_op}, batch=${batchSize}, zipfian=${zipf})"
                sleep 5
            done 
        done 
    done
done

# Calculate and display total execution time
total_end_time=$(date +%s)
total_execution_time=$((total_end_time - total_start_time))
echo "All benchmarks completed successfully in $total_execution_time seconds. Results saved in $outputDir/"
