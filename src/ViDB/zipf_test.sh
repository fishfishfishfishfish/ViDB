#!/bin/bash
source env.sh
db_name=${1:-vidb}
test_name=${2:-test_name}
echo "db_name: $db_name, test_name=$test_name"

# define experiment parameters
# small test
operationCounts=("1M") # Human-readable operation counts, n1
zipfFactors=(0.01 0.1 0.2 0.5 0.8 0.9 0.99 0.999 1.0) # n2, zipf分布的参数
zipfFactors=(0.8 0.9 0.99 0.999) # n2, zipf分布的参数
batchSizes=(5000) # n3, 查询的key的数量
batchSizes=(50) # n3, 查询的key的数量
valueSizes=(128)          
keySize=32
batchCount=5


# directory
dataPath="${datadir}"
outputDir="${resdir}/results_${db_name}/zipf_test_${test_name}"
mkdir -p ${dataPath}
mkdir -p ${outputDir}
rm -rf ${dataPath}/*
rm -rf ${outputDir}/*
# Timestamp for unique filenames
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
            for zipf in "${zipfFactors[@]}"; do
                mkdir -p ${dataPath}
                rm -rf ${dataPath}/*
                Log="${outputDir}/zipf_${human_op}_${batchSize}_${valueSize}_${timestamp}.log"
                Cmd="${builddir}/build_release_${db_name}/vidb get_zipf --n_1=$operationCount --n_3=$batchSize --p1=$zipf --zType=zipfian --batch_count=$batchCount --value_size=$valueSize --key_size=$keySize --dataPath=$dataPath" 
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
