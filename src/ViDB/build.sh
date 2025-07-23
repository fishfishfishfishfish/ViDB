#!/bin/bash
source env.sh

mkdir -p ${builddir}/build_release_vidb
rm -rf ${builddir}/build_release_vidb/*

cp vidbsvc/vidb ${builddir}/build_release_vidb/vidb
cp -r hyperbench3/ ${builddir}/build_release_vidb/hyperbench3
cd ${builddir}/build_release_vidb/hyperbench3
tar -xvzf hyperbench.tar.gz hyperbench
rm hyperbench.tar.gz