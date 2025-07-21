#!/bin/bash
source env.sh
db_name=${1:-vidb}
test_name=${2:-test}
echo "db_name: $db_name, test_name=$test_name"

# define experiment parameters
operationCounts=("1M" "10M") # Human-readable operation counts
valueSizes=(1024)
batchSizes=(25000)             
queryVersionsStart="1,2,4,10,20,40"
queryVersionsCount="1,10"
keySize=32


# directory
dataPath="${datadir}"
outputDir="${resdir}/results_${db_name}/lineage_benchmark_${test_name}"
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
for human_op in "${operationCounts[@]}"; do
    operationCount=$(human_to_number "$human_op")
    for batchSize in "${batchSizes[@]}"; do
        for valueSize in "${valueSizes[@]}"; do
            logFile="${outputDir}/dataLineage_${human_op}_${batchSize}_${timestamp}.log"
            cmd="${builddir}/build_release_${db_name}/vidb dataLineage --batchSize=$batchSize --queryVerStart=$queryVersionsStart --queryVerCount=$queryVersionsCount --dataPath=$dataPath --operationCount=$operationCount --valueSize=$valueSize --bfLoad=40 --afLoad=0"
            run_benchmark "$cmd" "$logFile" "dataLineage benchmark (ops=${human_op}, batch=${batchSize})"
            sleep 5
        done 
    done
done

# Calculate and display total execution time
total_end_time=$(date +%s)
total_execution_time=$((total_end_time - total_start_time))
echo "All benchmarks completed successfully in $total_execution_time seconds. Results saved in $outputDir/"
