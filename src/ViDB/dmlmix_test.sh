#!/bin/bash
source env.sh
db_name=${1:-vidb}
test_name=${2:-test}
echo "db_name: $db_name, test_name=$test_name"

# define experiment parameters
# small test
operationCounts=("1M") # Human-readable operation counts, n1
operationCounts2=("9M") # Human-readable operation counts, n2
op_comps=("1.0 0 0" "0 1.0 0" "0 0 1.0" "0 0.5 0.5" "0.5 0 0.5" "0.5 0.5 0" "0.33 0.33 0.34") # p1,p2,p3
# op_comps=("1.0 0 0" "0 1.0 0") # p1,p2,p3
batchSizes=(5000) # n3, random point query key count
rangeSizes=(5000) # n4, range query size
loadBatchSize=5000
insertBatchSize=5000
updateBatchSize=5000
deleteBatchSize=5000
valueSizes=(128)      
keySize=32


# directory
dataPath="${datadir}"
outputDir="${resdir}/results_${db_name}/dmlmix_test_${test_name}"
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
    for human_op2 in "${operationCounts2[@]}"; do
        operationCount2=$(human_to_number "$human_op2")
        for op_comp in "${op_comps[@]}"; do
            read -a portion <<< "$op_comp"
            for batchSize in "${batchSizes[@]}"; do
            for rangeSize in "${rangeSizes[@]}"; do
            for valueSize in "${valueSizes[@]}"; do
                mkdir -p ${dataPath}
                rm -rf ${dataPath}/*
                Log="${outputDir}/dmlmix_${human_op}_${batchSize}_${valueSize}_${op_comp}_${timestamp}.log"
                Cmd="${builddir}/build_release_${db_name}/vidb updateInsertDeleteMix --n_1=$operationCount --n_2=$operationCount2 --n_3=$batchSize --n_4=$rangeSize --p1=${portion[0]} --p2=${portion[1]} --p3=${portion[2]} --batch_size=$loadBatchSize --i_batch_size=$insertBatchSize --u_batch_size=$updateBatchSize --d_batch_size=$deleteBatchSize --value_size=$valueSize --key_size=$keySize --dataPath=$dataPath" 
                run_benchmark "$Cmd" "$Log" "dmlmix test (ops=${human_op}, batch=${batchSize})"
                sleep 5
            done 
            done 
            done 
        done 
    done
done

# Calculate and display total execution time
total_end_time=$(date +%s)
total_execution_time=$((total_end_time - total_start_time))
echo "All benchmarks completed successfully in $total_execution_time seconds. Results saved in $outputDir/"
