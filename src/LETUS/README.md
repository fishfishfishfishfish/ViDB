# LETUS
This README provides comprehensive instructions on building and running experiments for LETUS, which is proposed in the paper "LETUS: A Log-Structured Efficient Trusted Universal BlockChain Storage". We implement LETUS by ourselves based on the detailed description in the paper.

## Build
To build the LETUS, please navigate to the `exps` directory and run the build script:
```bash
cd exps
./build.sh --cxx g++ --build-type release
```
The script will build LETUS, 
The resulting binaries are in: `${respository_root}/builds/build_release_letus/bin`.

## Experiments
Each experiment will generate results in a CSV file. 
> ℹ️ We use the `{timestamp}` placeholder to represent the start time of each experiment.

### Point query and updates
The following commands start the evaluation:
```bash
cd exps
./run_micro_benchmark.sh 
```
The results will be stored in `${respository_root}/results/results_letus/micro_benchamrk_{timestamp}_summary.csv`.

### Range query
The following commands start the evaluation:
```bash
cd exps
./run_range_benchmark.sh 
```
The results will be stored in `${respository_root}/results/results_letus/range_benchamrk_{timestamp}_summary.csv`.

### Historical version query
The following commands start the evaluation:
```bash
cd exps
./run_lineage_benchmark.sh 
```
The results will be stored in `${respository_root}/results/results_letus/lineage_benchamrk_{timestamp}_summary.csv`.

### Version pruning
The following commands start the evaluation:
```bash
cd exps
./run_update_benchmark.sh 
```
The results will be stored in `${respository_root}/results/results_letus/update_benchmark_{timestamp}_summary.csv`.