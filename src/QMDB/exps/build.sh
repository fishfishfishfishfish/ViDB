#!/bin/bash
export PATH=$PATH:/home/${USER}/.cargo/bin
source env.sh

cargo clean --target-dir=${builddir}/build_release_qmdb
cargo build --release --target-dir=${builddir}/build_release_qmdb

head -c 10M </dev/urandom > ${builddir}/build_release_qmdb/randsrc.dat