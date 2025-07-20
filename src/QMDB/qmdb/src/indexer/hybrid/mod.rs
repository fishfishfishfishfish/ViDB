pub mod def;
pub mod file_reader;
pub mod file_writer;
pub mod index_cache;
pub mod overlay;
pub mod ref_unit;
pub mod tempfile;
pub mod unit;

use crate::def::{OP_CREATE, OP_DELETE, OP_READ, OP_WRITE, SHARD_COUNT, SHARD_DIV};
use crate::utils::activebits::ActiveBits;
use aes_gcm::Aes256Gcm;
use byteorder::{BigEndian, ByteOrder};
use def::{
    to_key_pos, MERGER_WAIT, MERGE_DIV, MERGE_RATIO, MERGE_THRES, TEMP_FILE_COUNT, UNIT_COUNT,
    UNIT_GROUP_SIZE,
};
use file_reader::FileReader;
use file_writer::FileWriter;
use parking_lot::RwLock;
use std::path::Path;
use std::sync::atomic::{AtomicBool, AtomicUsize, Ordering};
use std::sync::{Arc, Mutex};
use std::thread;
use std::time;
use tempfile::TempFile;
use unit::Unit;

const ZERO: AtomicUsize = AtomicUsize::new(0);

fn new_temp_file(dir: &str, num: usize, part: usize) -> Arc<RwLock<TempFile>> {
    let fname = format!("{}/{:#010x}-{:#04x}", dir, num, part);
    Arc::new(RwLock::new(TempFile::new(fname)))
}

fn split_k80(k80: &[u8]) -> (usize, u64) {
    let idx = BigEndian::read_u16(&k80[..2]) as usize;
    let k = BigEndian::read_u64(&k80[2..10]);
    (idx, k)
}

// 0 is used for initialization, merger starts from 1
const MERGER_START: usize = 1;

struct Merger {
    file_num: usize,
    f_rd: FileReader,
    f_wr: FileWriter,
    hi_arc: Arc<HybridIndexer>,
    merge_thres: usize,
}

impl Merger {
    fn new(hi_arc: Arc<HybridIndexer>, merge_thres: usize) -> Self {
        let new_file = new_temp_file(&hi_arc.dir, MERGER_START, 0);
        Self {
            file_num: MERGER_START,
            f_rd: FileReader::new(hi_arc.cipher.clone()),
            f_wr: FileWriter::new(new_file, hi_arc.cipher.clone()),
            hi_arc,
            merge_thres,
        }
    }

    fn wait_for_large_enough(&self, shard_id: usize) {
        loop {
            let size = self.hi_arc.sizes[shard_id].load(Ordering::SeqCst);
            let count = self.hi_arc.change_counts[shard_id].load(Ordering::SeqCst);
            if count as u128 * MERGE_DIV as u128 > size as u128 * MERGE_RATIO as u128
                && count > self.merge_thres
            {
                self.hi_arc.change_counts[shard_id].fetch_sub(count, Ordering::SeqCst);
                break;
            }
            thread::sleep(time::Duration::from_millis(MERGER_WAIT));
        }
    }

    fn run(&mut self) {
        loop {
            let mut last_shard_id = SHARD_COUNT + 1;
            for idx in 0..UNIT_COUNT {
                let shard_id = idx / SHARD_DIV;
                if last_shard_id != shard_id {
                    self.wait_for_large_enough(shard_id);
                    last_shard_id = shard_id;
                }

                let first = idx == 0 && self.file_num == MERGER_START;
                if idx % UNIT_GROUP_SIZE == 0 && !first {
                    let j = idx / UNIT_GROUP_SIZE;
                    let new_file = new_temp_file(&self.hi_arc.dir, self.file_num, j);
                    self.f_wr.load_file(new_file);
                }

                let mut unit = self.hi_arc.units[idx].lock().unwrap();
                unit.merge(&mut self.f_rd, &mut self.f_wr, idx);
            }
            self.file_num += 1;
        }
    }
}

