#include "DMMTrie.hpp"

#include <openssl/evp.h>
#include <openssl/sha.h>

#include <algorithm>
#include <array>
#include <cstring>
#include <iomanip>
#include <iostream>
#include <list>
#include <memory>
#include <set>
#include <sstream>
#include <tuple>
#include <unordered_map>
#include <vector>

#include "LSVPS.hpp"

using namespace std;

// string HashFunction(const string &input) {  // hash function SHA-256
//   EVP_MD_CTX *ctx = EVP_MD_CTX_new();       // create SHA-256 context

//   // initialize SHA-256 hash computation
//   EVP_DigestInit_ex(ctx, EVP_sha256(), nullptr);

//   // update the hash with input string
//   EVP_DigestUpdate(ctx, input.c_str(), input.size());

//   unsigned char hash[EVP_MAX_MD_SIZE];
//   unsigned int hash_len = 0;

//   EVP_DigestFinal_ex(ctx, hash, &hash_len);

//   EVP_DigestFinal_ex(ctx, hash, &hash_len);
//   EVP_MD_CTX_free(ctx);

//   return string(reinterpret_cast<char *>(hash), hash_len);
// }

string HashFunction(const string &input) { // SHA 1
  unsigned char hash[SHA_DIGEST_LENGTH];
  SHA1(reinterpret_cast<const unsigned char*>(input.c_str()),
        input.size(), hash);
  return string(reinterpret_cast<char *>(hash), SHA_DIGEST_LENGTH);
}

auto CompareStrings = [](const std::string &a, const std::string &b) {
  if (a.size() != b.size()) {
    return a.size() > b.size();  // first compare length
  }
  return a < b;  // then compare alphabetical order
};

// convert hexadecimal digit to corresponding index 0~15
int GetIndex(char ch) {
  if (isdigit(ch)) {
    return ch - '0';
  } else if (ch >= 'a' && ch <= 'f') {
    return ch - 'a' + 10;
  } else if (ch >= 'A' && ch <= 'F') {
    return ch - 'A' + 10;
  } else {
    return -1;
  }
}

void Node::CalculateHash() {}
void Node::AddChild(int index, Node *child, uint64_t version,
                    const string &hash) {}
Node *Node::GetChild(int index) const { return nullptr; }
bool Node::HasChild(int index) const { return false; }
void Node::SetChild(int index, uint64_t version, string hash) {}
string Node::GetChildHash(int index) {}
uint64_t Node::GetChildVersion(int index) {}
void Node::UpdateNode() {}
void Node::SetLocation(tuple<uint64_t, uint64_t, uint64_t> location) {}
NodeProof Node::GetNodeProof(int level, int index) {}

LeafNode::LeafNode(uint64_t V, const string &k,
                   const tuple<uint64_t, uint64_t, uint64_t> &l,
                   const string &h)
    : version_(V), key_(k), location_(l), hash_(h), is_leaf_(true) {}

void LeafNode::CalculateHash(const string &value) {
  // hash_ = HashFunction(key_ + value);
  hash_ = HashFunction(value);
}

/* serialized leaf node format (size in bytes):
   | is_leaf_node (1) | version (8) | key_size (8 in 64-bit system) | key
   (key_size) | location(8, 8, 8) | hash (32) |
*/
void LeafNode::SerializeTo(char *buffer, size_t &current_size,
                           bool is_root) const {
  memcpy(buffer + current_size, &is_leaf_,
         sizeof(bool));  // true means that the node is leafnode
  current_size += sizeof(bool);

  memcpy(buffer + current_size, &version_, sizeof(uint64_t));
  current_size += sizeof(uint64_t);

  size_t key_size = key_.size();
  memcpy(buffer + current_size, &key_size, sizeof(key_size));  // key size
  current_size += sizeof(key_size);
  memcpy(buffer + current_size, key_.c_str(), key_size);  // key
  current_size += key_size;

  memcpy(buffer + current_size, &get<0>(location_),
         sizeof(uint64_t));  // fileID
  current_size += sizeof(uint64_t);
  memcpy(buffer + current_size, &get<1>(location_),
         sizeof(uint64_t));  // offset
  current_size += sizeof(uint64_t);
  memcpy(buffer + current_size, &get<2>(location_),
         sizeof(uint64_t));  // size
  current_size += sizeof(uint64_t);

  memcpy(buffer + current_size, hash_.c_str(), hash_.size());
  current_size += HASH_SIZE;
}

void LeafNode::DeserializeFrom(char *buffer, size_t &current_size,
                               bool is_root) {
  version_ = *(reinterpret_cast<uint64_t *>(
      buffer + current_size));  // deserialize leafnode version
  current_size += sizeof(uint64_t);

  size_t key_size = *(reinterpret_cast<size_t *>(
      buffer + current_size));  // deserialize key_size
  current_size += sizeof(key_size);
  key_ = string(buffer + current_size, key_size);  // deserialize key
  current_size += key_size;

  uint64_t fileID = *(reinterpret_cast<uint64_t *>(
      buffer + current_size));  // deserialize fileID
  current_size += sizeof(uint64_t);
  uint64_t offset = *(reinterpret_cast<uint64_t *>(
      buffer + current_size));  // deserialize offset
  current_size += sizeof(uint64_t);
  uint64_t size = *(
      reinterpret_cast<uint64_t *>(buffer + current_size));  // deserialize size
  current_size += sizeof(uint64_t);
  location_ = make_tuple(fileID, offset, size);

  hash_ = string(buffer + current_size, HASH_SIZE);  // deserialize hash
  current_size += HASH_SIZE;
}

void LeafNode::UpdateNode(uint64_t version,
                          const tuple<uint64_t, uint64_t, uint64_t> &location,
                          const string &value, uint8_t location_in_page,
                          DeltaPage *deltapage) {
  version_ = version;
  location_ = location;
  if (value == "") {  // value是空字符串代表Delete节点，此时将哈希改为空串
    hash_ = "";
  } else {
    // hash_ = HashFunction(key_ + value);
    hash_ = HashFunction(value);
  }

  if (deltapage != nullptr) {
    deltapage->AddLeafNodeUpdate(location_in_page, version, hash_,
                                 get<0>(location), get<1>(location),
                                 get<2>(location));
  }
}

tuple<uint64_t, uint64_t, uint64_t> LeafNode::GetLocation() const {
  return location_;
}

void LeafNode::SetLocation(tuple<uint64_t, uint64_t, uint64_t> location) {
  location_ = location;
}

string LeafNode::GetHash() { return hash_; }
uint64_t LeafNode::GetVersion() { return version_; }
void LeafNode::SetVersion(uint64_t version) { version_ = version; }
void LeafNode::SetHash(string hash) { hash_ = hash; }

bool LeafNode::IsLeaf() const { return is_leaf_; }

IndexNode::IndexNode(uint64_t V, const string &h, uint16_t b)
    : version_(V), hash_(h), bitmap_(b), is_leaf_(false) {
  for (size_t i = 0; i < DMM_NODE_FANOUT; i++) {
    children_[i] =
        make_tuple(0, "", nullptr);  // initialize children to default
  }
}

IndexNode::IndexNode(
    uint64_t version, const string &hash, uint16_t bitmap,
    const array<tuple<uint64_t, string, Node *>, DMM_NODE_FANOUT> &children)
    : version_(version),
      hash_(hash),
      bitmap_(bitmap),
      children_(children),
      is_leaf_(false) {}

