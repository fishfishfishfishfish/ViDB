# rocksdb (≥ 5.8)
检查是否已经安装
```bash
$ rocksdb_dump --version
rocksdb_dump: command not found
```

```bash
~$ git clone https://github.com/facebook/rocksdb.git
Cloning into 'rocksdb'...
remote: Enumerating objects: 138060, done.
remote: Counting objects: 100% (471/471), done.
remote: Compressing objects: 100% (296/296), done.
remote: Total 138060 (delta 326), reused 177 (delta 175), pack-reused 137589 (from 3)
Receiving objects: 100% (138060/138060), 222.78 MiB | 1.36 MiB/s, done.
Resolving deltas: 100% (105734/105734), done.
~$ cd rocksdb/
~/rocksdb$ git checkout rocksdb-5.8
HEAD is now at 266ac245a Bumping version to 5.8
~/rocksdb$ make shared_lib
  GEN      util/build_version.cc
  GEN      util/build_version.cc
  CC       shared-objects/cache/clock_cache.o
  CC       shared-objects/cache/lru_cache.o
  CC       shared-objects/cache/sharded_cache.o
  CC       shared-objects/db/builder.o
  CC       shared-objects/db/c.o
  ...
  ln -fs librocksdb.so.5.8.0 librocksdb.so
  ln -fs librocksdb.so.5.8.0 librocksdb.so.5
  ln -fs librocksdb.so.5.8.0 librocksdb.so.5.8
~/rocksdb$ sudo make install-shared
  GEN      util/build_version.cc
install -d /usr/local/lib
for header_dir in `find "include/rocksdb" -type d`; do \
        install -d /usr/local/$header_dir; \
done
for header in `find "include/rocksdb" -type f -name *.h`; do \
        install -C -m 644 $header /usr/local/$header; \
done
install -C -m 755 librocksdb.so.5.8.0 /usr/local/lib && \
        ln -fs librocksdb.so.5.8.0 /usr/local/lib/librocksdb.so.5.8 && \
        ln -fs librocksdb.so.5.8.0 /usr/local/lib/librocksdb.so.5 && \
        ln -fs librocksdb.so.5.8.0 /usr/local/lib/librocksdb.so
```


# protobuf (≥ 2.6.1)
检查是否已经安装
```bash
$ protoc --version
libprotoc 3.0.0
```
安装1
1. 下载安装包
```bash
wget https://github.com/protocolbuffers/protobuf/releases/download/v2.6.1/protobuf-2.6.1.tar.gz
tar -xvzf protobuf-2.6.1.tar.gz
cd protobuf-2.6.1
```
2. 编译安装
```bash
./configure
make
make check
sudo make install
sudo ldconfig # refresh shared library cache.
```
3. 检查是否安装成功
```bash
protoc --version
libprotoc 2.6.1
```

安装2
1. 先安装bazel
https://github.com/bazelbuild/bazel/releases/tag/8.2.1
```bash
wget https://github.com/bazelbuild/bazel/releases/download/8.2.1/bazel-8.2.1-installer-linux-x86_64.sh
chmod +x bazel-8.2.1-installer-linux-x86_64.sh
./bazel-8.2.1-installer-linux-x86_64.sh
```

2. 下载protobuf
```bash
git clone https://github.com/protocolbuffers/protobuf.git
cd protobuf
git submodule update --init --recursive
git checkout v30.2
```
3. 编译protobuf
```bash
bazel build :protoc :protobuf
cp bazel-bin/protoc /usr/local/bin
```
4. 如果cmake仍然找出旧版本的protobuf
```bash
sudo ldconfig
# 检查 protobuf 库是否被正确识别
ldconfig -p | grep libprotobuf
# 确认库文件权限（应有执行权限）
ls -l /usr/local/lib/libprotobuf.so*
# 示例：删除Ubuntu自带版本
sudo apt remove libprotobuf-dev protobuf-compiler
```

