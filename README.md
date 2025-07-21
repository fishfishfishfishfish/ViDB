# ViDB: Cost-Efficient Ledger Database At Scale
This repository includes the implementations of Ledger databases, which are used in the paper "[ViDB: Cost-Efficient Ledger Database At Scale](doc/)".

## Code Structure

├── data/&emsp;&emsp;&emsp;&emsp;*directory resevered for storing data*<br>
├── results/&emsp;&emsp;&emsp;*directory resevered for storing experiment results*<br>
├── builds/&emsp;&emsp;&emsp;&nbsp;*directory resevered for storing builded executables*<br>
├── logs/&emsp;&emsp;&emsp;&emsp;*directory resevered for storing logs*<br>
├── doc/&emsp;&emsp;&emsp;&emsp;&nbsp;*the paper*<br>
├── src/<br>
&emsp;&emsp;├── LDB/&emsp;&emsp;&emsp;&nbsp;*source code of LedgerDB, SQL Ledger, and QLDB, a fork of [LedgerDatabase](https://github.com/nusdbsystem/LedgerDatabase)*<br>
&emsp;&emsp;├── QMDB/&emsp;&emsp;*source code of QMDB, a fork of [QMDB](https://github.com/LayerZero-Labs/qmdb)*<br>
&emsp;&emsp;├── LETUS/&emsp;&emsp;&nbsp;*source code of LETUS, reimplemented by us*<br>
&emsp;&emsp;├── ViDB/&emsp;&emsp;&emsp;*source code of ViDB, proposed in our work*<br>
&emsp;&emsp;└── tools/ <br>
├── .gitignore<br>
├── README.md<br>
├── LICENSE<br>
└── CMakeLists.txt<br>


## Dependency
* Ubuntu 22.04 LTS
* rocksdb (&geq; 5.8)
* boost (&geq; 1.67)
* protobuf (&geq; 2.6.1)
* libevent (&geq; 2.1.12)
* cryptopp (&geq; 6.1.0)
* cargo (1.84.1 66221abde 2024-11-19)
* Intel Threading Building Block (tbb_2020 version)
* openssl (&geq; 3.0.2)
* cmake (&geq; 3.12.2)
* gcc (&geq; 5.5)
* make (&geq; 4.3)
* unzip (6.00 of 20)
* python (&geq; 3.8.10)
* linux-libc-dev
* libclang-dev
* libjemalloc-dev

## Setup and Run experiments
To setup the experiments, we provide a script to build the database for each Ledger database by just running.
```bash
./build.sh
```
The builded database will be put in the folder `builds/build_release_{database_name}`.

To run the experiments, we provides scripts to run the benchmarks used in our paper.
So that you can run a benchmark by just running.
```
./run_{benchmark_name}_benchmark.sh
```
Then the experiment results will be put in the folder `results/result_{database_name}/{benchmark_name}`.

For details, the instructions to setup and run experiments are in the folder of each ledger databases.
- To run experiments on ViDB, please see [ViDB run experiment](ViDB/README.md).
- To run experiments on LedgerDB, SQL Ledger, and QLDB, please see [LDB run experiment](LDB/README.md).
- To run experiments on QMDB, please see [QMDB run experiment](QMDB/README.md).
- To run experiments on LETUS, please see [LETUS run experiment](LETUS/README.md).