IndexNode::IndexNode(const IndexNode &other)
    : version_(other.version_),
      hash_(other.hash_),
      bitmap_(other.bitmap_),
      is_leaf_(other.is_leaf_) {
  // Deep copy children array
  for (size_t i = 0; i < DMM_NODE_FANOUT; i++) {
    if (other.HasChild(i)) {
      Node *child = other.GetChild(i);
      if (child != nullptr) {
        if (child->IsLeaf()) {
          children_[i] =
              make_tuple(get<0>(other.children_[i]), get<1>(other.children_[i]),
                         new LeafNode(*dynamic_cast<LeafNode *>(child)));
        } else {
          children_[i] =
              make_tuple(get<0>(other.children_[i]), get<1>(other.children_[i]),
                         new IndexNode(*dynamic_cast<IndexNode *>(child)));
        }
      } else {
        children_[i] = make_tuple(get<0>(other.children_[i]),
                                  get<1>(other.children_[i]), nullptr);
      }

    } else {
      children_[i] = make_tuple(0, "", nullptr);
    }
  }
}

void IndexNode::CalculateHash() {
  string concatenated_hash;
  for (int i = 0; i < DMM_NODE_FANOUT; i++) {
    concatenated_hash += get<1>(children_[i]);
  }
  hash_ = HashFunction(concatenated_hash);
}

/* serialized index node format (size in bytes):
   | is_leaf_node (1) | version (8) | hash (32) | bitmap (2) | Vc (8) | Hc
   (32) | Vc (8) | Hc (32) | ... | child 1 | child 2 | ...
   the function doesn't serialize pointer and doesn't serialize empty child
   nodes
*/
void IndexNode::SerializeTo(char *buffer, size_t &current_size,
                            bool is_root) const {
  memcpy(buffer + current_size, &is_leaf_,
         sizeof(bool));  // false means that the node is indexnode
  current_size += sizeof(bool);

  memcpy(buffer + current_size, &version_, sizeof(uint64_t));
  current_size += sizeof(uint64_t);

  memcpy(buffer + current_size, hash_.c_str(), hash_.size());
  current_size += HASH_SIZE;

  memcpy(buffer + current_size, &bitmap_, sizeof(uint16_t));
  current_size += sizeof(uint16_t);

  for (int i = 0; i < DMM_NODE_FANOUT; i++) {
    if (bitmap_ & (1 << i)) {
      uint64_t child_version = get<0>(children_[i]);
      string child_hash = get<1>(children_[i]);

      memcpy(buffer + current_size, &child_version, sizeof(uint64_t));
      current_size += sizeof(uint64_t);
      memcpy(buffer + current_size, child_hash.c_str(), child_hash.size());
      current_size += HASH_SIZE;
    }
  }

  if (is_root) {  // if an index node is the root node of a page, serialize
                  // its children
    for (int i = 0; i < DMM_NODE_FANOUT; i++) {
      if (bitmap_ & (1 << i)) {  // only serialize children that exists
        Node *child = get<2>(children_[i]);
        child->SerializeTo(buffer, current_size, false);
      }
    }
  }
}

void IndexNode::DeserializeFrom(char *buffer, size_t &current_size,
                                bool is_root) {
  version_ = *(reinterpret_cast<uint64_t *>(buffer + current_size));
  current_size += sizeof(uint64_t);

  hash_ = string(buffer + current_size, HASH_SIZE);
  current_size += HASH_SIZE;

  bitmap_ = *(reinterpret_cast<uint16_t *>(buffer + current_size));
  current_size += sizeof(uint16_t);

  for (int i = 0; i < DMM_NODE_FANOUT; i++) {
    if (bitmap_ & (1 << i)) {
      uint64_t child_version =
          *(reinterpret_cast<uint64_t *>(buffer + current_size));
      current_size += sizeof(uint64_t);
      string child_hash(buffer + current_size, HASH_SIZE);
      current_size += HASH_SIZE;

      children_[i] = make_tuple(child_version, child_hash, nullptr);
    }
  }

  if (!is_root) {  // indexnode is in second level of a page, return
    return;
  }

  for (int i = 0; i < DMM_NODE_FANOUT; i++) {
    if (bitmap_ & (1 << i)) {
      // serialized data only stores children that exists
      bool child_is_leaf_node =
          *(reinterpret_cast<bool *>(buffer + current_size));
      current_size += sizeof(bool);

      if (child_is_leaf_node) {  // second level of page is leafnode
        Node *child = new LeafNode();
        child->DeserializeFrom(buffer, current_size, false);
        this->AddChild(i, child);  // add pointer to children in indexnode
      } else {                     // second level of page is indexnode
        Node *child = new IndexNode();
        child->DeserializeFrom(buffer, current_size, false);
        this->AddChild(i, child);
      }
    }
  }
}

void IndexNode::UpdateNode(uint64_t version, int index,
                           const string &child_hash, uint8_t location_in_page,
                           DeltaPage *deltapage) {
  version_ = version;
  bitmap_ |= (1 << index);
  get<0>(children_[index]) = version;
  get<1>(children_[index]) = child_hash;

  string concatenated_hash;
  for (int i = 0; i < DMM_NODE_FANOUT; i++) {
    concatenated_hash += get<1>(children_[i]);
  }
  hash_ = HashFunction(concatenated_hash);
  if (deltapage != nullptr) {
    deltapage->AddIndexNodeUpdate(location_in_page, version, hash_, index,
                                  child_hash);
  }
}

void IndexNode::AddChild(int index, Node *child, uint64_t version,
                         const string &hash) {
  if (index >= 0 && index < DMM_NODE_FANOUT) {
    children_[index] = make_tuple(version, hash, child);
    bitmap_ |= (1 << index);  // update bitmap
  } else
    throw runtime_error("AddChild out of range.");
}

Node *IndexNode::GetChild(int index) const {
  if (index >= 0 && index < DMM_NODE_FANOUT) {
    if (bitmap_ & (1 << index)) {
      return get<2>(children_[index]);
    } else
      throw runtime_error("GetChild: child doesn't exist");
  } else
    throw runtime_error("GetChild out of range.");
}

bool IndexNode::HasChild(int index) const {
  return bitmap_ & (1 << index) ? true : false;
}

void IndexNode::SetChild(int index, uint64_t version, string hash) {
  if (index >= 0 && index < DMM_NODE_FANOUT) {
    get<0>(children_[index]) = version;
    get<1>(children_[index]) = hash;
    bitmap_ |= (1 << index);  // update bitmap
  } else
    throw runtime_error("SetChild out of range.");
}

string IndexNode::GetChildHash(int index) { return get<1>(children_[index]); }
uint64_t IndexNode::GetChildVersion(int index) {
  return get<0>(children_[index]);
}

string IndexNode::GetHash() { return hash_; }
uint64_t IndexNode::GetVersion() { return version_; }
void IndexNode::SetVersion(uint64_t version) { version_ = version; }
void IndexNode::SetHash(string hash) { hash_ = hash; }

bool IndexNode::IsLeaf() const { return is_leaf_; }

