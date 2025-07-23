# ViDB: Cost-Efficient Ledger Database At Scale
This repository includes the implementations of five Ledger databases (ViDB, LETUS, QMDB, LedgerDB, SQL Ledger, QLDB), as presented in in the paper "[ViDB: Cost-Efficient Ledger Database At Scale](doc/)".

## Code Structure

├── data/&emsp;&emsp;&emsp;&emsp;*Reserved directory for data.*<br>
├── results/&emsp;&emsp;&emsp;*Reserved directory for experiment results.*<br>
├── builds/&emsp;&emsp;&emsp;&nbsp;*Reserved directory for built executables.*<br>
├── logs/&emsp;&emsp;&emsp;&emsp;*Reserved directory for execution logs.*<br>
├── doc/&emsp;&emsp;&emsp;&emsp;&nbsp;*Includes the research paper.*<br>
├── src/<br>
&emsp;&emsp;├── LDB/&emsp;&emsp;&emsp;&nbsp;*Source code of LedgerDB, SQL Ledger, and QLDB, forked from [LedgerDatabase](https://github.com/nusdbsystem/LedgerDatabase).*<br>
&emsp;&emsp;├── QMDB/&emsp;&emsp;*Source code of QMDB, forked from [QMDB](https://github.com/LayerZero-Labs/qmdb).*<br>
&emsp;&emsp;├── LETUS/&emsp;&emsp;&nbsp;*Source code of LETUS, re-implemented by the us.*<br>
&emsp;&emsp;└── ViDB/&emsp;&emsp;&emsp;*Source code of ViDB, proposed in our paper.*<br>
├── CMakeLists.txt<br>
├── .gitignore<br>
├── LICENSE<br>
└── README.md: *This documentation file.*<br>


## Dependency
Ensure the following packages are installed:
- **OS**: Ubuntu 22.04 LTS
- **Libraries:**
    * rocksdb (&geq; 5.8)
    * boost (&geq; 1.67)
    * protobuf (&geq; 2.6.1)
    * cryptopp (&geq; 6.1.0)
    * Intel Threading Building Block (tbb_2020 version)
    * openssl (&geq; 3.0.2)
    * libevent (&geq; 2.1.12)
    * linux-libc-dev, libclang-dev, libjemalloc-dev
- **Build Tools:**
    * cargo (1.84.1 66221abde 2024-11-19)
    * cmake (&geq; 3.12.2)
    * gcc (&geq; 5.5)
    * make (&geq; 4.3)
    * python (&geq; 3.8.10)
- **Other tools:**
    * unzip (6.00 of 20)
    * tar (1.34)

## Setup and Run experiments
To set up the experiments, we provide a build script for each ledger database. 
For each ledger database, enter the `exps` folder, and run the following command to compile:
```bash
cd src/{LDB,QMDB,LETUS,ViDB}/exps
./build.sh
```
The built database will be placed in `builds/build_release_{database}`.

Benchmark scripts are provided for each database. To execute a benchmark:
```
./run_{benchmark_name}_benchmark.sh
```
Results will be saved in: `results/result_{database_name}/{benchmark_name}`.

🌟**Note:** Each ledger database may require slightly different commands.
For details, the users can refer to the docs of each ledger databases.
- ViDB: please see [ViDB docs](src/ViDB/README.md).
- LedgerDB, SQL Ledger, and QLDB: please see [LDB docs](src/LDB/README.md).
- QMDB: please see [QMDB docs](src/QMDB/README.md).
- LETUS: please see [LETUS docs](src/LETUS/README.md).

