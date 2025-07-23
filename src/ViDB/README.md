# ViDB
This README provides comprehensive instructions on building and running experiments for ViDB, which is proposed in our paper "[ViDB: Cost-Efficient Ledger Database At Scale](doc/)". 

## build
To build the ViDB, please run the build script:
```bash
./build.sh
```
The script will build ViDB.
The resulting binaries are in: `${respository_root}/builds/build_release_vidb/bin`.

# run experiments
Each experiment will generate results in a CSV file or log file. 
> ℹ️ We use the `{timestamp}` placeholder to represent the start time of each experiment.

## Point query and updates
The following commands start the evaluation:
```bash
./run_micro_benchmark.sh 
```
The results will be stored in`${respository_root}/results/results_vidb/micro_benchamrk_{timestamp}_summary.csv`.


## Range query
The following commands start the evaluation:
```bash
./run_range_benchmark.sh 
```
The results will be stored in `${respository_root}/results/results_vidb/range_benchamrk_{timestamp}_summary.csv`.

## Historical version query
The following commands start the evaluation:
```bash
./run_lineage_benchmark.sh 
```
The results will be stored in `${respository_root}/results/results_vidb/lineage_benchamrk_{timestamp}_summary.csv`.

## Version pruning
The following commands start the evaluation:
```bash
./run_update_benchmark.sh 
```
The results will be stored in `${respository_root}/results/results_qmdb/update_benchmark_{timestamp}/updataMeta_{#rec}_{timestamp}.log`.
Each `.log` file is created for a specific number of data records (`{#rec}`), containing update latency, index disk usage, total disk usage for each version.

## Version rollback
The following commands start the evaluation:
```bash
./run_rollback_benchmark.sh 
```
The results will be stored in `${respository_root}/results/results_qmdb/rollback_benchmark_{timestamp}/rollback_{#rec}_{timestamp}.log`.
Each `.log` file is created for a specific number of data records (`{#rec}`), containing rollback latency each number of rollback version.

## Storage tiering
The following commands start the evaluation:
```bash
./run_tier_benchmark.sh
```
The results will be stored in `${respository_root}/results/results_vidb/tier_benchamrk_{timestamp}/coldHot_{#rec}_{cold_ratio}_summary.csv`.
Each `.log` file is created for a specific number of data records (`{#rec}`) and the ratio of cold data (`{cold_ratio}`), containing query latency under the corresponiding condition.