NodeProof IndexNode::GetNodeProof(int level, int index) {
  NodeProof node_proof = {level, index, bitmap_};
  for (int i = 0; i < DMM_NODE_FANOUT; i++) {
    node_proof.sibling_hash.push_back(GetChildHash(i));
  }
  return node_proof;
}

bool DeltaPage::Deserialize(std::ifstream &in) {
  if (!in.good()) return false;

  try {
    // Read last_pagekey_ info
    in.read(reinterpret_cast<char *>(&last_pagekey_.version), sizeof(uint64_t));
    in.read(reinterpret_cast<char *>(&last_pagekey_.tid), sizeof(uint64_t));
    in.read(reinterpret_cast<char *>(&last_pagekey_.type), sizeof(bool));

    // Read pid size and pid
    size_t pid_size;
    in.read(reinterpret_cast<char *>(&pid_size), sizeof(pid_size));
    std::vector<char> pid_buffer(pid_size);
    in.read(pid_buffer.data(), pid_size);
    last_pagekey_.pid = string(pid_buffer.data(), pid_size);
    // Read update_count_
    in.read(reinterpret_cast<char *>(&update_count_), sizeof(uint16_t));
    // Read number of DeltaItems
    size_t items_count;
    in.read(reinterpret_cast<char *>(&items_count), sizeof(items_count));
    // Clear existing deltaitems_
    deltaitems_.clear();
    // Read each DeltaItem
    for (size_t i = 0; i < items_count; ++i) {
      DeltaItem item;
      if (!item.Deserialize(in)) {
        return false;
      }
      deltaitems_.push_back(item);
    }

    return in.good();
  } catch (const std::exception &e) {
    return false;
  }
}

bool DeltaPage::Deserialize(char *buffer) {
  // #ifdef DEBUG
  //   cout << "new DeltaPage" << endl;
  // #endif

  // Page({0, 0, true, ""});  // 临时初始化，后面会更新
  ClearDeltaPage();
  b_update_count_ = 0;

  size_t current_size = 0;
  last_pagekey_.version =
      *(reinterpret_cast<uint64_t *>(buffer + current_size));
  current_size += sizeof(uint64_t);
  last_pagekey_.tid = *(reinterpret_cast<uint64_t *>(buffer + current_size));
  current_size += sizeof(uint64_t);
  last_pagekey_.type = *(reinterpret_cast<bool *>(buffer + current_size));
  current_size += sizeof(bool);
  size_t pid_size = *(reinterpret_cast<size_t *>(buffer + current_size));
  current_size += sizeof(pid_size);
  last_pagekey_.pid = string(buffer + current_size,
                             pid_size);  // deserialize pid (pid_size bytes)
  current_size += pid_size;
  update_count_ = *(reinterpret_cast<uint16_t *>(buffer + current_size));
  current_size += sizeof(uint16_t);
  for (int i = 0; i < update_count_; i++) {
    deltaitems_.push_back(DeltaItem(buffer, current_size));
  }

  // 反序列化完成后，更新 PageKey
  PageKey pagekey = {last_pagekey_.version, last_pagekey_.tid, true,
                     last_pagekey_.pid};
  this->SetPageKey(pagekey);

  return true;
}

void DeltaPage::SerializeTo(std::ofstream &out) const {
  // 写入 last_pagekey_ 信息
  out.write(reinterpret_cast<const char *>(&last_pagekey_.version),
            sizeof(uint64_t));
  out.write(reinterpret_cast<const char *>(&last_pagekey_.tid),
            sizeof(uint64_t));
  out.write(reinterpret_cast<const char *>(&last_pagekey_.type), sizeof(bool));
  size_t pid_size = last_pagekey_.pid.size();
  out.write(reinterpret_cast<const char *>(&pid_size), sizeof(pid_size));
  out.write(last_pagekey_.pid.c_str(), pid_size);
  // 写入 update_count_
  out.write(reinterpret_cast<const char *>(&update_count_), sizeof(uint16_t));
  // 写入实际的 deltaitems_ 数量
  size_t items_count = deltaitems_.size();
  out.write(reinterpret_cast<const char *>(&items_count), sizeof(items_count));
  // 序列化每个 DeltaItem
  for (const auto &item : deltaitems_) {
    item.SerializeTo(out);
  }
}

DeltaPage::DeltaItem::DeltaItem(uint8_t loc, bool leaf, uint64_t ver,
                                const string &h, uint64_t fID, uint64_t off,
                                uint64_t sz, uint8_t idx, const string &ch_hash)
    : location_in_page(loc),
      is_leaf_node(leaf),
      version(ver),
      hash(h),
      fileID(fID),
      offset(off),
      size(sz),
      index(idx),
      child_hash(ch_hash) {
  if (index >= DMM_NODE_FANOUT) {
    throw runtime_error("index out of range");
  }
}

DeltaPage::DeltaItem::DeltaItem(char *buffer, size_t &current_size) {
  location_in_page = *(reinterpret_cast<uint8_t *>(buffer + current_size));
  current_size += sizeof(uint8_t);
  is_leaf_node = *(reinterpret_cast<bool *>(buffer + current_size));
  current_size += sizeof(bool);
  version = *(reinterpret_cast<uint64_t *>(buffer + current_size));
  current_size += sizeof(uint64_t);
  hash = string(buffer + current_size, HASH_SIZE);
  current_size += HASH_SIZE;

  if (is_leaf_node) {
    fileID = *(reinterpret_cast<uint64_t *>(buffer + current_size));
    current_size += sizeof(uint64_t);
    offset = *(reinterpret_cast<uint64_t *>(buffer + current_size));
    current_size += sizeof(uint64_t);
    size = *(reinterpret_cast<uint64_t *>(buffer + current_size));
    current_size += sizeof(uint64_t);
  } else {
    index = *(reinterpret_cast<uint8_t *>(buffer + current_size));
    if (index >= DMM_NODE_FANOUT) {
      throw runtime_error("index out of range");
    }
    current_size += sizeof(uint8_t);
    child_hash = string(buffer + current_size, HASH_SIZE);
    current_size += HASH_SIZE;
  }
}

// ... existing code ...

void DeltaPage::DeltaItem::SerializeTo(char *buffer,
                                       size_t &current_size) const {
  // Write location_in_page
  memcpy(buffer + current_size, &location_in_page, sizeof(uint8_t));
  current_size += sizeof(uint8_t);

  // Write is_leaf_node
  memcpy(buffer + current_size, &is_leaf_node, sizeof(bool));
  current_size += sizeof(bool);

  // Write version
  memcpy(buffer + current_size, &version, sizeof(uint64_t));
  current_size += sizeof(uint64_t);

  // Write hash length and hash
  // uint32_t hash_length = hash.length();
  // memcpy(buffer + current_size, &hash_length, sizeof(uint32_t));
  // current_size += sizeof(uint32_t);
  memcpy(buffer + current_size, hash.c_str(), HASH_SIZE);
  current_size += HASH_SIZE;

  if (is_leaf_node) {
    // Write leaf node specific data
    memcpy(buffer + current_size, &fileID, sizeof(uint64_t));
    current_size += sizeof(uint64_t);
    memcpy(buffer + current_size, &offset, sizeof(uint64_t));
    current_size += sizeof(uint64_t);
    memcpy(buffer + current_size, &size, sizeof(uint64_t));
    current_size += sizeof(uint64_t);
  } else {
    // Write index node specific data
    memcpy(buffer + current_size, &index, sizeof(uint8_t));
    if (index >= DMM_NODE_FANOUT) {
      throw runtime_error("index out of range");
    }
    current_size += sizeof(uint8_t);
    // Write child_hash length and child_hash
    memcpy(buffer + current_size, child_hash.c_str(), HASH_SIZE);
    current_size += HASH_SIZE;
  }
}

