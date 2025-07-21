#!/bin/bash
source env.sh
db_name=${1:-vidb}
test_name=${2:-test}
echo "db_name: $db_name, test_name=$test_name"

# define experiment parameters
batchSizes=(500 1000 2000 3000 4000 5000)
meta_range="1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26,27,28,29,30,31,32,33,34,35,36,37,38,39,40,41,42,43,44,45,46,47,48,49,50,51,52,53,54,55,56,57,58,59,60,61,62,63,64,65,66,67,68,69,70,71,72,73,74,75,76,77,78,79,80,81,82,83,84,85,86,87,88,89,90,91,92,93,94,95,96,97,98,99,100"
valueSize=1024

# directory
dataPath="${datadir}"
outputDir="${resdir}/results_${db_name}/update_benchmark_${test_name}"
mkdir -p ${dataPath}
mkdir -p ${outputDir}
rm -rf ${dataPath}/*
rm -rf ${outputDir}/*

# unique testname
timestamp=${test_name}

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

for batchSize in "${batchSizes[@]}"; do
    logFile="${outputDir}/updateMeta_${batchSize}_${timestamp}.log"
    cmd="${builddir}/build_release_${db_name}/vidb updateMeta --dataPath=$dataPath --metas=$meta_range --batchSize=$batchSize"
    run_benchmark "$cmd" "$logFile" "updateMeta benchmark (ops=${human_op})"
    sleep 5
done

# Calculate and display total execution time
total_end_time=$(date +%s)
total_execution_time=$((total_end_time - total_start_time))
echo "All benchmarks completed successfully in $total_execution_time seconds. Results saved in $outputDir/"