pub struct HybridIndexer {
    dir: String,
    initializing: AtomicBool,
    units: Vec<Mutex<Unit>>,
    // for merging control
    sizes: [AtomicUsize; SHARD_COUNT],
    change_counts: [AtomicUsize; SHARD_COUNT],
    activebits: Vec<ActiveBits>,
    cipher: Arc<Option<Aes256Gcm>>,
}

impl HybridIndexer {
    pub fn new(n: usize) -> Self {
        if n != UNIT_COUNT {
            panic!("HybridIndexer must have {} units", UNIT_COUNT);
        }
        Self::_new("default_hybrid_dir".to_string(), Arc::new(None))
    }

    pub fn with_dir_and_cipher(dir: String, cipher: Arc<Option<Aes256Gcm>>) -> Self {
        Self::_new(dir, cipher)
    }

    pub fn with_dir(dir: String) -> Self {
        Self::_new(dir, Arc::new(None))
    }

    fn _new(dir: String, cipher: Arc<Option<Aes256Gcm>>) -> Self {
        if Path::new(&dir).exists() {
            std::fs::remove_dir_all(&dir).unwrap();
        }
        std::fs::create_dir(&dir).unwrap();
        let mut files = Vec::with_capacity(TEMP_FILE_COUNT);
        for i in 0..TEMP_FILE_COUNT {
            let f = new_temp_file(&dir, 0, i);
            files.push(f);
        }

        let mut units = Vec::with_capacity(UNIT_COUNT);
        for i in 0..UNIT_COUNT {
            let f = files[i / UNIT_GROUP_SIZE].clone();
            units.push(Mutex::new(Unit::new(f, cipher.clone())));
        }
        let mut v = Vec::new();
        for _ in 0..SHARD_COUNT {
            v.push(ActiveBits::with_capacity(1000));
        }

        Self {
            dir,
            initializing: AtomicBool::new(true),
            units,
            sizes: [ZERO; SHARD_COUNT],
            change_counts: [ZERO; SHARD_COUNT],
            activebits: v,
            cipher,
        }
    }

    pub fn dump_mem_to_file(&self, shard_id: usize) {
        if !self.initializing.load(Ordering::SeqCst) {
            panic!("Cannot dump_mem_to_file after initializing");
        }
        let unit_start = shard_id * SHARD_DIV;
        let unit_end = unit_start + SHARD_DIV;
        let mut f_wr = {
            let unit = self.units[unit_start].lock().unwrap();
            FileWriter::new(unit.ifof.clone(), self.cipher.clone())
        };
        for idx in unit_start..unit_end {
            let mut unit = self.units[idx].lock().unwrap();
            if idx != 0 && idx % UNIT_GROUP_SIZE == 0 {
                f_wr.load_file(unit.ifof.clone());
            }
            unit.start_pos = f_wr.get_file_size();
            let mut first_k_list = (Vec::with_capacity(0), Vec::with_capacity(0));
            // take unit.first_k_list out
            std::mem::swap(&mut unit.first_k_list, &mut first_k_list);
            first_k_list.0.clear();
            first_k_list.1.clear();
            for &(k, pos) in unit.overlay.new_kv.iter() {
                let key_pos = to_key_pos(k, pos);
                f_wr.write(&key_pos, &mut first_k_list);
            }
            first_k_list.0.shrink_to_fit();
            first_k_list.1.shrink_to_fit();
            // return unit.first_k_list back
            std::mem::swap(&mut unit.first_k_list, &mut first_k_list);
            unit.overlay.new_kv.clear();
            f_wr.flush();
            unit.end_pos = f_wr.get_file_size();
        }
        if shard_id == SHARD_COUNT - 1 {
            self.initializing.store(false, Ordering::SeqCst);
        }
    }

    pub fn start_compacting(hi: Arc<HybridIndexer>) {
        let mut merger = Merger::new(hi, MERGE_THRES);
        thread::spawn(move || {
            merger.run();
        });
    }

