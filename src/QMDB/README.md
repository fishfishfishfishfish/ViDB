# QMDB
This README provides comprehensive instructions on building and running experiments for QMDB, which is proposed in the paper "QMDB: Quick Merkle Database". The code is forked from https://github.com/LayerZero-Labs/qmdb.

## Build
To build the QMDB, please navigate to the `exps` directory and run the build script:
```bash
cd exps
./build.sh
```
The script will build QMDB.
The resulting binaries are in: `${respository_root}/builds/build_release_qmdb/bin`.

### Experiments
Each experiment will generate results in a CSV file. 
> ℹ️ We use the `{timestamp}` placeholder to represent the start time of each experiment.

### Point query and updates
The following commands start the evaluation:
```bash
cd exps
./run_micro_benchmark.sh 
```
The results will be stored in `${respository_root}/results/results_qmdb/micro_benchamrk_{timestamp}_summary.csv`.


### Range query
The following commands start the evaluation:
```bash
cd exps
./run_range_benchmark.sh 
```
The results will be stored in `${respository_root}/results/results_qmdb/range_benchamrk_{timestamp}_summary.csv`.

### Historical version query
The following commands start the evaluation:
```bash
cd exps
./run_lineage_benchmark.sh 
```
The results will be stored in `${respository_root}/results/results_qmdb/lineage_benchamrk_{timestamp}_summary.csv`.

## Version pruning
The following commands start the evaluation:
```bash
cd exps
./run_update_benchmark.sh 
```
The results will be stored in `${respository_root}/results/results_qmdb/update_benchmark_{timestamp}_summary.csv`.