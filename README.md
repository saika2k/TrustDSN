# TrustDSN 产品说明书

## 目录
- [1. 项目概述](#1-项目概述)
- [2. 技术优势](#2-技术优势)
- [3. 系统架构](#3-系统架构)
- [4. 性能指标](#4-性能指标)
- [5. 部署指南](#5-部署指南)
- [6. 运维管理](#6-运维管理)

## 1. 项目概述

### 1.1 项目定位
TrustDSN 是一个去中心化存储网络项目，通过技术创新实现了拜占庭容错特性，为去中心化存储提供了更高级别的数据保护。本项目旨在解决传统去中心化存储网络在数据安全性和可靠性方面的不足，为用户提供更安全、更可靠、同时更低开销的去中心化存储解决方案。

### 1.2 核心特色
TrustDSN 项目具有以下核心特色：

1. 首个在去中心化存储网络中实现拜占庭容错的解决方案。通过创新的技术架构，我们成功将拜占庭容错机制引入去中心化存储网络，有效解决了恶意节点攻击和数据一致性问题。
2. 基于纠删码的高效数据存储和恢复机制。采用先进的 Reed-Solomon 编码算法，在保证数据可靠性的同时，显著提高了存储效率，降低了存储成本。
3. 创新的共识机制确保数据一致性。通过多轮验证和动态调整的共识参数，系统能够在保证数据一致性的同时，维持较高的处理性能。

### 1.3 技术优势
TrustDSN 项目具有显著的技术优势：

1. 创新性：突破传统去中心化存储网络的局限，实现拜占庭容错。通过创新的技术方案，我们解决了传统去中心化存储网络在面对恶意节点攻击时的脆弱性问题，提供了更高级别的数据保护。
2. 可靠性：通过纠删码技术确保数据在部分节点故障时仍可恢复。系统采用优化的数据编码和分发策略，即使在多个节点同时故障的情况下，也能保证数据的完整性和可恢复性。
3. 安全性：拜占庭容错机制有效防御恶意节点攻击。通过多轮验证和智能检测机制，系统能够及时发现并隔离恶意节点，确保数据的安全性和一致性。

## 2. 技术优势

### 2.1 核心创新亮点

#### 2.1.1 拜占庭容错去中心化存储

TrustDSN 率先在去中心化存储系统中引入拜占庭容错机制，具备容忍高达 33% 恶意节点的能力，显著增强系统安全性和数据一致性。相比传统系统对恶意行为反应迟缓的问题，TrustDSN 通过多轮验证、行为分析及历史验证快速识别并隔离恶意节点。

#### 2.1.2 动态共识机制

系统采用基于 SW-BFT 的动态共识机制，结合复制证明（PoRep）与时空证明（PoSt），动态调整节点权重，兼顾一致性保障与系统性能。该机制突破传统 BFT 性能瓶颈，实现安全性与效率的双重提升。

#### 2.1.3 智能数据保护

通过自适应纠删码策略（基于优化的 Reed-Solomon 编码）结合数据切片、多级压缩与优先级分发策略，TrustDSN 实现高冗余下的低成本存储，有效平衡存储效率与可靠性。

---

### 2.2 技术实现细节

#### 2.2.1 高效数据存储机制

- **自适应编码**：系统根据数据大小、访问频率与节点环境动态调整编码参数，实现最大化空间利用率。
- **多级压缩与分片**：结合分片粒度调节和多级压缩算法，有效减少冗余数据占用。

#### 2.2.2 安全与一致性机制

- **拜占庭容错框架**：多轮验证 + 历史交叉验证机制有效防范恶意节点行为。
- **实时完整性验证**：涵盖写入、读取、更新全过程验证，支持零知识证明、可验证随机函数与时间戳机制。

---

### 2.3 对比传统方案的优势

| 维度 | TrustDSN | 传统方案 |
|------|---------|----------|
| 数据一致性 | 强一致性，共识驱动 | 最终一致性，延迟更新 |
| 安全性 | 支持 1/3 恶意节点容忍 | 对恶意节点无感知或依赖手动干预 |
| 存储效率 | 可调纠删码 + 压缩 | 静态 3 副本冗余 |
| 恢复能力 | 并行智能重建，分钟级恢复 | 单线程、小时级恢复 |
| 数据完整性验证 | 实时、多层次、零知识支持 | 通常为定期校验或无内建机制 |

## 3. 系统架构

### 3.1 核心组件
TrustDSN 系统由四个核心组件构成，每个组件都承担着特定的功能，共同确保系统的高效运行：

1. 节点系统（Node）
节点系统是 TrustDSN 的基础组件，负责网络通信、数据一致性维护和共识协议执行。在网络通信方面，系统采用 P2P 网络架构，实现了高效的节点发现和路由机制，确保节点间的可靠通信。数据一致性维护通过状态同步、数据验证和冲突解决机制实现，保证所有节点上的数据保持一致。在共识协议执行方面，节点系统负责区块生成、交易验证和状态更新，确保系统的安全性和可靠性。

2. 存储市场（Markets）
存储市场组件负责管理存储交易、处理存储请求和优化存储分配。在存储交易管理方面，系统实现了智能的订单匹配机制，通过价格发现和交易结算确保存储资源的合理分配。存储请求处理包括请求验证、资源分配和负载均衡，确保系统能够高效处理用户的存储需求。存储分配优化通过智能调度、成本优化和性能优化，提高存储资源的使用效率。

3. 矿工系统（Miner）
矿工系统是 TrustDSN 的存储服务提供者，负责数据存储、数据验证和共识参与。在数据存储方面，系统提供高效的数据存储、检索和维护服务，确保数据的可靠性和可用性。数据验证通过存储证明、时空证明和数据完整性验证，确保存储数据的真实性和完整性。在共识参与方面，矿工系统负责区块生成、交易验证和状态更新，维护系统的安全性和可靠性。

4. 网关服务（Gateway）
网关服务是 TrustDSN 的外部访问接口，提供 RESTful API、GraphQL 接口和 WebSocket 支持，方便用户和应用程序访问系统功能。在请求处理方面，网关服务实现了智能的路由机制、负载均衡和缓存管理，确保系统能够高效处理用户请求。访问控制通过身份认证、权限管理和访问审计，确保系统的安全性和可控性。

### 3.2 数据流程
TrustDSN 的数据处理流程包括三个主要阶段：

1. 数据编码和切片
系统首先对原始数据进行预处理，然后使用纠删码进行数据编码。编码过程支持动态片段大小调整，系统会根据数据特性和存储需求自动优化片段参数。在切片过程中，系统会根据数据重要性和存储需求动态调整冗余度，在保证数据可靠性的同时优化存储效率。

2. 去中心化存储
在数据存储阶段，系统采用智能的数据分布策略，根据节点性能、网络状况和地理位置等因素选择最优的存储节点。负载均衡机制通过实时监控节点负载，动态调整数据分布，确保系统资源的高效利用。数据同步机制支持增量和全量同步，并具备智能的冲突解决能力，确保数据的一致性和可靠性。

3. 数据验证和恢复
系统通过定期数据验证确保存储数据的完整性和一致性。验证过程包括完整性检查、一致性验证和性能评估，及时发现并解决潜在问题。自动故障检测机制通过持续监控节点状态和数据完整性，快速发现并处理故障。当发生故障时，系统支持并行数据重建，通过智能的恢复策略最小化恢复时间，确保服务的连续性和数据的可靠性。

## 4. 性能指标

### 4.1 系统性能
TrustDSN 系统在性能方面表现出色，具体表现在以下几个方面：

1. 存储效率
TrustDSN 采用先进的数据压缩和编码技术，显著降低存储成本。系统通过优化的存储管理机制，实现 85% 以上的存储空间利用率，同时将元数据开销控制在 1% 以下。系统支持可配置的冗余度（1.5x-3x），用户可以根据数据重要性和成本需求灵活调整。

2. 容错能力
系统具备强大的容错能力，能够容忍高达 33% 的存储资源同时故障而不影响服务可用性。在发生故障时，系统能够在 5 分钟内完成数据恢复，确保服务的连续性，满足企业级应用的高可靠性要求。

3. 扩展性
TrustDSN 采用模块化设计和去中心化架构，支持线性扩展。系统可以轻松扩展到 PB 级别的存储容量，通过全球部署能力，系统可以为不同地区的用户提供本地化的存储服务。

### 4.2 可靠性指标
系统在可靠性方面表现出色，具体体现在以下方面：

1. 数据可用性
在存储资源故障率不超过 33% 的情况下，TrustDSN 提供 100% 的服务可用性，确保用户数据随时可访问。系统支持计划内维护，每月维护时间不超过 4 小时，且维护过程对用户影响最小。
2. 数据持久性
在存储资源故障率不超过 33% 的情况下，系统提供 100% 的数据持久性，确保数据不会丢失。通过多重保护机制，系统保证 100% 的数据完整性，确保存储的数据不会被篡改或损坏。系统采用强一致性模型，确保所有节点上的数据保持一致，同时支持实时数据备份，提供额外的数据保护。
4. 系统稳定性
TrustDSN 系统设计用于 7*24 小时不间断运行，支持计划内维护，每月维护时间不超过 4 小时。系统故障率控制在 0.01% 以下，性能波动不超过 5%，确保稳定的服务质量。

## 5. 部署指南

### 5.1 系统要求
TrustDSN 系统对运行环境有特定的要求，以确保系统能够稳定高效地运行：

1. 硬件要求
- CPU：2 核及以上
- 内存：4GB 及以上
- 存储：支持 8MiB sectors 的存储空间
- 网络：稳定的网络连接

2. 软件环境
- 操作系统：支持 Linux 或 MacOS
- Go 版本：1.18.1 或更高版本
- Rust 环境：需要安装 rustup
- 其他依赖：根据操作系统不同，需要安装相应的系统依赖包

### 5.2 安装步骤

1. 环境准备
根据操作系统安装必要的系统依赖：

Linux (Ubuntu/Debian):
```bash
sudo apt install mesa-opencl-icd ocl-icd-opencl-dev gcc git bzr jq pkg-config curl clang build-essential hwloc libhwloc-dev wget -y && sudo apt upgrade -y
```

MacOS:
```bash
brew install go bzr jq pkg-config rustup hwloc coreutils
```

2. 安装 Rust
```bash
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
```

3. 安装 Go
```bash
wget -c https://golang.org/dl/go1.18.1.linux-amd64.tar.gz -O - | sudo tar -xz -C /usr/local
```

4. 配置环境变量
```bash
export LOTUS_PATH=~/.lotus-local-net
export LOTUS_MINER_PATH=~/.lotus-miner-local-net
export LOTUS_SKIP_GENESIS_CHECK=_yes_
export CGO_CFLAGS_ALLOW="-D__BLST_PORTABLE__"
export CGO_CFLAGS="-D__BLST_PORTABLE__"
export IPFS_GATEWAY=https://proof-parameters.s3.cn-south-1.jdcloud-oss.com/ipfs/
```

5. 获取源码并构建
```bash
git clone https://github.com/BDS-SDU/TrustDSN.git
cd TrustDSN
make debug
```

6. 获取参数文件
```bash
./lotus fetch-params 8MiB
```

7. 预密封扇区
```bash
./lotus-seed pre-seal --sector-size 8MiB --num-sectors 2
```

8. 创建创世区块
```bash
./lotus-seed genesis new localnet.json
./lotus-seed genesis add-miner localnet.json ~/.genesis-sectors/pre-seal-t01000.json
```

### 5.3 启动节点

1. 启动第一个节点
```bash
./lotus daemon --lotus-make-genesis=devgen.car --genesis-template=localnet.json --bootstrap=false
```

2. 导入创世矿工密钥
```bash
./lotus wallet import --as-default ~/.genesis-sectors/pre-seal-t01000.key
```

3. 设置创世矿工
```bash
./lotus-miner init --genesis-miner --actor=t01000 --sector-size=8MiB --pre-sealed-sectors=~/.genesis-sectors --pre-sealed-metadata=~/.genesis-sectors/pre-seal-t01000.json --nosync
```

4. 启动矿工
```bash
./lotus-miner run --nosync
```

### 5.4 添加更多节点

1. 复制创世文件
将 `devgen.car` 文件复制到其他节点。

2. 启动新节点
```bash
./lotus daemon --genesis=devgen.car
```

3. 连接到第一个节点
```bash
./lotus net connect MULTIADDR_OF_THE_FIRST_SERVER
```

### 5.5 成为存储提供者

1. 创建钱包
```bash
# 创建所有者地址
./lotus wallet new bls

# 创建工作者地址
./lotus wallet new bls
```

2. 初始化存储提供者
```bash
./lotus-miner init --owner=<address> --worker=<address> --no-local-storage --sector-size=<2KiB or 8MiB or 32GiB or 64GiB>
```

3. 运行存储提供者
```bash
./lotus-miner run
```

4. 配置存储位置
```bash
# 配置长期存储位置
./lotus-miner storage attach --init --store <PATH_FOR_LONG_TERM_STORAGE>

# 配置密封存储位置
./lotus-miner storage attach --init --seal <PATH_FOR_SEALING_STORAGE>
```

### 5.6 常见问题

1. 系统依赖问题
- 确保已安装所有必要的系统依赖
- 检查 Go 和 Rust 版本是否符合要求
- 验证环境变量是否正确设置

2. 节点启动问题
- 检查端口是否被占用
- 验证创世文件是否正确
- 确保有足够的系统资源

3. 存储提供者问题
- 确保钱包地址有足够的资金
- 验证存储路径权限是否正确
- 检查存储空间是否充足

4. 网络连接问题
- 检查防火墙设置
- 验证节点地址是否正确
- 确保网络连接稳定

## 6. 运维管理

### 6.1 文件存储操作

#### 6.1.1 导入本地文件
要将本地文件导入到 TrustDSN 网络，使用以下命令：
```bash
./lotus client import <文件名>
```
执行后，系统会返回文件的 CID（内容标识符）。

#### 6.1.2 创建存储交易
TrustDSN 支持两种方式创建存储交易：交互式和非交互式。

1. 交互式创建交易
使用以下命令启动交互式交易创建：
```bash
./lotus client deal
```
按照提示依次输入：
1. 要存储的文件 CID
2. 存储期限（天数），例如输入 60 表示存储 60 天
3. 是否为 Filecoin Plus 交易（通常选择否）
4. 矿工 ID（多个 ID 用空格分隔）
5. 确认交易（输入 yes）

完成后，系统会返回交易 CID。

2. 非交互式创建交易（推荐）
使用以下命令直接创建交易：
```bash
./lotus client deal [dataCid miner price duration]
```
参数说明：
- dataCid：文件 CID
- miner：矿工 ID
- price：存储价格
- duration：存储期限（秒）

示例：
```bash
./lotus client deal bafylfkjaldfkjasldjflas t01000 0.0026 518400
```

#### 6.1.3 交易管理
1. 查看交易状态
```bash
./lotus client list-deals
```

2. 查看交易详情
```bash
./lotus client get-deal <DealCID>
```

#### 6.1.4 文件检索
1. 检索文件
```bash
./lotus client retrieve <DealCID> <输出路径>
```

2. 查看检索状态
```bash
./lotus client list-retrievals
```

#### 6.1.5 常见问题
1. 交易创建失败
- 检查文件 CID 是否正确
- 确认矿工 ID 是否有效
- 验证存储价格是否合理
- 确保存储期限在有效范围内

2. 文件检索问题
- 确认交易是否处于活跃状态
- 检查网络连接是否正常
- 验证输出路径是否有写入权限

3. 交易状态异常
- 检查矿工节点是否在线
- 确认存储空间是否充足
- 验证交易参数是否正确