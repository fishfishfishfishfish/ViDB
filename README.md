# ViDB: Cost-Efficient Ledger Database At Scale
This repository includes the implementations of Ledger databases, which are used in the paper "[ViDB: Cost-Efficient Ledger Database At Scale](doc/)".

## Code Structure
```
├── data/
├── results/
├── builds/
├── doc/
├── exps/
    ├── scripts/
├── src/
    ├── LDB/    source code of LedgerDB, SQL Ledger, and QLDB, a fork of [LedgerDatabase](https://github.com/nusdbsystem/LedgerDatabase).
    ├── QMDB/  source code of QMDB, a fork of [QMDB]().
    ├── LETUS/  source code of LETUS , reimplemented by ours.
    ├── ViDB/  source code of ViDB,
    └── tools/
├── test/
├── .gitignore
├── README.md
├── LICENSE
```

## Dependency
* rocksdb (&geq; 5.8)
* boost (&geq; 1.67)
* protobuf (&geq; 2.6.1)
* libevent (&geq; 2.1.12)
* cryptopp (&geq; 6.1.0)
* Intel Threading Building Block (tbb_2020 version)
* openssl
* cmake (&geq; 3.12.2)
* gcc (&geq; 5.5)

## Setup


## Run experiments


### Example results
Experiment setting
```
Number of servers (nshards)   : 16
Client per node (nclients)    : 1 2 3 4 5 6
Number of client nodes        : 8
Write percentage (wperc)      : 50
Ziph factor (theta)           : 0
Duration in seconds (rtime)   : 300
Delay in milliseconds (delay) : 100
Transactino size (tlen)       : 10
```