void DeltaPage::DeltaItem::SerializeTo(std::ofstream &out) const {
  out.write(reinterpret_cast<const char *>(&location_in_page),
            sizeof(location_in_page));
  out.write(reinterpret_cast<const char *>(&is_leaf_node),
            sizeof(is_leaf_node));
  out.write(reinterpret_cast<const char *>(&version), sizeof(version));
  out.write(hash.c_str(), HASH_SIZE);

  if (is_leaf_node) {
    out.write(reinterpret_cast<const char *>(&fileID), sizeof(fileID));
    out.write(reinterpret_cast<const char *>(&offset), sizeof(offset));
    out.write(reinterpret_cast<const char *>(&size), sizeof(size));
  } else {
    out.write(reinterpret_cast<const char *>(&index), sizeof(index));
    if (index >= DMM_NODE_FANOUT) {
      throw runtime_error("index out of range");
    }
    out.write(child_hash.c_str(), HASH_SIZE);
  }
}

// ... existing code ...

bool DeltaPage::DeltaItem::Deserialize(std::ifstream &in) {
  if (!in.good()) return false;

  try {
    // Read basic fields
    in.read(reinterpret_cast<char *>(&location_in_page), sizeof(uint8_t));
    in.read(reinterpret_cast<char *>(&is_leaf_node), sizeof(bool));
    in.read(reinterpret_cast<char *>(&version), sizeof(uint64_t));

    // Read hash
    char hash_buffer[HASH_SIZE];
    in.read(hash_buffer, HASH_SIZE);
    hash = string(hash_buffer, HASH_SIZE);

    if (is_leaf_node) {
      // Read leaf node specific fields
      in.read(reinterpret_cast<char *>(&fileID), sizeof(uint64_t));
      in.read(reinterpret_cast<char *>(&offset), sizeof(uint64_t));
      in.read(reinterpret_cast<char *>(&size), sizeof(uint64_t));

      // Initialize unused index node fields
      index = 0;
      child_hash = "";
    } else {
      // Read index node specific fields

      in.read(reinterpret_cast<char *>(&index), sizeof(uint8_t));
      if (index >= DMM_NODE_FANOUT) {
        throw runtime_error("index out of range");
      }
      char child_hash_buffer[HASH_SIZE];
      in.read(child_hash_buffer, HASH_SIZE);
      child_hash = string(child_hash_buffer, HASH_SIZE);

      // Initialize unused leaf node fields
      fileID = 0;
      offset = 0;
      size = 0;
    }

    return in.good();
  } catch (const std::exception &e) {
    return false;
  }
}
// ... existing code ...

DeltaPage::DeltaPage(PageKey last_pagekey, uint16_t update_count,
                     uint16_t b_update_count)
    : last_pagekey_(last_pagekey),
      update_count_(update_count),
      b_update_count_(b_update_count) {
  // #ifdef DEBUG
  //   cout << "new DeltaPage" << endl;
  // #endif
}

DeltaPage::DeltaPage(const DeltaPage &other) : Page(other) {
  // #ifdef DEBUG
  //   cout << "new DeltaPage" << endl;
  // #endif

  // Copy DeltaPage specific members
  last_pagekey_ = other.last_pagekey_;
  update_count_ = other.update_count_;
  b_update_count_ = other.b_update_count_;
  deltaitems_ = other.deltaitems_;
}

DeltaPage::DeltaPage(char *buffer) : b_update_count_(0) {
  // #ifdef DEBUG
  //   cout << "new DeltaPage" << endl;
  // #endif

  Page({0, 0, true, ""});  // 临时初始化，后面会更新

  size_t current_size = 0;

  last_pagekey_.version =
      *(reinterpret_cast<uint64_t *>(buffer + current_size));
  current_size += sizeof(uint64_t);
  last_pagekey_.tid = *(reinterpret_cast<uint64_t *>(buffer + current_size));
  current_size += sizeof(uint64_t);
  last_pagekey_.type = *(reinterpret_cast<bool *>(buffer + current_size));
  current_size += sizeof(bool);
  size_t pid_size = *(reinterpret_cast<size_t *>(buffer + current_size));
  current_size += sizeof(pid_size);
  last_pagekey_.pid = string(buffer + current_size,
                             pid_size);  // deserialize pid (pid_size bytes)
  current_size += pid_size;
  update_count_ = *(reinterpret_cast<uint16_t *>(buffer + current_size));
  current_size += sizeof(uint16_t);
  for (int i = 0; i < update_count_; i++) {
    deltaitems_.push_back(DeltaItem(buffer, current_size));
  }

  // 反序列化完成后，更新 PageKey
  PageKey pagekey = {last_pagekey_.version, last_pagekey_.tid, true,
                     last_pagekey_.pid};
  this->SetPageKey(pagekey);
}

DeltaPage::~DeltaPage() {
  // #ifdef DEBUG
  //   cout << "delete DeltaPage" << endl;
  // #endif
}

void DeltaPage::AddIndexNodeUpdate(uint8_t location, uint64_t version,
                                   const string &hash, uint8_t index,
                                   const string &child_hash) {
  deltaitems_.push_back(
      DeltaItem(location, false, version, hash, 0, 0, 0, index, child_hash));
  ++update_count_;
  ++b_update_count_;
}

void DeltaPage::AddLeafNodeUpdate(uint8_t location, uint64_t version,
                                  const string &hash, uint64_t fileID,
                                  uint64_t offset, uint64_t size) {
  deltaitems_.push_back(
      DeltaItem(location, true, version, hash, fileID, offset, size));
  ++update_count_;
  ++b_update_count_;
}

void DeltaPage::SerializeTo() {
  char *buffer = this->GetData();
  memset(buffer, 0, PAGE_SIZE);
  size_t current_size = 0;
  memcpy(buffer + current_size, &last_pagekey_.version, sizeof(uint64_t));
  current_size += sizeof(uint64_t);
  memcpy(buffer + current_size, &last_pagekey_.tid, sizeof(uint64_t));
  current_size += sizeof(uint64_t);
  memcpy(buffer + current_size, &last_pagekey_.type, sizeof(bool));
  current_size += sizeof(bool);
  size_t pid_size = last_pagekey_.pid.size();
  memcpy(buffer + current_size, &pid_size, sizeof(pid_size));
  current_size += sizeof(pid_size);
  memcpy(buffer + current_size, last_pagekey_.pid.c_str(), pid_size);
  current_size += pid_size;

  memcpy(buffer + current_size, &update_count_, sizeof(uint16_t));
  current_size += sizeof(uint16_t);

  for (const auto &item : deltaitems_) {
    if (current_size + sizeof(DeltaItem) > PAGE_SIZE) {  // exceeds page size
      throw overflow_error(
          "DeltaPage exceeds PAGE_SIZE during serialization._");
    }
    item.SerializeTo(buffer, current_size);
  }
}

