---
# 🐙 OctoStudio

<div align="center">

![License](https://img.shields.io/badge/License-MIT-green.svg)
![Python](https://img.shields.io/badge/Python-3.10+-3776AB?style=flat-square&logo=python&logoColor=white)
![Go](https://img.shields.io/badge/Go-00ADD8?style=flat-square&logo=go&logoColor=white)
![TypeScript](https://img.shields.io/badge/TypeScript-3178C6?style=flat-square&logo=typescript&logoColor=white)
![NVIDIA](https://img.shields.io/badge/NVIDIA-76B900?style=flat-square&logo=nvidia&logoColor=white)
![Hugging Face](https://img.shields.io/badge/%F0%9F%A4%97-Hugging%20Face-FFD21E?style=flat-square)

</div>
<div align="center">
**一个全能型 LLM 工作站：集模型下载、多引擎部署、硬件监控于一体。**

[🏠 官方文档](https://www.google.com/search?q=%23) | [🚀 快速开始](https://www.google.com/search?q=%23) | [🛠️ 贡献指南](https://www.google.com/search?q=%23) | [💬 反馈建议](https://www.google.com/search?q=%23)

</div>
---

## 🌟 核心理念

**OctoStudio** (八爪鱼工作室) 旨在打破大模型使用的门槛。它不仅是一个聊天 UI，更是一个完整的**本地 AI 运营中心**。通过集成的多语言优势，实现从网络请求到深度推理的全栈优化。

### 核心能力一览：

- **多语言协同**:
- **Go**: 驱动高性能、高并发的模型下载引擎。
- **Python**: 深度适配 `llama-cpp-python` 与 `vLLM`。
- **TypeScript**: 构建基于 Next.js 的响应式、可视化监控面板。

- **一站式工作流**: 搜索 (HuggingFace) -> 下载 -> 部署 (Llama_cpp/vLLM) -> 聊天 -> 监控。
- **混合动力**: 灵活切换本地权重与 API（OpenAI, Claude, DeepSeek）。

---

## 📊 硬件资源实时监控

OctoStudio 内置了毫秒级响应的监控模块，你可以直接在聊天界面边框查看：

| 监控指标        | 技术实现         | 说明                                  |
| --------------- | ---------------- | ------------------------------------- |
| **GPU VRAM**    | `nvidia-smi` API | 实时显示显存占用，防止 OOM (显存溢出) |
| **Token Speed** | 推理引擎反馈     | 实时计算每秒生成 Token 数 (T/s)       |
| **CPU/RAM**     | `gopsutil` (Go)  | 监控后台进程对系统资源的整体消耗      |

---

## 🛠️ 技术栈

- **前端**: React 18, TypeScript, Tailwind CSS, Lucide Icons.
- **后端服务**:
- **Main API**: Go (Gin/Echo) —— 负责文件系统、系统监控。
- **Inference Server**: Python (FastAPI) —— 负责模型热加载与推理转换。

- **推理后端**:
- `llama.cpp` (GGUF 格式支持)
- `vLLM` (高性能并行推理)

---

## 🚀 快速安装

```bash
# 1. 克隆并进入目录
git clone https://github.com/YourName/OctoStudio.git
cd OctoStudio

# 2. 启动 Go 后端 (负责监控与下载)
cd backend-core && go run main.go

# 3. 启动 Python 推理服务
cd inference-engine && pip install -r requirements.txt && python server.py

# 4. 启动前端 UI
cd frontend && npm install && npm run dev

```

---

## 🗺️ 路线图 (Roadmap)

- [x] 基于 Go 的高性能下载器 (Multi-threaded)
- [x] Llama.cpp 后端集成
- [ ] vLLM 分布式部署支持
- [ ] 自定义硬件阈值报警
- [ ] 插件系统：支持 Web 搜索

---

## 📄 开源协议

本项目基于 **MIT License**。

---
