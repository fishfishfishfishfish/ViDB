#!/bin/bash
source env.sh
db_name=${1:-vidb}
test_name=${2:-test}
echo "db_name: $db_name, test_name=$test_name"

# define experiment parameters
# small test
operationCounts=("1M") # Human-readable operation counts, n1
tree_capacities=("1M" "700K" "650K" "600K" "550K" "500K" "450K" "400K" "350K" "300K" "250K" "200K" "160K" "150K" "100K" "50K" "20K" "10K") # n2
tree_capacities=("20K" "40K" "60K") # n2
batchSizes=(5000) # n3
valueSizes=(128)          
keySize=32
zipfFactor=0.99
bloomCap=5000
cacheCost="0"
VlogSize="$(($((1<<20)) * 40))"


# directory
dataPath="${datadir}"
outputDir="${resdir}/results_${db_name}/multitree_test_${test_name}"
mkdir -p ${dataPath}
mkdir -p ${outputDir}
rm -rf ${dataPath}/*
rm -rf ${outputDir}/*
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
            for human_treeCap in "${tree_capacities[@]}"; do
                mkdir -p ${dataPath}
                rm -rf ${dataPath}/*
                treeCap=$(human_to_number "$human_treeCap")
                Log="${outputDir}/multitree_${human_op}_${batchSize}_${valueSize}_${human_treeCap}_${timestamp}.log"
                Cmd="${builddir}/build_release_${db_name}/vidb multiTree --n_1=$treeCap --n_2=$operationCount --n_3=$batchSize --zipf=$zipfFactor --valueSize=$valueSize --keySize=$keySize --batchSize=$batchSize --VlogSize=$VlogSize --bloomCap=$bloomCap --cacheCost=$cacheCost --dataPath=$dataPath" 
                run_benchmark "$Cmd" "$Log" "multi-tree test (ops=${human_op}, batch=${batchSize}, treeCap=${treeCap})"
                sleep 5
            done 
        done 
    done
done

# Calculate and display total execution time
total_end_time=$(date +%s)
total_execution_time=$((total_end_time - total_start_time))
echo "All benchmarks completed successfully in $total_execution_time seconds. Results saved in $outputDir/"
