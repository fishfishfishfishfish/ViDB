# build
To build the LETUS for experiments, please enter the `exps` directory and run the build script:
```bash
cd exps
./build.sh --cxx g++ --build-type release
```
The script will build LETUS, and place the binaries in the `${respository_root}/builds/build_release_letus/bin` directory.

# run experiments
To run the experiments, you will need to have the LETUS binaries built as described above. 
Here we detail how to run the experiments to evaluate LETUS.

## Point query and updates
You can start the evaluation using the following commands.
```bash
cd exps
./run_micro_benchmark.sh 
```
This will run experiments, and summarize the results in `${respository_root}/results/results_letus/micro_benchamrk_{timestamp}_summary.csv`.
`{timestamp}` is the timestamp of when the experiments started.

## Range query
You can start the evaluation using the following commands.
```bash
cd exps
./run_range_benchmark.sh 
```
This will run experiments, and summarize the results in `${respository_root}/results/results_letus/range_benchamrk_{timestamp}_summary.csv`.
`{timestamp}` is the timestamp of when the experiments started.

## Historical version query
You can start the evaluation using the following commands.
```bash
cd exps
./run_lineage_benchmark.sh 
```
This will run experiments, and summarize the results in `${respository_root}/results/results_letus/lineage_benchamrk_{timestamp}_summary.csv`.
`{timestamp}` is the timestamp of when the experiments started.

## Version pruning
You can start the evaluation using the following commands.
```bash
cd exps
./run_update_benchmark.sh 
```
This will run experiments, and summarize the results in `${respository_root}/results/results_letus/update_benchmark_{timestamp}_summary.csv`.
`{timestamp}` is the timestamp of when the experiments started.