#!/bin/bash
source env.sh
./build_release.sh qldb > ${logdir}/build_release_qldb.log 2>&1
./build_release.sh sqlledger >> ${logdir}/build_release_sqlledger.log 2>&1
./build_release.sh ledgerdb >> ${logdir}/build_release_ledgerdb.log 2>&1
