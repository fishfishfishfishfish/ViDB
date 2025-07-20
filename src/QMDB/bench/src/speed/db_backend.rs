use parking_lot::RwLock;
use qmdb::test_helper::SimpleTask;
use std::path::Path;

#[cfg(feature = "use_mdbx")]
pub const NAME: &str = "MDBX";

#[cfg(feature = "use_rocksdb")]
pub const NAME: &str = "RocksDB";

#[cfg(all(not(feature = "use_mdbx"), not(feature = "use_rocksdb")))]
pub const NAME: &str = "QMDB";

#[cfg(feature = "use_mdbx")]
pub fn init(db_dir: &str, _: usize) {
    let path = Path::new(db_dir).join("MDBX");
    crate::speed::mdbx::init(path.to_str().unwrap());
}

#[cfg(feature = "use_rocksdb")]
pub fn init(db_dir: &str, _: usize) {
    let path = Path::new(db_dir).join("ROCKSDB");
    crate::speed::rocksdb::init(path.to_str().unwrap());
}

#[cfg(all(not(feature = "use_mdbx"), not(feature = "use_rocksdb")))]
pub fn init(db_dir: &str, table_id: usize) {
    let path = Path::new(db_dir).join("ADS");
    crate::speed::qmdb::init(path.to_str().unwrap(), table_id);
}

#[cfg(feature = "use_mdbx")]
pub fn create_kv(_: usize, height: i64, task_list: Vec<RwLock<Option<SimpleTask>>>) {
    crate::speed::mdbx::create_kv(height, task_list);
}

#[cfg(feature = "use_rocksdb")]
pub fn create_kv(_: usize, height: i64, task_list: Vec<RwLock<Option<SimpleTask>>>) {
    crate::speed::rocksdb::create_kv(height, task_list);
}

#[cfg(all(not(feature = "use_mdbx"), not(feature = "use_rocksdb")))]
pub fn create_kv(table_id: usize, height: i64, task_list: Vec<RwLock<Option<SimpleTask>>>) {
    crate::speed::qmdb::create_kv(table_id, height, task_list);
}

#[cfg(feature = "use_mdbx")]
pub fn update_kv(_: usize, height: i64, task_list: Vec<RwLock<Option<SimpleTask>>>) {
    crate::speed::mdbx::update_kv(height, task_list);
}

#[cfg(feature = "use_rocksdb")]
pub fn update_kv(_: usize, height: i64, task_list: Vec<RwLock<Option<SimpleTask>>>) {
    crate::speed::rocksdb::update_kv(height, task_list);
}

#[cfg(all(not(feature = "use_mdbx"), not(feature = "use_rocksdb")))]
pub fn update_kv(table_id: usize, height: i64, task_list: Vec<RwLock<Option<SimpleTask>>>) {
    crate::speed::qmdb::update_kv(table_id, height, task_list);
}

#[cfg(feature = "use_mdbx")]
pub fn delete_kv(height: i64, task_list: Vec<RwLock<Option<SimpleTask>>>) {
    crate::speed::mdbx::delete_kv(height, task_list);
}

#[cfg(feature = "use_rocksdb")]
pub fn delete_kv(_: usize, height: i64, task_list: Vec<RwLock<Option<SimpleTask>>>) {
    crate::speed::rocksdb::delete_kv(height, task_list);
}

#[cfg(all(not(feature = "use_mdbx"), not(feature = "use_rocksdb")))]
pub fn delete_kv(table_id: usize, height: i64, task_list: Vec<RwLock<Option<SimpleTask>>>) {
    crate::speed::qmdb::delete_kv(table_id, height, task_list);
}

#[cfg(feature = "use_mdbx")]
pub fn read_kv(_: usize, key_list: &Vec<[u8; 52]>) {
    crate::speed::mdbx::read_kv(key_list);
}

#[cfg(feature = "use_rocksdb")]
pub fn read_kv(_: usize, key_list: &Vec<[u8; 52]>) {
    crate::speed::rocksdb::read_kv(key_list);
}

#[cfg(all(not(feature = "use_mdbx"), not(feature = "use_rocksdb")))]
pub fn read_kv(table_id: usize, height: i64, key_list: &Vec<Vec<u8>>) -> Vec<Vec<u8>> {
    let value_list = crate::speed::qmdb::read_kv(table_id, height, key_list);
    value_list
}

#[cfg(any(feature = "use_mdbx", feature = "use_rocksdb"))]
pub fn flush(_: usize) {}

#[cfg(all(not(feature = "use_mdbx"), not(feature = "use_rocksdb")))]
pub fn flush(table_id: usize) {
    crate::speed::qmdb::flush(table_id);
}