    pub fn len(&self, shard_id: usize) -> usize {
        self.sizes[shard_id].load(Ordering::SeqCst)
    }

    pub fn add_kv(&self, k80: &[u8], v_in: i64, sn: u64) {
        let (idx, k) = split_k80(k80);
        //if k80[0]==0x81 && k80[1]==0x8b || k80[2]==0xa1 && k80[3]==0x9a || k80[0]==0x4c && k80[1]==0xe2 {
        //    println!("ADD_KV k80={} v_in={:#010x} sn={:#010x}", hex::encode(k80), v_in, sn);
        //}
        let mut unit = self.units[idx].lock().unwrap();
        if v_in % 8 != 0 {
            panic!("value not 8x v_in={}", v_in);
        }
        let v_in = v_in / 8;
        unit.overlay.add_kv(k, v_in);
        if cfg!(feature = "check_hybridindexer") {
            unit.ref_u.as_mut().unwrap().add_kv(0, k, v_in);
        }
        //if k80[0]==0x81 && k80[1]==0x8b || k80[2]==0xa1 && k80[3]==0x9a || k80[0]==0x4c && k80[1]==0xe2 {
        //    println!("overlay.add_kv k={:#016x}, v_in={:#010x} len={}", k, v_in, unit.overlay.new_kv.len());
        //}
        let shard_id = idx / SHARD_DIV;
        self.activebits[shard_id].set(sn);
        self.sizes[shard_id].fetch_add(1, Ordering::SeqCst);
        self.change_counts[shard_id].fetch_add(1, Ordering::SeqCst);
    }

    pub fn erase_kv(&self, k80: &[u8], v_in: i64, sn: u64) {
        let (idx, k) = split_k80(k80);
        //if k80[0]==0x81 && k80[1]==0x8b || k80[2]==0xa1 && k80[3]==0x9a || k80[0]==0x4c && k80[1]==0xe2 {
        //    println!("ERASE_KV k80={} v_in={:#010x} sn={:#010x}", hex::encode(k80), v_in, sn);
        //}
        if self.initializing.load(Ordering::SeqCst) {
            panic!("Cannot erase_kv during initializing");
        }
        let mut unit = self.units[idx].lock().unwrap();
        if v_in % 8 != 0 {
            panic!("value not 8x");
        }
        let v_in = v_in / 8;
        unit.overlay.erase_kv(k, v_in);
        if cfg!(feature = "check_hybridindexer") {
            unit.ref_u.as_mut().unwrap().erase_kv(0, k, v_in);
        }
        //if k80[0]==0x81 && k80[1]==0x8b || k80[2]==0xa1 && k80[3]==0x9a || k80[0]==0x4c && k80[1]==0xe2 {
        //    println!("overlay.erase_kv k={:#016x}, v_in={:#010x} len={}", k, v_in, unit.overlay.new_kv.len());
        //}
        let shard_id = idx / SHARD_DIV;
        self.activebits[shard_id].clear(sn);
        self.sizes[shard_id].fetch_sub(1, Ordering::SeqCst);
        self.change_counts[shard_id].fetch_add(1, Ordering::SeqCst);
    }

