use clap::{Args, Parser};
use serde::Serialize;

#[derive(Parser, Serialize, Debug, Clone)]
pub struct Cli {
    /// Number of entries to populate database with
    #[arg(short, long, default_value_t = 500_000_000)]
    pub entry_count: u64,

    /// Directory to store the database's persistent files (e.g., /mnt/nvme/QMDB)
    #[arg(long, default_value = "/tmp/QMDB")]
    pub db_dir: String,

    /// Source of randomness for workload generation  
    #[arg(short, long, default_value = "./randsrc.dat")]
    pub randsrc_filename: String,

    /// Output file for benchmark results
    #[arg(short, long, default_value = "results.json")]
    pub output_filename: String,

    /// Number of changesets per task/transaction
    #[arg(long, default_value_t = 2, alias = "changesets-per-txn")]
    pub changesets_per_task: u64,

    /// Target number of operations per block. Will be rounded to nearest achievable number based on changesets-per-task and used to calculate transactions (tasks) per block.
    #[arg(long, default_value_t = 100_000)]
    pub ops_per_block: u64,

    #[command(flatten)]
    pub hover: HoverWorkloadArgs,

    /// Number of blocks to run TPS benchmarks
    #[arg(long, default_value_t = if cfg!(debug_assertions) { 10 } else { 50 })]
    pub tps_blocks: u64,
}

#[derive(Args, Debug, Serialize, Clone)]
pub struct HoverWorkloadArgs {
    // TODO: Make this a function of entries rather than by block
    /// Interval between hover tasks (dev=10, test=1000)
    #[arg(long, default_value_t = if cfg!(debug_assertions) { 10 } else { 10000 })]
    pub hover_interval: u64,

    /// Number of blocks? to insert in hover tasks (dev=2, test=50)
    #[arg(long, default_value_t = if cfg!(debug_assertions) { 2 } else { 50 })]
    pub hover_recreate_block: u64,

    /// Number of blocks? to write in hover tasks (dev=2, test=50)
    #[arg(long, default_value_t = if cfg!(debug_assertions) { 2 } else { 50 })]
    pub hover_write_block: u64,

    /// Number of read latency samples to collect
    #[arg(long, default_value_t = 100000)]
    pub num_read_latency_samples: u64,
}

#[derive(Parser, Serialize, Debug, Clone)]
pub struct MicroBenchCli {
    /// Number of entries to populate database with
    #[arg(short, long, default_value_t = 5000)]
    pub entry_count: u64,

    /// Directory to store the database's persistent files (e.g., /mnt/nvme/QMDB)
    #[arg(long, default_value = "/tmp/QMDB")]
    pub db_dir: String,

    /// Source of randomness for workload generation  
    #[arg(short, long, default_value = "./randsrc.dat")]
    pub randsrc_filename: String,

    /// Output file for benchmark results
    #[arg(short, long, default_value = "results.csv")]
    pub output_filename: String,

    /// Number of changesets per task/transaction
    #[arg(long, default_value_t = 1, alias = "changesets-per-txn")]
    pub changesets_per_task: u64,

    /// Target number of operations per block. Will be rounded to nearest achievable number based on changesets-per-task and used to calculate transactions (tasks) per block.
    #[arg(long, default_value_t = 500)]
    pub ops_per_block: u64,

    /// Number of blocks to run updates
    #[arg(long, default_value_t = if cfg!(debug_assertions) { 10 } else { 50 })]
    pub tps_blocks: u64,

    /// size of the key
    #[arg(long, default_value_t = 32)]
    pub key_size: u64,

    /// size of the value
    #[arg(long, default_value_t = 1024)]
    pub val_size: u64,
}

#[derive(Parser, Serialize, Debug, Clone)]
pub struct UpdateBenchCli {
    /// Number of entries to populate database with
    #[arg(short, long, default_value_t = 5000)]
    pub entry_count: u64,

    /// Directory to store the database's persistent files (e.g., /mnt/nvme/QMDB)
    #[arg(long, default_value = "/tmp/QMDB")]
    pub db_dir: String,

    /// Source of randomness for workload generation  
    #[arg(short, long, default_value = "./randsrc.dat")]
    pub randsrc_filename: String,

    /// Output file for benchmark results
    #[arg(short, long, default_value = "results.csv")]
    pub output_filename: String,

    /// Number of changesets per task/transaction
    #[arg(long, default_value_t = 1, alias = "changesets-per-txn")]
    pub changesets_per_task: u64,

    /// Number of blocks to run updates
    #[arg(long, default_value_t = if cfg!(debug_assertions) { 10 } else { 50 })]
    pub tps_blocks: u64,

    /// List of versions for lineage tests
    #[arg(long, value_parser, num_args = 1.., value_delimiter = ',')]
    pub version_list: Vec<u64>,

    /// size of the key
    #[arg(long, default_value_t = 32)]
    pub key_size: u64,

    /// size of the value
    #[arg(long, default_value_t = 1024)]
    pub val_size: u64,
}

#[derive(Parser, Serialize, Debug, Clone)]
pub struct RangeBenchCli {
    /// Number of entries to populate database with
    #[arg(short, long, default_value_t = 5000)]
    pub entry_count: u64,

    /// Directory to store the database's persistent files (e.g., /mnt/nvme/QMDB)
    #[arg(long, default_value = "/tmp/QMDB")]
    pub db_dir: String,

    /// Source of randomness for workload generation  
    #[arg(short, long, default_value = "./randsrc.dat")]
    pub randsrc_filename: String,

    /// Output file for benchmark results
    #[arg(short, long, default_value = "results.csv")]
    pub output_filename: String,

    /// Target number of operations per block. Will be rounded to nearest achievable number based on changesets-per-task and used to calculate transactions (tasks) per block.
    #[arg(long, default_value_t = 500)]
    pub ops_per_block: u64,

    /// Number of changesets per task/transaction
    #[arg(long, default_value_t = 1, alias = "changesets-per-txn")]
    pub changesets_per_task: u64,

    /// List of versions for lineage tests
    #[arg(long, value_parser, num_args = 1.., value_delimiter = ',')]
    pub range_list: Vec<u64>,

    /// size of the key
    #[arg(long, default_value_t = 32)]
    pub range_test_count: u64,

    /// size of the key
    #[arg(long, default_value_t = 32)]
    pub key_size: u64,

    /// size of the value
    #[arg(long, default_value_t = 1024)]
    pub val_size: u64,
}
