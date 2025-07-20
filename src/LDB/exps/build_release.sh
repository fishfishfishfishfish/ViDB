#!/bin/bash

source env.sh

srcdir="${PWD}/.."
if [ $(echo "$1" | awk '{print tolower($0)}') == 'qldb' ]
then
  qldbopt=ON
  ledgerdbopt=OFF
  sqlledgeropt=OFF
  db_build_dir=${builddir}/build_release_qldb
else
  if [ $(echo "$1" | awk '{print tolower($0)}') == 'ledgerdb' ]
  then
    qldbopt=OFF
    ledgerdbopt=ON
    sqlledgeropt=OFF
    db_build_dir=${builddir}/build_release_ledgerdb
  else
    qldbopt=OFF
    ledgerdbopt=OFF
    sqlledgeropt=ON
    db_build_dir=${builddir}/build_release_sqlledger
  fi  
fi

echo ${db_build_dir}
mkdir -p ${db_build_dir}
rm -rf ${db_build_dir}/*
cd ${db_build_dir}

cmake -DLEDGERDB=${ledgerdbopt} -DAMZQLDB=${qldbopt} -DSQLLEDGER=${sqlledgeropt} ${srcdir}
make -j6 VERBOSE=1 