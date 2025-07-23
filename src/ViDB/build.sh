#!/bin/bash
source env.sh

mkdir -p ${builddir}/build_release_vidb
rm -rf ${builddir}/build_release_vidb/*

cp vidbsvc/vidb ${builddir}/build_release_vidb/vidb