void DeltaPage::ClearDeltaPage() {
  deltaitems_.clear();
  update_count_ = 0;
}

const vector<DeltaPage::DeltaItem> &DeltaPage::GetDeltaItems() const {
  return deltaitems_;
}

PageKey DeltaPage::GetLastPageKey() const { return last_pagekey_; }

void DeltaPage::SetLastPageKey(PageKey pagekey) { last_pagekey_ = pagekey; }

uint16_t DeltaPage::GetDeltaPageUpdateCount() { return update_count_; }

uint16_t DeltaPage::GetBasePageUpdateCount() { return b_update_count_; }

void DeltaPage::ClearBasePageUpdateCount() { b_update_count_ = 0; }

BasePage::BasePage(DMMTrie *trie, Node *root, const string &pid)
    : trie_(trie), root_(root), Page({0, 0, false, pid}) {
  // #ifdef DEBUG
  //   cout << "new BasePage" << endl;
  // #endif
}

BasePage::BasePage(const BasePage &other) : Page(other), trie_(other.trie_) {
  // #ifdef DEBUG
  //   cout << "new BasePage" << endl;
  // #endif
  // Deep copy the root node
  if (other.root_) {
    if (other.root_->IsLeaf()) {
      root_ = new LeafNode(*dynamic_cast<LeafNode *>(other.root_));
    } else {
      root_ = new IndexNode(*dynamic_cast<IndexNode *>(other.root_));
    }
  } else {
    root_ = nullptr;
  }
}

BasePage::BasePage(DMMTrie *trie, char *buffer) : trie_(trie) {
  Page({0, 0, false, ""});  // 临时初始化，后面会更新
                            // #ifdef DEBUG
                            //   cout << "new BasePage" << endl;
                            // #endif
  size_t current_size = 0;

  uint64_t version = *(reinterpret_cast<uint64_t *>(
      buffer + current_size));  // deserialize version
  current_size += sizeof(uint64_t);

  uint64_t tid = *(reinterpret_cast<uint64_t *>(
      buffer + current_size));  // deserialize DMMTrie id
  current_size += sizeof(uint64_t);

  bool page_type = *(reinterpret_cast<bool *>(
      buffer + current_size));  // deserialize page type (1 byte)
  current_size += sizeof(bool);

  size_t pid_size = *(reinterpret_cast<size_t *>(
      buffer + current_size));  // deserialize pid_size (8 bytes for size_t)
  current_size += sizeof(pid_size);
  string pid(buffer + current_size,
             pid_size);  // deserialize pid (pid_size bytes)
  current_size += pid_size;

  bool is_leaf_node = *(reinterpret_cast<bool *>(buffer + current_size));
  current_size += sizeof(bool);

  if (is_leaf_node) {  // the root node of page is leafnode
    root_ = new LeafNode();
    root_->DeserializeFrom(buffer, current_size, true);
  } else {  // the root node of page is indexnode
    root_ = new IndexNode();
    root_->DeserializeFrom(buffer, current_size, true);
  }

  // 反序列化完成后，更新 PageKey
  PageKey pagekey = {version, tid, page_type, pid};
  this->SetPageKey(pagekey);
}

BasePage::BasePage(DMMTrie *trie, string key, string pid, string nibbles)
    : trie_(trie), Page({0, 0, false, pid}) {
  // #ifdef DEBUG
  //   cout << "new BasePage" << endl;
  // #endif
  if (nibbles.size() == 0) {  // leafnode
    root_ = new LeafNode(0, key, {}, "");
  } else if (nibbles.size() == 1) {  // indexnode->leafnode
    Node *child_node = new LeafNode(0, key, {}, "");
    root_ = new IndexNode(0, "", 0);

    int index = GetIndex(nibbles[0]);
    root_->AddChild(index, child_node, 0, "");
  } else {  // indexnode->indexnode
    int index = GetIndex(nibbles[1]);
    // second level of indexnode should route its child by bitmap
    Node *child_node = new IndexNode(0, "", 1 << index);
    root_ = new IndexNode(0, "", 0);

    index = GetIndex(nibbles[0]);
    root_->AddChild(index, child_node, 0, "");
  }
}

BasePage::~BasePage() {
  // #ifdef DEBUG
  //   cout << "delete BasePage" << endl;
  // #endif
  for (int i = 0; i < DMM_NODE_FANOUT; i++) {
    if (root_->HasChild(i)) {
      delete root_->GetChild(i);
    }
  }
  delete root_;
}

/* serialized BasePage format (size in bytes):
   | version (8) | tid (8) | tp (1) | pid_size (8 in 64-bit system) | pid
   (pid_size) | root node | */
void BasePage::SerializeTo() {
  char *buffer = this->GetData();
  size_t current_size = 0;

  uint64_t version = root_->GetVersion();
  memcpy(buffer + current_size, &version, sizeof(uint64_t));
  current_size += sizeof(uint64_t);

  uint64_t tid = 0;
  memcpy(buffer + current_size, &tid, sizeof(uint64_t));
  current_size += sizeof(uint64_t);

  bool page_type = false;  // Tp is false means basepage
  memcpy(buffer + current_size, &page_type, sizeof(bool));
  current_size += sizeof(bool);

  size_t pid_size = GetPageKey().pid.size();
  memcpy(buffer + current_size, &pid_size, sizeof(pid_size));  // pid size
  current_size += sizeof(pid_size);
  memcpy(buffer + current_size, GetPageKey().pid.c_str(), pid_size);  // pid
  current_size += pid_size;

  root_->SerializeTo(buffer, current_size, true);  // serialize nodes
}

