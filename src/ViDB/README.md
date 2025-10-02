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

# write operation test
The following commands start the evaluation:
```bash
./dml_test.sh 
```
The results will be stored in `${respository_root}/results/results_vidb/dml_test_test/`.

# query under different write operations
The following commands start the evaluation:
```bash
./dmlmix_test.sh 
```
The results will be stored in `${respository_root}/results/results_vidb/dmlmix_test_test/`.

# performance for skewed access pattern
The following commands start the evaluation:
```bash
./zipf_test.sh 
```
The results will be stored in `${respository_root}/results/results_vidb/zipf_test_test/`.

# performance for scaling data volume
The following commands start the evaluation:
```bash
./grow_test.sh 
```
The results will be stored in `${respository_root}/results/results_vidb/grow_test_test/`.

# performance under various capacity
The following commands start the evaluation:
```bash
./multi_tree_test.sh 
```
The results will be stored in `${respository_root}/results/results_vidb/multitree_test_test`.


# latency breakdown for point query
```bash
./point_breakdown.sh
```
The results will be stored in `${respository_root}/results/results_vidb/point_brk_test_test/`.


# latency breakdown for range query
```bash
./range_breakdown.sh
```
The results will be stored in `${respository_root}/results/results_vidb/range_brk_test_test/`.


