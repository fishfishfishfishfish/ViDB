> 本文由 [简悦 SimpRead](http://ksria.com/simpread/) 转码， 原文地址 [blog.csdn.net](https://blog.csdn.net/lian740930980/article/details/126659488)

1. 安装 gcc9

```
# Install GCC 9
sudo add-apt-repository ppa:ubuntu-toolchain-r/test
sudo apt-get update
sudo apt-get install gcc-9 g++-9
```

2. 安装 tbb

可以用代码下载，如果网络不好也可以网盘自取，下载后记得将文件名字改为 tbb-2019_U8，这样粘后面部分命令会方便很多。

链接：https://pan.baidu.com/s/1gPVgz5viGINtjgSzLj2BoA?pwd=yari   
提取码：yari 

```
wget https://github.com/intel/tbb/archive/2019_U8.tar.gz
tar zxvf 2019_U8.tar.gz
rm 2019_U8.tar.gz
 
cd tbb-2019_U8
```

*   编辑 linux.gcc-9.inc

```
cp build/linux.gcc.inc build/linux.gcc-9.inc 
 
vi build/linux.gcc-9.inc
 
# 把大约是15,16　行的CPLUS ?= g++　CONLY ?= gcc修改为下面语句，然后保存
CPLUS ?= g++-9
CONLY ?= gcc-9
```

*   编译 tbb

```
sudo mkdir /usr/local/tbb-2019_U8
sudo cp -r include /usr/local/tbb-2019_U8/include
sudo ln -s /usr/local/tbb-2019_U8/include/tbb /usr/local/include/tbb
sudo cp -r build/my_tbb_build_release /usr/local/tbb-2019_U8/lib
 
sudo ln -s /usr/local/tbb-2019_U8/lib/libtbb.so.2 /usr/local/lib/libtbb.so
sudo ln -s /usr/local/tbb-2019_U8/lib/libtbbmalloc.so.2 /usr/local/lib/libtbbmalloc.so
sudo ln -s /usr/local/tbb-2019_U8/lib/libtbbmalloc_proxy.so.2 /usr/local/lib/libtbbmalloc_proxy.so
echo 'export LD_LIBRARY_PATH=/usr/local/tbb-2019_U8/lib:$LD_LIBRARY_PATH' >> ~/.bashrc
source ~/.bashrc
```

*   安装 tbb

```
#include <algorithm>
#include <chrono>
#include <execution>
#include <iostream>
#include <random>
#include <vector>
 
void printDuration(std::chrono::steady_clock::time_point start, std::chrono::steady_clock::time_point end, const char *message)
{
    auto diff = end - start;
    std::cout << message << ' ' << std::chrono::duration<double, std::milli>(diff).count() << " msn";
}
template <typename T>
void test(const T &policy, const std::vector<double> &data, const int repeat, const char *message)
{
    for (int i = 0; i < repeat; ++i)
    {
        std::vector<double> curr_data(data);
 
        const auto start = std::chrono::steady_clock::now();
        std::sort(policy, curr_data.begin(), curr_data.end());
        const auto end = std::chrono::steady_clock::now();
        printDuration(start, end, message);
    }
    std::cout << 'n';
}
 
int main()
{
    // Test samples and repeat factor
    constexpr size_t samples{5'000'000};
    constexpr int repeat{10};
 
    // Fill a vector with samples numbers
    std::random_device rd;
    std::mt19937_64 mre(rd());
    std::uniform_real_distribution<double> urd(0.0, 1.0);
 
    std::vector<double> data(samples);
    for (auto &e : data)
    {
        e = urd(mre);
    }
 
    // Sort data using different execution policies
    std::cout << "std::execution::seqn";
    test(std::execution::seq, data, repeat, "Elapsed time");
 
    std::cout << "std::execution::parn";
    test(std::execution::par, data, repeat, "Elapsed time");
}
```

**３．测试**

*   **测试代码**

**创建 t0.cpp**

```
#include <algorithm>
#include <chrono>
#include <execution>
#include <iostream>
#include <random>
#include <vector>
 
void printDuration(std::chrono::steady_clock::time_point start, std::chrono::steady_clock::time_point end, const char *message)
{
    auto diff = end - start;
    std::cout << message << ' ' << std::chrono::duration<double, std::milli>(diff).count() << " msn";
}
template <typename T>
void test(const T &policy, const std::vector<double> &data, const int repeat, const char *message)
{
    for (int i = 0; i < repeat; ++i)
    {
        std::vector<double> curr_data(data);
 
        const auto start = std::chrono::steady_clock::now();
        std::sort(policy, curr_data.begin(), curr_data.end());
        const auto end = std::chrono::steady_clock::now();
        printDuration(start, end, message);
    }
    std::cout << 'n';
}
 
int main()
{
    // Test samples and repeat factor
    constexpr size_t samples{5'000'000};
    constexpr int repeat{10};
 
    // Fill a vector with samples numbers
    std::random_device rd;
    std::mt19937_64 mre(rd());
    std::uniform_real_distribution<double> urd(0.0, 1.0);
 
    std::vector<double> data(samples);
    for (auto &e : data)
    {
        e = urd(mre);
    }
 
    // Sort data using different execution policies
    std::cout << "std::execution::seqn";
    test(std::execution::seq, data, repeat, "Elapsed time");
 
    std::cout << "std::execution::parn";
    test(std::execution::par, data, repeat, "Elapsed time");
}
```

*   编译程序

```
g++-9 -std=c++17 -Wall -Wextra -pedantic -O2 t0.cpp -o t0_opt -L /usr/local/lib/ -ltbb
```

*   运行代码

```
./t0_opt
```

参考：[(1 条消息) linux tbb 安装_Ubuntu18.04 GCC9 安装_天街踏尽公卿骨的博客 - CSDN 博客](https://blog.csdn.net/weixin_32207065/article/details/112270765 "(1条消息) linux tbb 安装_Ubuntu18.04 GCC9 安装_天街踏尽公卿骨的博客-CSDN博客")