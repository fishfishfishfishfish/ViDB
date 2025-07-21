#!/bin/bash

# Configuration
dataPath="testdata/letus"
valueSize=1024
rs="5,50,100,200,300,400,500,1000,2000"
cr_values=(0.2 0.5 0.8)
meta_values="2,4,10,20,40"
meta_range="1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26,27,28,29,30,31,32,33,34,35,36,37,38,39,40,41,42,43,44,45,46,47,48,49,50,51,52,53,54,55,56,57,58,59,60,61,62,63,64,65,66,67,68,69,70,71,72,73,74,75,76,77,78,79,80,81,82,83,84,85,86,87,88,89,90,91,92,93,94,95,96,97,98,99,100"
rollback_values="5,10,20,30"
# Human-readable operation counts
operationCounts=("1M")  # Much clearer than 1000000, 10000000
batchSizes=(500)              # Could also use "500" or "0.5K" if preferred

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
            writeLog="${outputDir}/range_prewrite_${human_op}_${batchSize}_${timestamp}.log"
            writeCmd="./letus-vidb write --batchSize=$batchSize --dataPath=$dataPath --operationCount=$operationCount --valueSize=$valueSize"
            run_benchmark "$writeCmd" "$writeLog" "range query pre-write (ops=${human_op}, batch=${batchSize})"
            sleep 5

            RandomGetLog="${outputDir}/random_get_${human_op}_${batchSize}_${timestamp}.log"
            RandomGetCmd="./letus-vidb randomget --batchSize=$batchSize --dataPath=$dataPath --operationCount=$operationCount --valueSize=$valueSize"
            run_benchmark "$RandomGetCmd" "$RandomGetLog" "random-get benchmark (ops=${human_op}, batch=${batchSize})"
            sleep 5

            RandomPutLog="${outputDir}/random_write_${human_op}_${batchSize}_${timestamp}.log"
            RandomPutCmd="./letus-vidb randomput --batchSize=$batchSize --dataPath=$dataPath --operationCount=$operationCount --valueSize=$valueSize"
            run_benchmark "$RandomPutCmd" "$RandomPutLog" "random-put benchmark (ops=${human_op}, batch=${batchSize})"
            sleep 5
        done
#    for batchSize in "${batchSizes[@]}"; do
#        writeLog="${outputDir}/range_prewrite_${human_op}_${batchSize}_${timestamp}.log"
#        writeCmd="./letus-vidb write --batchSize=$batchSize --dataPath=$dataPath --operationCount=$operationCount --valueSize=$valueSize"
#        run_benchmark "$writeCmd" "$writeLog" "range query pre-write (ops=${human_op}, batch=${batchSize})"
#        sleep 5
#
#        rangeLog="${outputDir}/range_query_${human_op}_${batchSize}_${timestamp}.log"
#        rangeCmd="./letus-vidb rang_query --batchSize=$batchSize --dataPath=$dataPath --operationCount=$operationCount --rs=$rs --valueSize=$valueSize"
#        run_benchmark "$rangeCmd" "$rangeLog" "range query benchmark (ops=${human_op}, batch=${batchSize})"
#
#        readLog="${outputDir}/read_${human_op}_${batchSize}_${timestamp}.log"
#        readCmd="./letus-vidb read --batchSize=$batchSize --dataPath=$dataPath --operationCount=$operationCount --rs=$rs --valueSize=$valueSize"
#        run_benchmark "$readCmd" "$readLog" "read benchmark (ops=${human_op}, batch=${batchSize})"
#        sleep 5
#    done
#
#    for batchSize in "${batchSizes[@]}"; do
#        for cr in "${cr_values[@]}"; do
#            logFile="${outputDir}/coldHot_${human_op}_${batchSize}_${cr}_${timestamp}.log"
#            cmd="./letus-vidb coldHot --batchSize=$batchSize --cr=$cr --dataPath=$dataPath --operationCount=$operationCount --valueSize=$valueSize"
#            run_benchmark "$cmd" "$logFile" "cold/hot benchmark (ops=${human_op}, batch=${batchSize}, cr=${cr})"
#            sleep 5
#        done
#    done
#
#    logFile="${outputDir}/dataLineage_${human_op}_${batchSize}_${timestamp}.log"
#    cmd="./letus-vidb dataLineage --batchSize=$batchSize --metas=$meta_values --dataPath=$dataPath --operationCount=$operationCount --valueSize=$valueSize"
#    run_benchmark "$cmd" "$logFile" "dataLineage benchmark (ops=${human_op}, batch=${batchSize})"
#    sleep 5

#   for batchSize in "${batchSizes[@]}"; do
#    logFile="${outputDir}/updateMeta_${human_op}_${timestamp}.log"
#    cmd="./letus-vidb updateMeta --dataPath=$dataPath --metas=$meta_range --batchSize=$batchSize"
#    run_benchmark "$cmd" "$logFile" "updateMeta benchmark (ops=${human_op})"
#    sleep 5
#   done


#
#    logFile="${outputDir}/rollback_${human_op}_${timestamp}.log"
#    cmd="./letus-vidb rollback --dataPath=$dataPath --rollback=$rollback_values"
#    run_benchmark "$cmd" "$logFile" "rollback benchmark (ops=${human_op})"
#    sleep 5

done

# Calculate and display total execution time
total_end_time=$(date +%s)
total_execution_time=$((total_end_time - total_start_time))
echo "All benchmarks completed successfully in $total_execution_time seconds. Results saved in $outputDir/"
