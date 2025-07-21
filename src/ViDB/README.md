# build
To build the ViDB for experiments, please enter the `exps` directory and run the build script:
```bash
./build.sh
```
The script will build ViDB, and place the binaries in the `${respository_root}/builds/build_release_vidb/bin` directory.

# run experiments
To run the experiments, you will need to have the ViDB binaries built as described above. 
Here we detail how to run the experiments to evaluate ViDB.

## Point query and updates
You can start the evaluation using the following command.
```bash
./run_micro_benchmark.sh 
```
This will run experiments, and summarize the results in `${respository_root}/results/results_vidb/micro_benchamrk_{timestamp}_summary.csv`.
`{timestamp}` is the timestamp of when the experiments started.

## Range query
You can start the evaluation using the following command.
```bash
./run_range_benchmark.sh 
```
This will run experiments, and summarize the results in `${respository_root}/results/results_vidb/range_benchamrk_{timestamp}_summary.csv`.
`{timestamp}` is the timestamp of when the experiments started.

## Historical version query
You can start the evaluation using the following command.
```bash
./run_lineage_benchmark.sh 
```
This will run experiments, and summarize the results in `${respository_root}/results/results_vidb/lineage_benchamrk_{timestamp}_summary.csv`.
`{timestamp}` is the timestamp of when the experiments started.

## Version pruning
You can start the evaluation using the following command.
```bash
./run_update_benchmark.sh 
```
This will run experiments, and place the results in `${respository_root}/results/results_qmdb/update_benchmark_{timestamp}/updataMeta_{#rec}_{timestamp}.log`.
`{timestamp}` is the timestamp of when the experiments started.
Each `.log` file is created for a specific number of data records (`{#rec}`), containing update latency, index disk usage, total disk usage for each version.

## Version rollback
You can start the evaluation using the following command.
```bash
./run_rollback_benchmark.sh 
```
This will run experiments, and place the results in `${respository_root}/results/results_qmdb/rollback_benchmark_{timestamp}/rollback_{#rec}_{timestamp}.log`.
`{timestamp}` is the timestamp of when the experiments started.
Each `.log` file is created for a specific number of data records (`{#rec}`), containing rollback latency each number of rollback version.

## Storage tiering
You can start the evaluation using the following command.
```bash
./run_tier_benchmark.sh
```
This will run experiments, and summarize the results in `${respository_root}/results/results_vidb/tier_benchamrk_{timestamp}/coldHot_{#rec}_{cold_ratio}_summary.csv`.
`{timestamp}` is the timestamp of when the experiments started.
Each `.log` file is created for a specific number of data records (`{#rec}`) and the ratio of cold data (`{cold_ratio}`), containing query latency under the corresponiding condition.