# cryptopp (≥ 6.1.0)
检查是否已经安装
```bash
$ dpkg -l | grep cryptopp
$  find /usr/lib /usr/local/lib -name "libcryptopp*"
/usr/local/lib/libcryptopp.a
$ find /usr/include /usr/local/include -name "cryptopp"
/usr/local/include/cryptopp
```
安装
```bash
~$ git clone https://github.com/weidai11/cryptopp.git
Cloning into 'cryptopp'...
remote: Enumerating objects: 28997, done.
remote: Counting objects: 100% (185/185), done.
remote: Compressing objects: 100% (57/57), done.
remote: Total 28997 (delta 148), reused 128 (delta 128), pack-reused 28812 (from 2)
Receiving objects: 100% (28997/28997), 27.73 MiB | 1.64 MiB/s, done.
Resolving deltas: 100% (21117/21117), done.
~$ cd cryptopp/
~/cryptopp$ git checkout CRYPTOPP_6_1_0
Note: checking out 'CRYPTOPP_6_1_0'.

You are in 'detached HEAD' state. You can look around, make experimental
changes and commit them, and you can discard any commits you make in this
state without impacting any branches by performing another checkout.

If you want to create a new branch to retain commits you create, you may
do so (now or later) by using -b with the checkout command again. Example:

  git checkout -b <new-branch-name>

HEAD is now at 5be140bc Prepare for Crypto++ 6.1 release
~/cryptopp$ make
g++ -DNDEBUG -g2 -O3 -fPIC -pthread -pipe -c cryptlib.cpp
g++ -DNDEBUG -g2 -O3 -fPIC -pthread -pipe -c cpu.cpp
g++ -DNDEBUG -g2 -O3 -fPIC -pthread -pipe -c integer.cpp
...
g++ -DNDEBUG -g2 -O3 -fPIC -pthread -pipe -c dlltest.cpp
g++ -DNDEBUG -g2 -O3 -fPIC -pthread -pipe -c fipsalgt.cpp
g++ -o cryptest.exe -DNDEBUG -g2 -O3 -fPIC -pthread -pipe adhoc.o test.o bench1.o bench2.o validat0.o validat1.o validat2.o validat3.o validat4.o datatest.o regtest1.o regtest2.o regtest3.o dlltest.o fipsalgt.o ./libcryptopp.a  
~/cryptopp$ sudo make install
install -m 644 *.h /usr/local/include/cryptopp
install -m 644 libcryptopp.a /usr/local/lib
install cryptest.exe /usr/local/bin
install -m 644 TestData/*.dat /usr/local/share/cryptopp/TestData
install -m 644 TestVectors/*.txt /usr/local/share/cryptopp/TestVectors
```

# boost (≥ 1.67)
检查是否安装
```bash
$ dpkg -l | grep boost
ii  libboost-all-dev            1.65.1.0ubuntu1     amd64       Boost C++ Libraries development files (ALL) (default version)
ii  libboost-atomic-dev:amd64   1.65.1.0ubuntu1     amd64       atomic data types, operations, and memory ordering constraints (default version)
...

```
安装
```bash
$ wget https://archives.boost.io/release/1.67.0/source/boost_1_67_0.tar.gz
--2025-04-16 09:50:45--  https://archives.boost.io/release/1.67.0/source/boost_1_67_0.tar.gz
...
Saving to: ‘boost_1_67_0.tar.gz’
boost_1_67_0.tar.gz                100%[===============================================================>]  98.58M   926KB/s    in 2m 0s   
2025-04-16 09:52:45 (844 KB/s) - ‘boost_1_67_0.tar.gz’ saved [103363944/103363944]
$ tar -xzvf boost_1_67_0.tar.gz 
...
boost_1_67_0/tools/quickbook/test/xml_escape-1_5.quickbook
boost_1_67_0/tools/quickbook/Jamfile.v2
boost_1_67_0/tools/quickbook/_clang-format
boost_1_67_0/tools/quickbook/index.html
boost_1_67_0/tools/Jamfile.v2
boost_1_67_0/tools/index.html
boost_1_67_0/tools/make-cputime-page.pl
$ ./b2 install
(base) xinyu.chen@246:~/Ledgerdatabase_deps/boost_1_67_0$ ./b2 install
/home/xinyu.chen/Ledgerdatabase_deps/boost_1_67_0/libs/predef/check/../tools/check/predef.jam:46: Unescaped special character in argument $(language)::$(expression)
Performing configuration checks

    - default address-model    : 64-bit
    - default architecture     : x86
    - symlinks supported       : yes
    - C++11 mutex              : yes
    - lockfree boost::atomic_flag : yes
    - Boost.Config Feature Check: cxx11_auto_declarations : yes
    - Boost.Config Feature Check: cxx11_constexpr : yes
    - Boost.Config Feature Check: cxx11_defaulted_functions : yes
    - Boost.Config Feature Check: cxx11_final : yes
    - Boost.Config Feature Check: cxx11_hdr_mutex : yes
    - Boost.Config Feature Check: cxx11_hdr_regex : yes
    - Boost.Config Feature Check: cxx11_hdr_tuple : yes
    - Boost.Config Feature Check: cxx11_lambdas : yes
    - Boost.Config Feature Check: cxx11_noexcept : yes
    - Boost.Config Feature Check: cxx11_nullptr : yes
    - Boost.Config Feature Check: cxx11_rvalue_references : yes
    - Boost.Config Feature Check: cxx11_template_aliases : yes
    - Boost.Config Feature Check: cxx11_thread_local : yes
    - Boost.Config Feature Check: cxx11_variadic_templates : yes
    - has_icu builds           : yes
...

```

# Intel Threading Building Block (tbb_2020 version)
检查是否安装
```bash
$ dpkg -l | grep tbb
ii  libtbb2                   2020.3-0ubuntu1       amd64        Intel Threading Building Blocks
ii  libtbb-dev                2020.3-0ubuntu1       amd64        Intel Threading Building Blocks development files
```