#!/bin/bash

# Configuration
dataPath="testdata/letus"
# Human-readable operation counts
# operationCounts=("1M")  # Much clearer than 1000000, 10000000
operationCounts=("100") 
batchSizes=(10)              # Could also use "500" or "0.5K" if preferred
# batchSizes=(500 1000 2000 3000 4000 5000)              # Could also use "500" or "0.5K" if preferred
# valueSizes=(256 512 1024 2048)
keySize=8
valueSizes=(2048)


# Output directory
outputDir="benchmark_results"
rm -rf "$outputDir"
mkdir -p "$outputDir"

# Timestamp for unique filenames
timestamp=$(date +"%Y%m%d_%H%M%S")

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
            WriteLog="${outputDir}/write_${human_op}_${batchSize}_${valueSize}_${timestamp}.log"
            WriteCmd="./letus-vidb write --batchSize=$batchSize --dataPath=$dataPath --operationCount=$operationCount --valueSize=$valueSize --keySize=$keySize"
            run_benchmark "$WriteCmd" "$WriteLog" "write benchmark (ops=${human_op}, batch=${batchSize})"
            sleep 5
        done 
    done
done

# Calculate and display total execution time
total_end_time=$(date +%s)
total_execution_time=$((total_end_time - total_start_time))
echo "All benchmarks completed successfully in $total_execution_time seconds. Results saved in $outputDir/"