    pub fn change_kv(&self, k80: &[u8], v_old: i64, v_new: i64, sn_old: u64, sn_new: u64) {
        let (idx, k) = split_k80(k80);
        //if k80[0]==0x81 && k80[1]==0x8b || k80[2]==0xa1 && k80[3]==0x9a || k80[0]==0x4c && k80[1]==0xe2 {
        //    println!("CHANGE_KV k80={} v:{:#010x}->{:#010x} sn:{:#010x}->{:#010x}", hex::encode(k80), v_old, v_new, sn_old, sn_new);
        //}
        if self.initializing.load(Ordering::SeqCst) {
            panic!("Cannot change_kv during initializing");
        }
        let mut unit = self.units[idx].lock().unwrap();
        if v_old % 8 != 0 {
            panic!("value not 8x");
        }
        let v_old = v_old / 8;
        if v_new % 8 != 0 {
            panic!("value not 8x");
        }
        let v_new = v_new / 8;
        unit.overlay.change_kv(k, v_old, v_new);
        if cfg!(feature = "check_hybridindexer") {
            unit.ref_u.as_mut().unwrap().change_kv(0, k, v_old, v_new);
        }
        //if k80[0]==0x81 && k80[1]==0x8b || k80[2]==0xa1 && k80[3]==0x9a || k80[0]==0x4c && k80[1]==0xe2 {
        //    println!("overlay.change_kv k={:#016x}, v={:#010x}->{:#010x} len={}", k, v_old, v_new, unit.overlay.new_kv.len());
        //}
        let shard_id = idx / SHARD_DIV;
        self.activebits[shard_id].clear(sn_old);
        self.activebits[shard_id].set(sn_new);
        self.change_counts[shard_id].fetch_add(2, Ordering::SeqCst);
    }

    // for_each is only used by prefetcher for warmup
    pub fn for_each<F>(&self, height: i64, op: u8, k80: &[u8], mut access: F)
    where
        F: FnMut(&[u8], i64) -> bool,
    {
        if op == OP_CREATE || op == OP_DELETE {
            self._for_each_adjacent_value::<F>(height, true, k80, access);
        } else if op == OP_WRITE || op == OP_READ {
            // OP_READ is only for test
            self.for_each_value_warmup(height, k80, |offset| access(k80, offset));
        }
    }

    pub fn for_each_value<F>(&self, h: i64, k80: &[u8], access: F)
    where
        F: FnMut(i64) -> bool,
    {
        self._for_each_value(h, false, k80, access);
    }

    pub fn for_each_value_warmup<F>(&self, h: i64, k80: &[u8], access: F)
    where
        F: FnMut(i64) -> bool,
    {
        self._for_each_value(h, true, k80, access);
    }

    fn _for_each_value<F>(&self, height: i64, warmup: bool, k80: &[u8], mut access: F)
    where
        F: FnMut(i64) -> bool,
    {
        //if k80[0]==0x81 && k80[1]==0x8b || k80[2]==0xa1 && k80[3]==0x9a || k80[0]==0x4c && k80[1]==0xe2 {
        //    println!("FEV k80={}", hex::encode(k80));
        //}
        let (idx, k) = split_k80(k80);
        let mut unit = self.units[idx].lock().unwrap();
        unit.for_each_value(height, warmup, k, k80, |v| access(v * 8));
    }

    pub fn for_each_adjacent_value<F>(&self, h: i64, k80: &[u8], access: F)
    where
        F: FnMut(&[u8], i64) -> bool,
    {
        self._for_each_adjacent_value(h, false, k80, access);
    }

    fn _for_each_adjacent_value<F>(&self, height: i64, warmup: bool, k80: &[u8], mut access: F)
    where
        F: FnMut(&[u8], i64) -> bool,
    {
        //if k80[0]==0x81 && k80[1]==0x8b || k80[2]==0xa1 && k80[3]==0x9a || k80[0]==0x4c && k80[1]==0xe2 {
        //    println!("FEA k80={}", hex::encode(k80));
        //}
        let (idx, k) = split_k80(k80);
        let mut unit = self.units[idx].lock().unwrap();
        let mut buf = [0u8; 10];
        buf[..].copy_from_slice(&k80[..10]);
        unit.for_each_adjacent_value(height, warmup, k, &buf, |k, v| access(k, v * 8));
    }

    pub fn key_exists(&self, k80: &[u8], file_pos: i64, sn: u64) -> bool {
        let (idx, _) = split_k80(k80);
        let shard_id = idx / SHARD_DIV;
        let res = self.activebits[shard_id].get(sn);
        if file_pos % 8 != 0 {
            panic!("value not 8x");
        }
        let file_pos = file_pos / 8;
        if cfg!(feature = "check_hybridindexer") {
            let (idx, k) = split_k80(k80);
            let unit = self.units[idx].lock().unwrap();
            let correct = unit.ref_u.as_ref().unwrap().key_exists(0, k, file_pos);
            if correct && !res {
                panic!(
                    "Mismatch key_exists r={} vs i={}: k80={} pos={:#010x} sn={:#010x}",
                    correct,
                    res,
                    hex::encode(k80),
                    file_pos,
                    sn
                );
            }
        }
        res
    }
}

