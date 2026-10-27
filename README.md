<p align="center">
  <img src="docs/assets/ouroboros-logo.png" alt="Ouroboros logo" width="360">
</p>

# Ouroboros · 授权安全测试 IDE 与编译器

Ouroboros 由宁夏大学鑫启网安俱乐部（[NingXia-XinQiCyberClub](https://github.com/NingXia-XinQiCyberClub)）发起，是一个面向**授权渗透测试、攻防演练与 CTF 竞赛**的 TUI/GUI 双形态安全工作台。它把"明确范围的目标"编译成可审计的计划、工具链动作、运行时工作节点状态、可回放证据与操作员控制的执行门。

> **合规声明**：本项目仅用于已获书面授权的测试环境、靶场与教学研究。使用前请确认你已获得目标所有者的明确授权，并遵守所在地法律法规。严禁对任何未授权系统实施测试。

---

## 一、项目定位

传统做法里，扫描脚本、利用工具与结果整理散落各处，过程无法复现。Ouroboros 的思路是把整个测试过程做成一个**可编译、可回放、可审计**的工作区：

| 能力 | 说明 |
| :-- | :-- |
| 目标驱动 | 以"目标（Goal）"为输入，域包（Web/API、内网、CTF/Lab）产出假设与语义动作 |
| 计划编译 | 编译器把语义动作编译为具体工具计划，LLM 不直接选择原始工具 |
| 证据链 | 每一步的输入输出落盘为事件（replay），证据图（Evidence Graph）只认证据事件 |
| 执行门 | 工具执行、越权旁路、报告生成均需操作员显式批准，旁路可审计、可回收 |
| 回放与报告 | `ouro replay` 完整回放时间线，一键生成 Markdown 报告 |
| AI 辅助 | 可接入 OpenAI 兼容模型做规划预览与提示，默认离线，密钥只经环境变量注入 |

与"AI 工具路由器"的区别：模型永远不直接跑工具——域包出假设，编译器出计划，状态只在证据事件上推进。

## 二、快速开始

要求：Go 1.26+（`go.mod` 为准）。国内网络环境 Makefile 已默认接入 `GOPROXY=https://goproxy.cn,direct`。

```bash
# 克隆
git clone https://github.com/NingXia-XinQiCyberClub/ouroboros.git
cd ouroboros

# 构建并跑全量测试（CI 同款）
make build
make test

# 初始化一个工作区并体验 TUI
go run ./cmd/ouro init ./workspace
go run ./cmd/ouro run ./workspace --tui
go run ./cmd/ouro replay ./workspace --tui
```

桌面 IDE（Electron + Go 后端）：

```bash
cd gui
npm install
npm run dev        # 开发模式
npm run build      # 构建
```

常用命令速查（详见 [英文版 README](docs/README.en.md#common-commands)）：

```bash
make doctor              # 工作区体检
make doctor-strict       # 严格模式：未授权的安全旁路直接阻断
make smoke-ide           # "现在能用吗"冒烟
make smoke-workflow      # init→run→replay→workflow-status 全链路冒烟
```

## 三、目录结构

```text
ouroboros/
├── cmd/
│   ├── ouro/          # CLI 入口（init/run/replay/provider/toolchain…）
│   └── desktop/       # 桌面 IDE 构建与契约检查入口
├── internal/
│   ├── app/           # 应用层：运行编排、安全门、工作节点租约
│   ├── compiler/      # 语义动作 → 工具计划编译器
│   ├── domain/        # 域包：Web/API、内网、CTF/Lab
│   ├── evidence/      # 证据图与证据门
│   ├── eventlog/      # 回放事件存储
│   ├── sisyphus/      # 任务规划引擎（西西弗斯循环）
│   ├── swarm/         # 多代理工作节点编排
│   ├── toolchain/     # 工具链清单、托管负载、执行隔离
│   ├── tui/           # Bubble Tea TUI
│   ├── scopegate/     # 目标范围与授权边界
│   └── …              # sessionguard / modelruntime / nativetools / patterns 等
├── gui/               # Electron 桌面 IDE（React + TypeScript）
├── toolchains/        # 域包工具清单与托管 bundle
├── configs/           # 策略模式配置、provider 示例
├── scripts/           # build / release / install 脚本
├── test/              # 测试夹具
└── docs/              # 架构设计书、里程碑记录、开发文档
```

## 四、安全边界

工作区把"范围"和"授权"作为一等公民：

- **目标范围**：`scope.target_policy: enforce` 时，命令行参数中的目标字面量（含 `--url=`、`targets=a,b` 等嵌入形式）都会做边界校验；越界直接阻断。
- **可审计的旁路**：`scope.target_policy: disabled` 与 `execution.isolation: disabled` 是为误报/遗留工具场景保留的显式关闭路径，旁路状态始终记录在案；`authorization.yaml` 支持过期时间，过期后严格模式立刻阻断。
- **十二个副作用门**：运行时工作节点的队列认领、派发、回传、租约心跳等十二个直接副作用入口全部要求当前有效的安全旁路授权。
- **证据最小化**：各解析器（nuclei/httpx/nmap/subfinder/naabu/dnsx/katana/gau/ffuf/dalfox 等）只持久化范围内目标的最小事实字段，原始流量、响应体、查询串、越界主机一律不落盘。
- **产物边界**：本地产物扫描限定 `workspace/artifacts/` 之下，符号链接逃逸与路径穿越在执行前失败，且不受 `target_policy: disabled` 豁免。

完整边界语义见 [英文版 README · Safety Boundary](docs/README.en.md#safety-boundary)。

## 五、GateMesh 竞赛代理

内置 [GateMesh](Dockerfile.gatemesh) 模式：把 Ouroboros 打包为竞赛环境镜像，包含 canary 检测、凭据登录、受限收获（gated harvest）、范围门、认知门（进度/预算）、证据门与传送带等完整合规链路，用于 CTF 平台上的自动化参赛代理。构建方式见 `Dockerfile.gatemesh`。

## 六、文档

| 文档 | 内容 |
| :-- | :-- |
| [docs/README.en.md](docs/README.en.md) | 英文完整版 README（能力清单、安全边界全量细节） |
| [docs/development.md](docs/development.md) | 开发指南 |
| [docs/Ouroboros_平台架构设计书.md](docs/Ouroboros_平台架构设计书.md) | 平台架构设计书（中文） |
| [docs/progress/](docs/progress/) | 里程碑记录 |
| [DESIGN.md](DESIGN.md) · [PRODUCT.md](PRODUCT.md) | 设计与产品定义 |
| [AGENTS.md](AGENTS.md) | 代码库导航（面向 AI 代理与贡献者） |

## 七、参与贡献

按组织[贡献指南](https://github.com/NingXia-XinQiCyberClub/.github/blob/main/CONTRIBUTING.md)执行：分支命名 `<类型>/<描述>`、Conventional Commits、变更请求合入。入门任务见各仓库 `good first issue` 标签。

讨论与问答在组织 [Discussions](https://github.com/NingXia-XinQiCyberClub/.github/discussions) 进行。

## 许可证

[MIT](LICENSE) · © 2026 Zhang_Haiyang（SeaC-25）与宁夏大学鑫启网安俱乐部贡献者