void BasePage::UpdatePage(uint64_t version,
                          tuple<uint64_t, uint64_t, uint64_t> location,
                          const string &value, const string &nibbles,
                          const string &child_hash, DeltaPage *deltapage,
                          PageKey pagekey) {
  // parameter "nibbles" are the first two nibbles after pid
  if (nibbles.size() == 0) {
    // page has one leafnode, eg. page "abcdef" for key "abcdef"
    if (!root_) {
      root_ = new LeafNode(0, pagekey.pid, {}, "");
    }
    static_cast<LeafNode *>(root_)->UpdateNode(version, location, value, 0,
                                               deltapage);
  } else if (nibbles.size() == 1) {
    // page has one indexnode and one level of leafnodes, eg. page "abcd" for
    // key "abcde"
    if (!root_) {
      root_ = new IndexNode(0, "", 0);
    }
    int index = GetIndex(nibbles[0]);
    if (!root_->HasChild(index)) {
      Node *child_node =
          new LeafNode(0, pagekey.pid + to_string(index), {}, "");
      root_->AddChild(index, child_node, 0, "");
    }
    static_cast<LeafNode *>(root_->GetChild(index))
        ->UpdateNode(version, location, value, index + 1, deltapage);

    string child_hash_2 = root_->GetChild(index)->GetHash();
    static_cast<IndexNode *>(root_)->UpdateNode(version, index, child_hash_2, 0,
                                                deltapage);
  } else {
    // page has two levels of indexnodes , eg. page "ab" for key "abcdef"
    if (!root_) {
      root_ = new IndexNode(0, "", 0);
    }
    int index = GetIndex(nibbles[0]), child_index = GetIndex(nibbles[1]);
    if (!root_->HasChild(index)) {
      Node *child_node = new IndexNode(0, "", 1 << child_index);
      root_->AddChild(index, child_node, 0, "");
    }
    static_cast<IndexNode *>(root_->GetChild(index))
        ->UpdateNode(version, child_index, child_hash, index + 1, deltapage);

    string child_hash_2 = root_->GetChild(index)->GetHash();
    static_cast<IndexNode *>(root_)->UpdateNode(version, index, child_hash_2, 0,
                                                deltapage);
  }

  this->SetPageKey(pagekey);
  if (deltapage != nullptr) {
    PageKey deltapage_pagekey = {version, 0, true, pagekey.pid};
    deltapage->SetPageKey(deltapage_pagekey);

    if (deltapage->GetDeltaPageUpdateCount() >= Td_) {
      // When a DeltaPage accumulates 𝑇𝑑 updates, it is frozen and a new active
      // one is initiated

      DeltaPage *deltapage_copy = new DeltaPage(*deltapage);
      deltapage_copy->SerializeTo();
      // store frozen deltapage in cache
      trie_->WritePageCache(deltapage_pagekey, deltapage_copy);

      deltapage->ClearDeltaPage();  // delete all DeltaItems in DeltaPage
      // record the PageKey of DeltaPage passed to LSVPS
      deltapage->SetLastPageKey(deltapage_pagekey);
      trie_->AddDeltaPageVersion(pagekey.pid, version);
    }
    if (deltapage->GetBasePageUpdateCount() >= Tb_) {
      // Each page generates a checkpoint as BasePage after every 𝑇𝑏 updates
      BasePage *basepage_copy =
          new BasePage(*this);  // not deep copy!!!!!!!!!!!!!!!!!!!!!!!!!
      basepage_copy->SerializeTo();
      trie_->WritePageCache(pagekey, basepage_copy);  // store basepage in cache

      trie_->UpdatePageVersion(pagekey, version, version);
      deltapage->ClearBasePageUpdateCount();
      deltapage->SetLastPageKey(pagekey);
      return;
    }
  }

  pair<uint64_t, uint64_t> page_version = trie_->GetPageVersion(pagekey);
  trie_->UpdatePageVersion(pagekey, version, page_version.second);
}

void BasePage::UpdateDeltaItem(const DeltaPage::DeltaItem &deltaitem) {
  // add one update from deltapage to basepage
  Node *node = nullptr;
  if (root_ == nullptr) {
    // create root if replay function in LSVPS has no basepage to start from
    if (deltaitem.location_in_page == 0 && deltaitem.is_leaf_node == true) {
      root_ = new LeafNode();
    } else {
      root_ = new IndexNode();
    }
  }

  if (deltaitem.is_leaf_node) {
    if (deltaitem.location_in_page == 0) {
      node = root_;
    } else if (!root_->HasChild(deltaitem.location_in_page - 1)) {
      node = new LeafNode();
      root_->AddChild(deltaitem.location_in_page - 1, node, 0, "");
    } else {
      node = root_->GetChild(deltaitem.location_in_page - 1);
    }

    node->SetVersion(deltaitem.version);
    node->SetLocation(
        make_tuple(deltaitem.fileID, deltaitem.offset, deltaitem.size));
    node->SetHash(deltaitem.hash);
  } else {
    if (deltaitem.location_in_page == 0) {
      node = root_;
    } else if (!root_->HasChild(deltaitem.location_in_page - 1)) {
      node = new IndexNode();
      root_->AddChild(deltaitem.location_in_page - 1, node, 0, "");
    } else {
      node = root_->GetChild(deltaitem.location_in_page - 1);
    }

    node->SetVersion(deltaitem.version);
    node->SetHash(deltaitem.hash);
    node->SetChild(deltaitem.index, deltaitem.version, deltaitem.child_hash);
  }
  const auto &old_pagekey = GetPageKey();
  SetPageKey(PageKey{deltaitem.version, old_pagekey.tid, old_pagekey.type,
                     old_pagekey.pid});
}

Node *BasePage::GetRoot() const { return root_; }

DMMTrie::DMMTrie(uint64_t tid, LSVPS *page_store, VDLS *value_store,
                 uint64_t current_version)
    : tid(tid),
      page_store_(page_store),
      value_store_(value_store),
      current_version_(current_version),
      root_page_(nullptr) {
  lru_cache_.clear();
  pagekeys_.clear();
  active_deltapages_.clear();
  page_versions_.clear();
  page_cache_.clear();
  put_cache_.clear();
  deltapage_versions_.clear();
}

DMMTrie::~DMMTrie() {
  while (lru_cache_.size()) {  // cache is full
    PageKey last_key = pagekeys_.back().first;
    auto last_iter = lru_cache_.find(last_key);
    delete last_iter->second->second;  // release memory of basepage

    // remove the page whose pagekey is at the tail of list
    lru_cache_.erase(last_key);
    pagekeys_.pop_back();
  }
}

bool DMMTrie::Put(uint64_t tid, uint64_t version, const string &key,
                  const string &value) {
  if (version < current_version_) {  // version invalid
    cout << "Version " << version << " is outdated!" << endl;
    return false;
  }
  if (value == "") {
    cout << "Value cannot be empty string" << endl;
    return false;
  }
  current_version_ = version;
  put_cache_[key] = value;
  return true;
}

string DMMTrie::Get(uint64_t tid, uint64_t version, const string &key) {
  string nibble_path = key;
  uint64_t page_version = version;
  LeafNode *leafnode = nullptr;
  for (int i = 0; i <= key.size(); i += 2) {
    string pid = nibble_path.substr(0, i);
    BasePage *page =
        GetPage({page_version, 0, false, pid});  // false means basepage
    if (page == nullptr || page->GetRoot() == nullptr) {
      cout << "Key " << key << " not found at version " << version << endl;
      return "";
    }

    if (!page->GetRoot()->IsLeaf()) {  // first level in page is indexnode

      if (!page->GetRoot()->HasChild(GetIndex(nibble_path[i]))) {
        cout << "Child not found" << endl;
        cout << "Key " << key << " not found at version " << version << endl;
        return "";
      }
      if (!page->GetRoot()->GetChild(GetIndex(nibble_path[i]))->IsLeaf()) {
        // second level is indexnode
        // TODO: child的版本比Root高是正常的吗？
        page_version = page->GetRoot()
                           ->GetChild(GetIndex(nibble_path[i]))
                           ->GetChildVersion(GetIndex(nibble_path[i + 1]));
      } else {  // second level is leafnode
        leafnode = static_cast<LeafNode *>(
            page->GetRoot()->GetChild(GetIndex(nibble_path[i])));
      }
    } else {  // first level is leafnode
      leafnode = static_cast<LeafNode *>(page->GetRoot());
    }
  }
  tuple<uint64_t, uint64_t, uint64_t> location = leafnode->GetLocation();
#ifdef DEBUG
  cout << "location:" << get<0>(location) << " " << get<1>(location) << " "
       << get<2>(location) << endl;
#endif
  string value = value_store_->ReadValue(leafnode->GetLocation());
#ifdef DEBUG
  cout << "Key " << key << " has value " << value << " at version " << version
       << endl;
#endif
  return value;
}

