#!/bin/bash
source env.sh

# 颜色定义
BLUE='\033[1;34m'
RED='\033[1;31m'
NC='\033[0m' # No Color

# 输出函数
print_status() {
    echo -e "${BLUE}[BUILD]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

BUILD_TYPE=Release
CC=gcc
CXX=g++
ADDITIONAL_FLAGS=""

# 解析命令行参数
while [[ $# -gt 0 ]]; do
    case $1 in
        --build-type)
        BUILD_TYPE="$2"
        shift # 移动到下一个参数
        shift # 移动到下一个参数
        ;;
        --cc)
        CC="$2"
        shift # 移动到下一个参数
        shift # 移动到下一个参数
        ;;
        --cxx)
        CXX="$2"
        shift # 移动到下一个参数
        shift # 移动到下一个参数
        ;;
        --add-flags)
        ADDITIONAL_FLAGS="$2"
        shift # 移动到下一个参数
        shift # 移动到下一个参数
        ;;
        *)
        print_error "Unknown option $1"
        exit 1
        ;;
    esac
done

# 设置默认构建类型并转换为小写
BUILD_TYPE=$(echo "$BUILD_TYPE" | tr '[:upper:]' '[:lower:]')
# 验证构建类型
case "$BUILD_TYPE" in
    debug|release)
        print_status "Build type: ${BUILD_TYPE}"
        ;;
    *)
        print_error "Invalid build type. Use 'debug' or 'release'"
        exit 1
        ;;
esac

# 根据构建类型设置目录
BUILD_DIR="${builddir}/build_${BUILD_TYPE}_letus"

# 清理旧的构建目录
print_status "Cleaning old build directories..."
rm -rf ${BUILD_DIR}/

# 清理旧的数据
print_status "Cleaning old data directories..."
rm -rf ${datadir}/data
rm -rf ${datadir}/index

# 创建必要的目录
print_status "Creating directories..."
mkdir -p ${BUILD_DIR}/
mkdir -p ${datadir}/data/
mkdir -p ${datadir}/index/

# 进入构建目录并保存项目根目录
PROJECT_ROOT="${rootdir}/src/LETUS"
cd "${BUILD_DIR}" || {
    print_error "Failed to enter build directory"
    exit 1
}

# 运行 CMAKE
print_status "Running CMAKE..."
BUILD_TYPE_UPPER=$(echo "$BUILD_TYPE" | tr '[:lower:]' '[:upper:]')
if ! cmake -DCMAKE_BUILD_TYPE=${BUILD_TYPE_UPPER} \
          -DCMAKE_C_COMPILER=${CC} \
          -DCMAKE_CXX_COMPILER=${CXX} \
          -DCMAKE_CXX_FLAGS="${ADDITIONAL_FLAGS}" \
          -DCMAKE_VERBOSE_MAKEFILE=ON \
          "${PROJECT_ROOT}"; then
    print_error "CMAKE configuration failed"
    exit 1
fi

# 编译项目
print_status "Building project..."
if ! make -j$(nproc); then
    print_error "Build failed"
    exit 1
fi

print_status "Build completed successfully!"