#[cfg(not(feature = "tee_cipher"))]
#[cfg(test)]
mod tests {
    use super::*;
    use crate::{test_helper::TempDir, utils::byte0_to_shard_id};

    fn to_k80(k: &[u8]) -> [u8; 10] {
        let mut result = [0u8; 10];
        result.copy_from_slice(k);
        result
    }

    impl HybridIndexer {
        pub fn init(&mut self) {
            // let sn = 0i64;
            let mut ones80 = [0xFFu8; 10];
            let mut zero80 = [0u8; 10];
            for shard_id in 0..SHARD_COUNT {
                let sn_start = (shard_id * SHARD_DIV) as u64;
                let sn_end = (shard_id * SHARD_DIV + SHARD_DIV) as u64;
                for sn in sn_start..sn_end {
                    BigEndian::write_u16(&mut zero80[..2], sn as u16);
                    self.add_kv(&zero80[..], sn as i64 * 8, sn);
                    BigEndian::write_u16(&mut ones80[..2], sn as u16);
                    self.add_kv(&ones80[..], sn as i64 * 8, sn);
                }
                self.dump_mem_to_file(shard_id);
            }
        }

        fn get(&self, k80: &[u8]) -> Option<i64> {
            let mut retrieved_value = None;
            self.for_each_value(-1, k80, |offset| {
                retrieved_value = Some(offset);
                true
            });
            retrieved_value
        }
    }

    #[test]
    fn test_split_k80() {
        let k80 = [0x12, 0x34, 0x56, 0x78, 0x9A, 0xBC, 0xDE, 0xF0, 0x11, 0x22];
        let (idx, k) = split_k80(&k80);
        assert_eq!(idx, 0x1234);
        assert_eq!(k, 0x56789ABCDEF01122);
    }

    #[test]
    fn test_new_hybrid_indexer() {
        let dir = "test_new_hybrid_indexer";
        let tmp_dir = TempDir::new(dir);
        let indexer = HybridIndexer::with_dir(dir.to_string());

        assert!(indexer.initializing.load(Ordering::SeqCst));
        assert_eq!(indexer.units.len(), UNIT_COUNT);
        assert_eq!(indexer.sizes.len(), SHARD_COUNT);
        assert_eq!(indexer.change_counts.len(), SHARD_COUNT);
        println!("{}", tmp_dir.list().join("\n"));
    }

    #[test]
    fn test_init_kv() {
        let dir = "test_init_kv";
        let _tmp_dir = TempDir::new(dir);
        let mut indexer = HybridIndexer::with_dir(dir.to_string());
        indexer.init();
    }

    #[test]
    fn test_add_and_retrieve_kv() {
        let dir = "test_add_and_retrieve_kv";
        let _tmp_dir = TempDir::new(dir);
        let mut indexer = HybridIndexer::with_dir(dir.to_string());
        indexer.init();
        let n1 = indexer.sizes[1].load(Ordering::SeqCst);
        let c1 = indexer.change_counts[1].load(Ordering::SeqCst);

        let k80 = to_k80(b"0123456789"); // 10-byte key
        let v_in = 48i64;
        let sn = 123456u64;

        indexer.add_kv(&k80, v_in, sn);
        let shard_id = byte0_to_shard_id(k80[0]);

        assert_eq!(indexer.get(&k80).unwrap(), v_in);
        assert_eq!(indexer.sizes[shard_id].load(Ordering::SeqCst), n1 + 1);
        assert_eq!(
            indexer.change_counts[shard_id].load(Ordering::SeqCst),
            c1 + 1
        );
    }