void DMMTrie::Delete(uint64_t tid, uint64_t version, const string &key) {
  if (version < current_version_) {  // version invalid
    cout << "Version " << version << " is outdated!" << endl;
    return;
  }
  current_version_ = version;
  put_cache_[key] = "";
}

// deprecated
void DMMTrie::Commit(uint64_t version) { CalcRootHash(0, version); }

void DMMTrie::CalcRootHash(uint64_t tid, uint64_t version) {
  if (version != current_version_) {
    cout << "Commit version incompatible" << endl;
  }

  map<string, set<string>, decltype(CompareStrings)> updates(CompareStrings);

  for (const auto &it : put_cache_) {
    for (int i = it.first.size() % 2 == 0 ? it.first.size()
                                          : it.first.size() - 1;
         i >= 0; i -= 2) {
      // store the pid and nibbles of each page updated in every put
      updates[it.first.substr(0, i)].insert(it.first.substr(i, 2));
    }
  }

  // unordered_map<string, DeltaPage *> active_deltapages;
  set<string> pids;

  for (const auto &it : put_cache_) {
    for (int i = it.first.size() % 2 == 0 ? it.first.size()
                                          : it.first.size() - 1;
         i >= 0; i -= 2) {
      pids.insert(it.first.substr(0, i));
    }
  }

  // get the needed active deltapages from LSVPS
  // for (string pid : pids) {
  //   active_deltapages[pid] = page_store_->GetActiveDeltaPage(pid);
  // }

  for (const auto &it : updates) {
    string pid = it.first;
    bool if_exceed = false;
    // get the latest version number of a page
    uint64_t page_version = GetPageVersion({0, 0, false, pid}).first;
    PageKey pagekey = {version, 0, false, pid},
            old_pagekey = {page_version, 0, false, pid};
    BasePage *page = GetPage(old_pagekey);  // load the page into lru cache

    if (page == nullptr) {
      // GetPage returns nullptr means that the pid is new
      page = new BasePage(this, nullptr, pid);
      PutPage(pagekey, page);  // add the newly generated page into cache
    }

    // DeltaPage *deltapage = GetDeltaPage(pid);

    DeltaPage *deltapage = page_store_->GetActiveDeltaPage(pid);

    if (2 * it.second.size() + deltapage->GetDeltaPageUpdateCount() >=
        2 * Td_) {
      // the updates in page is more than the capacity of two deltapages
      if_exceed = true;
      if (deltapage->GetDeltaPageUpdateCount() != 0) {
        PageKey deltapage_pagekey = {version, 0, true, pagekey.pid};

        DeltaPage *deltapage_copy = new DeltaPage(*deltapage);
        deltapage_copy->SetPageKey(deltapage_pagekey);
        deltapage_copy->SerializeTo();
        // store frozen deltapage in cache
        WritePageCache(deltapage_pagekey, deltapage_copy);

        deltapage->ClearDeltaPage();  // delete all DeltaItems in DeltaPage
        // record the PageKey of DeltaPage passed to LSVPS
        deltapage->SetLastPageKey(deltapage_pagekey);
        AddDeltaPageVersion(pagekey.pid, version);
      }
    }

    for (const auto &nibbles : it.second) {
      // path is key when page is leaf page, pid of child page when page is
      // index page
      string path = pid + nibbles;
      tuple<uint64_t, uint64_t, uint64_t> location;
      string value, child_hash;
      if (nibbles.size() == 2) {  // indexnode + indexnode
        child_hash = GetPage({version, 0, false, path})->GetRoot()->GetHash();
      } else {  // (indexnode + leafnode) or leafnode
        value = put_cache_[path];
        location = value_store_->WriteValue(version, path, value);
      }
      if (if_exceed) {
        page->UpdatePage(version, location, value, nibbles, child_hash, nullptr,
                         pagekey);
      } else {
        page->UpdatePage(version, location, value, nibbles, child_hash,
                         deltapage, pagekey);
      }
    }

    if (if_exceed) {
      BasePage *basepage_copy = new BasePage(*page);
      basepage_copy->SerializeTo();
      WritePageCache(pagekey, basepage_copy);  // store basepage in cache

      UpdatePageVersion(pagekey, version, version);
      deltapage->ClearBasePageUpdateCount();
      deltapage->SetLastPageKey(pagekey);
    }
    UpdatePageKey(old_pagekey, pagekey);
    // deltapage->SerializeTo();
    page_store_->StoreActiveDeltaPage(deltapage);
  }

  for (const auto &it : page_cache_) {
    page_store_->StorePage(it.second);
#ifdef DEBUG
    std::cout << "Commit" << version
              << " Store Page: " << it.second->GetPageKey() << std::endl;
#endif
  }

  // send the active deltapages back to LSVPS
  // for (const auto &it : active_deltapages) {
  //   page_store_->StoreActiveDeltaPage(it.second);
  // }
  for (auto &pair : page_cache_) {
    delete pair.second;
  }
  page_cache_.clear();
  put_cache_.clear();
#ifdef DEBUG
  cout << "Version " << version << " committed" << endl;
  cout << "Active delta pages: " << active_deltapages_.size() << endl;
  cout << "Active delta page size: " << sizeof(active_deltapages_.end()->second)
       << endl;
  cout << "LRU pages:" << lru_cache_.size() << endl;
  cout << "page_cache_:" << page_cache_.size() << endl;
  cout << "pagekeys_:" << pagekeys_.size() << endl;

  std::ifstream file("/proc/self/status");
  std::string line;
  while (std::getline(file, line)) {
    std::istringstream iss(line);
    std::string key;
    int value;
    if (iss >> key >> value) {
      if (key == "VmSize:")
        std::cout << "Virtual memory used: " << value << " kB" << endl;
      else if (key == "VmRSS:")
        std::cout << "Physical memory used: " << value << " kB" << endl;
    }
  }
#endif
}

string DMMTrie::GetRootHash(uint64_t tid, uint64_t version) {
  return GetPage({version, tid, false, ""})->GetRoot()->GetHash();
}

DMMTrieProof DMMTrie::GetProof(uint64_t tid, uint64_t version,
                               const string &key) {
  DMMTrieProof merkle_proof;
  string nibble_path = key;
  uint64_t page_version = version;
  LeafNode *leafnode = nullptr;
  for (int i = 0; i < key.size() + 1; i += 2) {
    string pid = nibble_path.substr(0, i);
    BasePage *page =
        GetPage({page_version, 0, false, pid});  // false means basepage
    if (page == nullptr || page->GetRoot() == nullptr) {
      cout << "Key " << key << " not found at version " << version << endl;
      merkle_proof.value = "";
      return merkle_proof;
    }

    if (!page->GetRoot()->IsLeaf()) {

      if (!page->GetRoot()->HasChild(GetIndex(nibble_path[i]))) {

        cout << "Key " << key << " not found at version " << version << endl;
        merkle_proof.value = "";
        return merkle_proof;
      }
      // first level in page is indexnode
      merkle_proof.proofs.push_back(
          page->GetRoot()->GetNodeProof(i, GetIndex(nibble_path[i])));
      if (!page->GetRoot()->GetChild(GetIndex(nibble_path[i]))->IsLeaf()) {
        // second level is indexnode
        merkle_proof.proofs.push_back(
            page->GetRoot()
                ->GetChild(GetIndex(nibble_path[i]))
                ->GetNodeProof(i + 1, GetIndex(nibble_path[i + 1])));
        page_version = page->GetRoot()
                           ->GetChild(GetIndex(nibble_path[i]))
                           ->GetChildVersion(GetIndex(nibble_path[i + 1]));
      } else {  // second level is leafnode
        leafnode = static_cast<LeafNode *>(
            page->GetRoot()->GetChild(GetIndex(nibble_path[i])));
      }
    } else {  // first level is leafnode
      leafnode = static_cast<LeafNode *>(page->GetRoot());
    }
  }
  merkle_proof.value = value_store_->ReadValue(leafnode->GetLocation());
  reverse(merkle_proof.proofs.begin(), merkle_proof.proofs.end());
  return merkle_proof;
}

