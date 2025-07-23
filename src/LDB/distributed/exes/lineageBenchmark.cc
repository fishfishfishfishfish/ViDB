// put random keys before fully load data

#include <chrono>
#include <fstream>

#include "distributed/lib/generator.h"
#include "distributed/store/common/backend/versionstore.h"

inline char RandomPrintChar(uint64_t num) {
  static const char charset[] =
      "0123456789"
      "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
      "abcdefghijklmnopqrstuvwxyz";
  return charset[num % (sizeof(charset) - 1)];
}

std::string BuildKeyName(uint64_t key_num, int key_len) {
  std::string key_num_str = std::to_string(key_num);
  int zeros = key_len - key_num_str.length();
  zeros = std::max(0, zeros);
  std::string key_name = "";
  return key_name.append(zeros, '0').append(key_num_str);
}

int main(int argc, char** argv) {
  uint64_t num_accout = 5000;  // 40,000,000(40M) 2,000,000(2M)
  uint64_t load_batch_size = 100;
  uint64_t num_txn_version = 40;
  uint64_t num_txn_account = 50;
  uint64_t key_len = 32;
  uint64_t value_len = 1024;
  bool is_get_with_proof = false;
  std::vector<size_t> query_versions = {2, 4, 10, 20, 40};
  std::string data_path = "data/";
  std::string result_path = "exps/results/test.csv";

  int opt;
  while ((opt = getopt(argc, argv, "a:b:t:z:p:k:v:l:d:i:r:")) != -1) {
    switch (opt) {
      case 'a':  // num_accout
      {
        char* strtolPtr;
        num_accout = strtoul(optarg, &strtolPtr, 10);
        if ((*optarg == '\0') || (*strtolPtr != '\0') || (num_accout <= 0)) {
          std::cerr << "option -a requires a numeric arg\n" << std::endl;
        }
        break;
      }

      case 'b':  // load_batch_size
      {
        char* strtolPtr;
        load_batch_size = strtoul(optarg, &strtolPtr, 10);
        if ((*optarg == '\0') || (*strtolPtr != '\0') ||
            (load_batch_size <= 0)) {
          std::cerr << "option -b requires a numeric arg\n" << std::endl;
        }
        break;
      }

      case 't':  // num_txn
      {
        char* strtolPtr;
        num_txn_version = strtoul(optarg, &strtolPtr, 10);
        if ((*optarg == '\0') || (*strtolPtr != '\0') ||
            (num_txn_version < 0)) {
          std::cerr << "option -t requires a numeric arg\n" << std::endl;
        }
        break;
      }

      case 'z':  // num_txn_account
      {
        char* strtolPtr;
        num_txn_account = strtoul(optarg, &strtolPtr, 10);
        if ((*optarg == '\0') || (*strtolPtr != '\0') ||
            (num_txn_account <= 0)) {
          std::cerr << "option -z requires a numeric arg\n" << std::endl;
        }
        break;
      }

      case 'p':  // is get with proof?
      {
        is_get_with_proof = true;
        break;
      }

      case 'k':  // length of key.
      {
        char* strtolPtr;
        key_len = strtoul(optarg, &strtolPtr, 10);
        if ((*optarg == '\0') || (*strtolPtr != '\0') || (key_len <= 0)) {
          std::cerr << "option -k requires a numeric arg\n" << std::endl;
        }
        break;
      }

      case 'v':  // length of value.
      {
        char* strtolPtr;
        value_len = strtoul(optarg, &strtolPtr, 10);
        if ((*optarg == '\0') || (*strtolPtr != '\0') || (value_len <= 0)) {
          std::cerr << "option -v requires a numeric arg\n" << std::endl;
        }
        break;
      }

      case 'l': {  // range for range query
        std::string str = optarg;
        std::stringstream ss(str);
        std::string item;
        query_versions.clear();
        while (std::getline(ss, item, ',')) {
          try {
            query_versions.push_back(std::stoul(item));
          } catch (const std::invalid_argument& e) {
            std::cerr << "Invalid number format for option -l: " << item
                      << std::endl;
            return 1;
          } catch (const std::out_of_range& e) {
            std::cerr << "Number out of range for option -l: " << item
                      << std::endl;
            return 1;
          }
        }
        break;
      }

      case 'd':  // data path
      {
        data_path = optarg;
        break;
      }

      case 'r':  // result path
      {
        result_path = optarg;
        break;
      }

      default:
        std::cerr << "Unknown argument " << argv[optind] << std::endl;
        break;
    }
  }
  // init database
  // timeout = 100ms, the time between update the tree
  VersionedKVStore store(data_path, 100);

  // prepare result file
  std::ofstream rs_file;
  rs_file.open(result_path, std::ios::trunc);
  rs_file << "version,operation,latency,throughput" << std::endl;
  rs_file.close();
  rs_file.open(result_path, std::ios::app);

  // generate random keys
  uint64_t random_keys[num_txn_account];
  UniformGenerator txn_key_generator(1, num_accout);
  for (int i = 0; i < num_txn_account; i++) {
    random_keys[i] = txn_key_generator.Next();
  }

  uint64_t num_load_version = num_accout / load_batch_size;
  uint64_t version = 0;
  uint64_t cur_version = 0;

  // updates
  for (; version <= num_txn_version; version++) {
    std::vector<std::string> put_keys;
    std::vector<std::string> put_values;
    strongstore::proto::Reply reply;
    auto start = std::chrono::system_clock::now();
    for (int i = 0; i < int(num_txn_account); i++) {
      uint64_t num = random_keys[i];
      std::string key = BuildKeyName(num, key_len);
      std::string val = "";
      val = val.append(value_len, RandomPrintChar(num));
      put_keys.push_back(key);
      put_values.push_back(val);
    }
    store.put(put_keys, put_values, Timestamp(version), &reply);
    auto end = std::chrono::system_clock::now();
    auto duration =
        std::chrono::duration_cast<std::chrono::nanoseconds>(end - start);
    double put_latency = double(duration.count()) *
                         std::chrono::nanoseconds::period::num /
                         std::chrono::nanoseconds::period::den;
    sleep(1);
#ifndef AMZQLDB
    cur_version = reply.values(0).estimate_block();
#endif
#ifdef AMZQLDB
    cur_version = version;
#endif
    std::cout << "version " << version << "/" << cur_version
              << ", put latnecy:" << put_latency << ","
              << "put throughput:" << num_txn_account / put_latency
              << std::endl;
    rs_file << version << ",PUT," << put_latency << ","
            << num_txn_account / put_latency << std::endl;
  }

  // load data
  CounterGenerator key_generator(1);
  for (; version <= num_txn_version + num_load_version; version++) {
    std::vector<std::string> keys;
    std::vector<std::string> values;
    strongstore::proto::Reply reply;
    auto start = std::chrono::system_clock::now();
    for (int i = 0; i < int(load_batch_size); i++) {
      uint64_t num = key_generator.Next();
      std::string key = BuildKeyName(num, key_len);
      std::string val = "";
      val = val.append(value_len, RandomPrintChar(num));
      keys.push_back(key);
      values.push_back(val);
    }
    store.put(keys, values, Timestamp(version), &reply);
    auto end = std::chrono::system_clock::now();
    auto duration =
        std::chrono::duration_cast<std::chrono::nanoseconds>(end - start);
    double load_latency = double(duration.count()) *
                          std::chrono::nanoseconds::period::num /
                          std::chrono::nanoseconds::period::den;
#ifndef AMZQLDB
    cur_version = reply.values(0).estimate_block();
#endif
#ifdef AMZQLDB
    cur_version = version;
#endif
    if (version % 1 == 0) {
      std::cout << "version " << version << "/" << cur_version
                << ", load latnecy:" << load_latency
                << ", load throughput:" << load_batch_size / load_latency
                << std::endl;
      rs_file << version << ",LOAD," << load_latency << ","
              << load_batch_size / load_latency << std::endl;
    }
  }
  sleep(1);

  // version get
  for (size_t nv : query_versions) {
    for (int t = 0; t < int(num_txn_account); t++) {
      strongstore::proto::Reply reply;
      double get_latency = 0;
      std::string key = BuildKeyName(random_keys[t], key_len);
      // test get
      if (is_get_with_proof) {
        std::vector<std::pair<std::string, size_t>> key_nv;
        std::map<uint64_t, std::vector<std::string>> blk_keys;
        auto start = std::chrono::system_clock::now();
        key_nv.push_back(std::make_pair(key, nv));
        store.GetNVersions(key_nv, &reply);
#ifndef AMZQLDB
        for (int vi = 0; vi < reply.values_size(); ++vi) {
          blk_keys[reply.values(vi).estimate_block()].push_back(
              reply.values(vi).key());
        }
        store.GetProof(blk_keys, &reply);
#endif
        auto end = std::chrono::system_clock::now();
        auto duration =
            std::chrono::duration_cast<std::chrono::nanoseconds>(end - start);
        get_latency = double(duration.count()) *
                      std::chrono::nanoseconds::period::num /
                      std::chrono::nanoseconds::period::den;
      } else {
        std::vector<std::pair<std::string, size_t>> ver_keys;
        auto start = std::chrono::system_clock::now();
        ver_keys.push_back(std::make_pair(key, nv));
        store.GetNVersions(ver_keys, &reply);
        auto end = std::chrono::system_clock::now();
        auto duration =
            std::chrono::duration_cast<std::chrono::nanoseconds>(end - start);
        get_latency = double(duration.count()) *
                      std::chrono::nanoseconds::period::num /
                      std::chrono::nanoseconds::period::den;
      }
      std::cout << "nv: " << nv << ", get latnecy:" << get_latency << ","
                << "get throughput:" << nv / get_latency << std::endl;
      rs_file << nv << ",GET," << get_latency << "," << nv / get_latency
              << std::endl;
    }
  }
  std::cout << "finished" << std::endl;
  rs_file.close();

  return 0;
}