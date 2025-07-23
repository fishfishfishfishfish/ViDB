# LedgerDB, SQL Ledger and QLDB
This README provides comprehensive instructions on building and running experiments for three ledger databases: LedgerDB, SQL Ledger, and QLDB.
- **LedgerDB** is proposed in the paper "LedgerDB: a centralized ledger database for universal audit and verification".
- **SQL Ledger** is proposed in the paper "SQL ledger: Cryptographically verifiable data in azure SQL database".
- **QLDB** is [Amazon Quantum Ledger Database](https://aws.amazon.com/qldb/).

The code is forked from [LedgerDatabase](https://github.com/nusdbsystem/LedgerDatabase).

## Build
To build the LedgerDB, SQL Ledger, and QLDB, please navigate to the `exps` directory and run the build script:
```bash
cd exps
./build_release_all.sh
```
This script will build all three ledger databases. 
The resulting binaries can be found in the following directories:
- LedgerDB: `${respository_root}/builds/build_release_ledgerdb/bin`
- SQL Ledger: `${respository_root}/builds/build_release_sqlledger/bin`
- QLDB: `${respository_root}/builds/build_release_qldb/bin`

## Experiments
Each experiment will generate results in a CSV file. 
> ℹ️ We use the `{database}` placeholder to refer to one of LedgerDB, SQL Ledger, or QLDB, and use  the`{timestamp}` placeholder to represent the start time of each experiment.

### Point query and updates
The following commands start the evaluation:
```bash
cd exps
./micro_benchmark_all.sh 
```
The results will be stored in `${respository_root}/results/results_{database}/micro_benchamrk_{timestamp}_summary.csv`.

### Range query
To start the evaluation:
```bash
cd exps
./range_benchmark_all.sh 
```
The results will be stored in `${respository_root}/results/results_{database}/range_benchmark_{timestamp}_summary.csv`.


### Historical version query
To start the evaluation:
```bash
cd exps
./lineage_benchmark_all.sh 
```
The results will be stored in `${respository_root}/results/results_{database}/lineage_benchmark_{timestamp}_summary.csv`.


### Version pruning
To start the evaluation:
```bash
cd exps
./update_benchmark_all.sh 
```
The results will be stored in `${respository_root}/results/results_{database}/update_benchmark_{timestamp}_summary.csv`.
