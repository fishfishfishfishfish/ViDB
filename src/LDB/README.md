# build
To build the LedgerDB, SQL Ledger, and QLDB, please enter the `exps` directory and run the build script:
```bash
cd exps
./build_release_all.sh
```
The script will build all three ledger databases, and place the binaries in the `${respository_root}/builds/build_release_ledgerdb/bin`, `${respository_root}/builds/build_release_sqlledger/bin`, and `${respository_root}/builds/build_release_qldb/bin` directory.

# run experiments

## Point query and updates
You can start the evaluation using the following commands.
```bash
cd exps
./micro_benchmark_all.sh 
```
This will run experiments, and summarize the results in `${respository_root}/results/results_{database}/micro_benchamrk_{timestamp}_summary.csv`.`{database}` can be one of `LedgerDB`, `SQL Ledger`, and `QLDB`, and `{timestamp}` is the time stamp when the experiment starts.

## Range query
You can start the evaluation using the following commands.
```bash
cd exps
./range_benchmark_all.sh 
```
This will run experiments, and summarize the results in `${respository_root}/results/results_{database}/range_benchmark_{timestamp}_summary.csv`.
`{database}` and `{timestamp}` work the same as in the point query and update experiments.

## Historical version query
You can start the evaluation using the following commands.
```bash
cd exps
./lineage_benchmark_all.sh 
```
This will run experiments, and summarize the results in `${respository_root}/results/results_{database}/lineage_benchmark_{timestamp}_summary.csv`.
`{database}` and `{timestamp}` work the same as in the point query and update experiments.

## Version pruning
You can start the evaluation using the following commands.
```bash
cd exps
./update_benchmark_all.sh 
```
This will run experiments, and summarize the results in `${respository_root}/results/results_{database}/update_benchmark_{timestamp}_summary.csv`.
`{database}` and `{timestamp}` work the same as in the point query and update experiments.