bool DMMTrie::Verify(uint64_t tid, const string &key, const string &value,
                     string root_hash, DMMTrieProof proof) {
  // string hash = HashFunction(key + value);
  string hash = HashFunction(value);
  for (const auto &node_proof : proof.proofs) {
    string concatenated_hash;
    for (int i = 0; i < DMM_NODE_FANOUT; i++) {
      if (i == node_proof.index) {
        concatenated_hash += hash;
      } else {
        concatenated_hash += node_proof.sibling_hash[i];
      }
    }
    hash = HashFunction(concatenated_hash);
  }
  return hash == root_hash;
}

bool DMMTrie::Verify(uint64_t tid, uint64_t version, string root_hash) {
  return RecursiveVerify({version, tid, false, ""}) == root_hash;
}

string DMMTrie::RecursiveVerify(PageKey pagekey) {
  BasePage *page = GetPage(pagekey);
  if (page == nullptr) {
    return "";
  }

  if (page->GetRoot()->IsLeaf()) {
    // first level is indexnode
    string value = value_store_->ReadValue(
        static_cast<LeafNode *>(page->GetRoot())->GetLocation());
    // return HashFunction(pagekey.pid + value);
    return HashFunction(value);
  }

  string concatenated_hash;
  for (int i = 0; i < DMM_NODE_FANOUT; i++) {
    if (!page->GetRoot()->HasChild(i)) {
      continue;
    }
    Node *child = page->GetRoot()->GetChild(i);
    if (!child->IsLeaf()) {
      // second level is indexnode
      string child_concatenated_hash;
      for (int j = 0; j < DMM_NODE_FANOUT; j++) {
        if (!child->HasChild(j)) {
          continue;
        }
        // call RecusiveVerify to calculate hash in child page
        child_concatenated_hash +=
            RecursiveVerify({pagekey.version, pagekey.tid, false,
                             pagekey.pid + to_string(i) + to_string(j)});
      }
      concatenated_hash += HashFunction(child_concatenated_hash);
    } else {
      string value = value_store_->ReadValue(
          static_cast<LeafNode *>(child)->GetLocation());
      // concatenated_hash += HashFunction(pagekey.pid + to_string(i) + value);
      concatenated_hash += HashFunction(value);
    }
  }
  return HashFunction(concatenated_hash);
}

void DMMTrie::Flush(uint64_t tid, uint64_t version) { page_store_->Flush(); }

void DMMTrie::Revert(uint64_t tid, uint64_t version) {}

DeltaPage *DMMTrie::GetDeltaPage(const string &pid) {
  auto it = active_deltapages_.find(pid);
  if (it != active_deltapages_.end()) {
    return &it->second;  // return deltapage if it exiests
  } else {
    DeltaPage new_page;
    new_page.SetLastPageKey(PageKey{0, 0, false, pid});
    active_deltapages_[pid] = new_page;
    return &active_deltapages_[pid];
  }
}

pair<uint64_t, uint64_t> DMMTrie::GetPageVersion(PageKey pagekey) {
  auto it = page_versions_.find(pagekey.pid);
  if (it != page_versions_.end()) {
    return it->second;
  }
  return {0, 0};
}

PageKey DMMTrie::GetLatestBasePageKey(PageKey pagekey) const {
  auto it = page_versions_.find(pagekey.pid);
  if (it != page_versions_.end()) {
    return {it->second.second, pagekey.tid, false, pagekey.pid};
  }
  return PageKey{0, 0, false, pagekey.pid};
}

void DMMTrie::UpdatePageVersion(PageKey pagekey, uint64_t current_version,
                                uint64_t latest_basepage_version) {
  page_versions_[pagekey.pid] = {current_version, latest_basepage_version};
}

void DMMTrie::WritePageCache(PageKey pagekey, Page *page) {
  page_cache_[pagekey] = page;
}

void DMMTrie::AddDeltaPageVersion(const string &pid, uint64_t version) {
  deltapage_versions_[pid].push_back(version);
}

uint64_t DMMTrie::GetVersionUpperbound(const string &pid, uint64_t version) {
  if (deltapage_versions_.find(pid) == deltapage_versions_.end()) {
    return 0;  // no deltapage of this pid
  }

  vector<uint64_t> &versions = deltapage_versions_[pid];
  auto it = upper_bound(versions.begin(), versions.end(), version);
  if (it == versions.end()) {
    // no deltapage has version larger than requested
    return current_version_;
  }
  return *it;
}

BasePage *DMMTrie::GetPage(
    const PageKey &pagekey) {  // get a page by its pagekey
  auto it = lru_cache_.find(pagekey);
  if (it != lru_cache_.end()) {  // page is in cache
    // move the accessed page to the front
    pagekeys_.splice(pagekeys_.begin(), pagekeys_, it->second);
    it->second = pagekeys_.begin();  // update iterator
    return it->second->second;
  }
  // page is not in cache, fetch it from LSVPS
  BasePage *page = page_store_->LoadPage(pagekey);
  if (!page) {  // page is not found in disk
    return nullptr;
  }
  // if (!page->GetRoot()) {  // page is not found in disk
  //   return nullptr;
  // }
  PutPage(pagekey, page);
  return page;
}

void DMMTrie::PutPage(const PageKey &pagekey,
                      BasePage *page) {        // add page to cache
  if (lru_cache_.size() >= max_cache_size_) {  // cache is full
    PageKey last_key = pagekeys_.back().first;
    auto last_iter = lru_cache_.find(last_key);
    delete last_iter->second->second;  // release memory of basepage

    // remove the page whose pagekey is at the tail of list
    lru_cache_.erase(last_key);
    pagekeys_.pop_back();
  }
  auto it = lru_cache_.find(pagekey);
  if (it != lru_cache_.end()) {
    delete it->second->second;
    pagekeys_.erase(it->second);
    lru_cache_.erase(it);
  }
  // insert the pair of PageKey and BasePage* to the front
  pagekeys_.push_front(make_pair(pagekey, page));
  lru_cache_[pagekey] = pagekeys_.begin();
}

void DMMTrie::UpdatePageKey(
    const PageKey &old_pagekey,
    const PageKey &new_pagekey) {  // update pagekey in lru cache
  auto it = lru_cache_.find(old_pagekey);
  if (it != lru_cache_.end()) {
    // save the basepage indexed by old pagekey
    BasePage *basepage = it->second->second;

    pagekeys_.erase(it->second);  // delete old pagekey item
    lru_cache_.erase(it);

    pagekeys_.push_front(make_pair(new_pagekey, basepage));
    lru_cache_[new_pagekey] = pagekeys_.begin();
  }
}