    #[test]
    fn test_erase_kv() {
        let dir = "test_erase_kv";
        let _tmp_dir = TempDir::new(dir);
        let mut indexer = HybridIndexer::with_dir(dir.to_string());
        indexer.init();
        let n1 = indexer.sizes[1].load(Ordering::SeqCst);
        let c1 = indexer.change_counts[1].load(Ordering::SeqCst);

        // Add and then erase a key-value pair
        let k80 = to_k80(b"1234567890");
        let shard_id = byte0_to_shard_id(k80[0]);
        let v_in = 160i64;
        let sn = 567890u64;

        indexer.add_kv(&k80, v_in, sn);
        assert_eq!(indexer.sizes[shard_id].load(Ordering::SeqCst), n1 + 1);
        assert_eq!(
            indexer.change_counts[shard_id].load(Ordering::SeqCst),
            c1 + 1
        );
        indexer.erase_kv(&k80, v_in, sn);
        assert_eq!(indexer.sizes[shard_id].load(Ordering::SeqCst), n1);
        assert_eq!(
            indexer.change_counts[shard_id].load(Ordering::SeqCst),
            c1 + 2
        );

        // Try to retrieve the erased value
        assert_eq!(indexer.get(&k80), None);
    }

    #[test]
    fn test_change_kv() {
        let dir = "test_change_kv";
        let _tmp_dir = TempDir::new(dir);
        let mut indexer = HybridIndexer::with_dir(dir.to_string());
        indexer.init();

        // Add, change, and retrieve a key-value pair
        let k80 = to_k80(b"2345678901");
        let v_old = 200i64;
        let v_new = 304i64;
        let sn_old = 987654u64;
        let sn_new = 987655u64;

        indexer.add_kv(&k80, v_old, sn_old);
        indexer.change_kv(&k80, v_old, v_new, sn_old, sn_new);

        // Retrieve the changed value
        assert_eq!(indexer.get(&k80).unwrap(), v_new);
    }

    #[test]
    fn test_key_exists() {
        let dir = "test_key_exists";
        let _tmp_dir = TempDir::new(dir);
        let mut indexer = HybridIndexer::with_dir(dir.to_string());
        indexer.init();

        let k80 = to_k80(b"3456789012");
        let v_in = 400i64;
        let sn = 112233u64;
        // println!("aaddf:{}", indexer.key_exists(&k80, v_in, sn));

        assert!(!indexer.key_exists(&k80, v_in, sn));
        indexer.add_kv(&k80, v_in, sn);

        assert!(indexer.key_exists(&k80, v_in, sn));
        assert!(!indexer.key_exists(&k80, v_in, sn + 1)); // Non-existent sn
    }

    #[test]
    fn test_merger() {
        let dir = "test_merger";
        let _tmp_dir = TempDir::new(dir);
        let indexer = HybridIndexer::with_dir(dir.to_string());

        let hi_arc = Arc::new(indexer);
        let mut merger = Box::new(Merger::new(hi_arc.clone(), 0));
        let merger_p = &mut *merger as *mut Merger;

        thread::spawn(move || {
            merger.run();
        });

        thread::sleep(time::Duration::from_millis(2000));

        unsafe {
            let merger = &mut (*merger_p);
            assert_eq!(merger.file_num, 1);
        }

        for shard_id in 0..SHARD_COUNT {
            hi_arc.sizes[shard_id].store(1000, Ordering::SeqCst);
            hi_arc.change_counts[shard_id]
                .store(1000 * MERGE_RATIO / MERGE_DIV + 1, Ordering::SeqCst);
        }

        thread::sleep(time::Duration::from_millis(2000));

        unsafe {
            let merger = &mut (*merger_p);
            assert_eq!(merger.file_num, 2);
        }

        for idx in 0..UNIT_COUNT {
            assert_eq!(
                hi_arc.change_counts[idx / SHARD_DIV].load(Ordering::SeqCst),
                0
            );
        }
    }
}
