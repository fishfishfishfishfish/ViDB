use std::collections::HashMap;

use crate::{
    def::{DEFAULT_FILE_SIZE, ENTRY_BASE_LENGTH, LEAF_COUNT_IN_TWIG, SMALL_BUFFER_SIZE},
    entryfile::{entry, Entry},
};

use super::{
    tree::{NodePos, Tree},
    twig::Twig,
};

pub fn build_test_tree(
    dir_name: &str,
    deact_sn_list: &Vec<u64>,
    count_before: i32,
    count_after: i32,
) -> (Tree, Vec<Result<i64, std::io::Error>>, u64, Vec<Vec<u8>>) {
    let mut entry_bzs = Vec::new();
    let mut tree = Tree::new(
        0,
        SMALL_BUFFER_SIZE as usize,
        DEFAULT_FILE_SIZE,
        dir_name.to_string(),
        "".to_string(),
        true,
        None,
    );
    let mut entry = Entry {
        key: b"key".as_slice(),
        value: b"value".as_slice(),
        next_key_hash: &pad32(b"nextkey".as_slice()),
        version: 100,
        serial_number: 0,
    };

    let total_len0 = ((ENTRY_BASE_LENGTH + &entry.key.len() + &entry.value.len() + 7) / 8) * 8;
    let mut bz = vec![0u8; total_len0];
    let entry_bz = entry::entry_to_bytes(&entry, &[], &mut bz);

    let mut pos_list = Vec::with_capacity((LEAF_COUNT_IN_TWIG + 10) as usize);
    pos_list.push(tree.append_entry(&entry_bz));
    entry_bzs.push(entry_bz.bz.to_vec());

    for i in 1..count_before {
        entry.serial_number = i as u64;
        let entry_bz = entry::entry_to_bytes(&entry, &[], &mut bz);
        pos_list.push(tree.append_entry(&entry_bz));
        entry_bzs.push(entry_bz.bz.to_vec());
    }
    for sn in deact_sn_list {
        tree.deactive_entry(*sn);
    }

    let total_len1 = total_len0 + &deact_sn_list.len() * 8;
    let mut bz1 = vec![0u8; total_len1];
    entry.serial_number = count_before as u64;
    let entry_bz = entry::entry_to_bytes(&entry, deact_sn_list.as_slice(), &mut bz1);
    pos_list.push(tree.append_entry(&entry_bz));
    entry_bzs.push(entry_bz.bz.to_vec());

    for _ in 0..count_after - 1 {
        entry.serial_number += 1;
        let entry_bz = entry::entry_to_bytes(&entry, &[], &mut bz);
        pos_list.push(tree.append_entry(&entry_bz));
        entry_bzs.push(entry_bz.bz.to_vec());
    }

    (tree, pos_list, entry.serial_number, entry_bzs)
}

pub fn pad32(bz: &[u8]) -> [u8; 32] {
    let mut res = [0; 32];
    res[..bz.len()].copy_from_slice(bz);
    res
}

#[allow(dead_code)]
fn compare_tree_nodes(tree_a: &Tree, tree_b: &Tree) {
    if tree_a.upper_tree.nodes.len() != tree_b.upper_tree.nodes.len() {
        panic!("Different nodes count");
    }
    let mut all_same = true;
    for nth_nodes in tree_a.upper_tree.nodes.clone() {
        for hash_a in nth_nodes {
            for pos in hash_a.keys() {
                if tree_b.upper_tree.get_node(*pos).unwrap()
                    != tree_b.upper_tree.get_node(*pos).unwrap()
                {
                    println!("Different Hash {}-{}", pos.level(), pos.nth());
                    all_same = false;
                }
            }
        }
    }
    if !all_same {
        panic!("Nodes Differ");
    }
}

#[allow(dead_code)]
fn compare_tree_twigs(tree_a: &Tree, tree_b: &mut Tree) {
    for active_bit_shards in &tree_a.active_bit_shards {
        for twig_id in active_bit_shards.keys() {
            if tree_a.get_active_bits(*twig_id) != tree_b.get_active_bits(*twig_id) {
                panic!("Twig ID {} not found in tree_b", twig_id);
            }
        }
    }

    if tree_a.upper_tree.active_twig_shards.len() != tree_b.upper_tree.active_twig_shards.len() {
        panic!("Different twig shards count");
    }
    for twig_map in &tree_a.upper_tree.active_twig_shards {
        for (twig_id, twig) in twig_map {
            if let Some(b) = tree_b.upper_tree.get_twig(*twig_id) {
                compare_twig(*twig_id, twig, b);
            } else {
                panic!("Twig ID {} not found in tree_b", twig_id);
            }
        }
    }
}

#[allow(dead_code)]
fn compare_twig(twig_id: u64, a: &Twig, b: &Twig) {
    for i in 0..a.active_bits_mtl1.len() {
        if a.active_bits_mtl1[i] != b.active_bits_mtl1[i] {
            panic!("activeBitsMTL1[{}] differ at twig {}", i, twig_id);
        }
    }
    for i in 0..a.active_bits_mtl2.len() {
        if a.active_bits_mtl2[i] != b.active_bits_mtl2[i] {
            panic!("activeBitsMTL2[{}] differ at twig {}", i, twig_id);
        }
    }
    if a.active_bits_mtl3 != b.active_bits_mtl3 {
        panic!("activeBitsMTL3 differ at twig {}", twig_id);
    }
    if a.left_root != b.left_root {
        panic!("leftRoot differ at twig {}", twig_id);
    }
    if a.twig_root != b.twig_root {
        panic!("twigRoot differ at twig {}", twig_id);
    }
}

#[allow(dead_code)]
fn compare_nodes(
    nodes_a: &Vec<Vec<HashMap<NodePos, [u8; 32]>>>,
    nodes_b: &Vec<Vec<HashMap<NodePos, [u8; 32]>>>,
) {
    assert_eq!(nodes_a.len(), nodes_b.len(), "Nodes count differs");
    let mut all_same = true;
    for (level, nth_nodes) in nodes_b.iter().enumerate() {
        for (nth, hash_a) in nth_nodes.iter().enumerate() {
            for pos in hash_a.keys() {
                if nodes_a[level][nth].get(pos).unwrap() != nodes_b[level][nth].get(pos).unwrap() {
                    println!("Different Hash {}-{}", pos.level(), pos.nth());
                    all_same = false;
                }
            }
        }
    }
    if !all_same {
        panic!("Nodes Differ");
    }
}

#[allow(dead_code)]
fn compare_twigs_map(twig_map_a: &HashMap<i64, Twig>, twig_map_b: &HashMap<i64, Twig>) {
    assert_eq!(twig_map_a.len(), twig_map_b.len(), "Twigs count differs");
    for (&id, twig_a) in twig_map_a {
        let twig_b = twig_map_b
            .get(&id)
            .expect("Twig ID not found in twig_map_b");
        assert_eq!(
            twig_a.active_bits_mtl1, twig_b.active_bits_mtl1,
            "activeBitsMTL1 differ at twig {}",
            id
        );
        assert_eq!(
            twig_a.active_bits_mtl2, twig_b.active_bits_mtl2,
            "activeBitsMTL2 differ at twig {}",
            id
        );
        assert_eq!(
            twig_a.active_bits_mtl3, twig_b.active_bits_mtl3,
            "activeBitsMTL3 differ at twig {}",
            id
        );
        assert_eq!(
            twig_a.left_root, twig_b.left_root,
            "leftRoot differ at twig {}",
            id
        );
        assert_eq!(
            twig_a.twig_root, twig_b.twig_root,
            "twigRoot differ at twig {}",
            id
        );
    }
}
