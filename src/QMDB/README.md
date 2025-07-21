# build
To build the QMDB for experiments, please enter the `exps` directory and run the build script:
```bash
cd exps
./build.sh
```
The script will build QMDB, and place the binaries in the `${respository_root}/builds/build_release_qmdb/bin` directory.

# run experiments
To run the experiments, you will need to have the QMDB binaries built as described above. 
Here we detail how to run the experiments to evaluate QMDB.

## Point query and updates
You can start the evaluation using the following commands.
```bash
cd exps
./run_micro_benchmark.sh 
```
This will run experiments, and summarize the results in `${respository_root}/results/results_qmdb/micro_benchamrk_{timestamp}_summary.csv`.
`{timestamp}` is the timestamp of when the experiments started.

## Range query
You can start the evaluation using the following commands.
```bash
cd exps
./run_range_benchmark.sh 
```
This will run experiments, and summarize the results in `${respository_root}/results/results_qmdb/range_benchamrk_{timestamp}_summary.csv`.
`{timestamp}` is the timestamp of when the experiments started.

## Historical version query
You can start the evaluation using the following commands.
```bash
cd exps
./run_lineage_benchmark.sh 
```
This will run experiments, and summarize the results in `${respository_root}/results/results_qmdb/lineage_benchamrk_{timestamp}_summary.csv`.
`{timestamp}` is the timestamp of when the experiments started.

## Version pruning
You can start the evaluation using the following commands.
```bash
cd exps
./run_update_benchmark.sh 
```
This will run experiments, and summarize the results in `${respository_root}/results/results_qmdb/update_benchmark_{timestamp}_summary.csv`.
`{timestamp}` is the timestamp of when the experiments started.