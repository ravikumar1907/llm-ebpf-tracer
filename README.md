# llm-ebpf-tracer
# 🔍 LLM eBPF Tracer

**Trace and analyze Large Language Model (LLM) inference workloads (e.g., PyTorch + CUDA) at the Linux kernel level using eBPF.**

This project enables system-level observability for modern AI inference pipelines — revealing what happens when LLMs interact with memory, CPUs, GPUs, and the filesystem under the hood.

---

## 🧠 Why This Project?

LLMs like LLaMA, Mistral, and GPT-variants require optimized infrastructure. While most focus on model design or training, inference performance depends heavily on OS-level behaviors like:

- File I/O for model weights (multi-GB mmap)
- CPU-GPU data paths and device access
- Thread pinning, scheduling, and NUMA locality
- Page locking (mlock) and memory pressure

Traditional profilers miss this. **eBPF (extended Berkeley Packet Filter)** allows us to trace these interactions with zero modification to the model code or OS kernel.

---

## 📦 Features

- ✅ Trace `mmap`, `mlock`, `openat` during model loading
- ✅ Visualize `sched_switch` to track inference thread behavior
- ✅ Monitor access to `/dev/nvidia*` (GPU devices)
- ✅ Support for PyTorch, TensorRT, and Hugging Face model servers
- ✅ Plug-and-play tracing using `bpftrace` and `bcc`
- 🧪 Optional integration with Prometheus and Grafana dashboards

---

## 🖼️ Architecture Overview

![LLM Tracing Architecture](./A_digital_illustration_presents_a_technological_ec.png)

---

## ⚙️ Example Tracing Scripts

### 🧩 Track `mmap` Usage (Model Weights Loading)

```bash
bpftrace -e 'tracepoint:syscalls:sys_enter_mmap { @[comm] = count(); }'
```

### 🔒 Monitor Memory Locking (PyTorch’s `mlock` Calls)

```bash
bpftrace -e 'tracepoint:syscalls:sys_enter_mlock { @[comm] = count(); }'
```

### 🔁 Visualize Scheduler Switches (CPU Thread Behavior)

```bash
bpftrace -e 'tracepoint:sched:sched_switch { @[prev_comm, next_comm] = count(); }'
```

### 🎮 Watch GPU Access (Device File Usage)

```bash
bpftrace -e 'tracepoint:syscalls:sys_enter_openat /str(args->filename) =~ "/dev/nvidia.*/" / { @[comm] = count(); }'
```

---

## 📊 Planned Dashboard (Work in Progress)

- 🔹 Model load phase duration
- 🔹 `mmap` vs `mlock` event frequency
- 🔹 CPU usage per inference worker
- 🔹 NUMA locality and scheduling efficiency
- 🔹 Live latency trends per token

---

## 📁 Repo Layout

```bash
.
├── scripts/               # bpftrace tracing scripts
│   ├── mmap.bt
│   ├── mlock.bt
│   ├── sched.bt
│   └── gpu_access.bt
├── dashboards/            # (WIP) Grafana dashboards for LLM workload visibility
├── benchmark/             # Scripts to benchmark LLM inference performance
├── README.md              # This file
└── LICENSE
```

---

## 💡 Use Cases

- Debug slow inference cold starts due to poor I/O or mmap
- Optimize thread pinning and scheduling latency
- Monitor GPU usage contention in shared environments
- Profile PyTorch and TensorRT workloads at system level
- Tune NUMA policies for large model inference

---
---

## 🧪 Example Outputs

### 📄 mmap.bt

#### 🔹 Output:
```
@[python3, 12345]: 5
@[llm_server, 9981]: 14
```

#### 💡 What it means:
- `python3` or `llm_server`: the process name
- `12345`: the process ID (PID)
- `5`, `14`: number of times the process called `mmap()`

#### 🧠 Use Case:
Model files (like `.pt`, `.bin`) are usually memory-mapped using `mmap()`. This shows which processes are mapping large files — useful for understanding cold start behavior and I/O footprint.

---

### 📄 mlock.bt

#### 🔹 Output:
```
@[torchrun, 9123]: 3
@[vllm_worker, 12077]: 6
```

#### 💡 What it means:
- Shows how many times each process invoked `mlock()` (to lock memory pages)

#### 🧠 Use Case:
Inference engines pin memory to avoid page swapping. This helps you detect memory pinning behavior and spot excessive locked memory usage.

---

### 📄 sched.bt

#### 🔹 Output:
```
@["llm_worker", "kworker"]: 102
@["vllm_server", "python3"]: 28
```

#### 💡 What it means:
- Represents (prev_comm, next_comm): task switch pairs
- Shows how frequently thread switches occurred

#### 🧠 Use Case:
Reveals inference thread behavior — helpful for:
- CPU core pinning validation
- Detecting noisy neighbors
- Optimizing latency

---

### 📄 gpu_access.bt

#### 🔹 Output:
```
@[vllm_worker, 14321]: 4
@[python3, 15432]: 2
```

#### 💡 What it means:
- Tracks how often processes accessed GPU-related device files like `/dev/nvidia0`

#### 🧠 Use Case:
Useful to:
- Confirm actual GPU usage
- Identify idle or blocked GPU tasks
- Debug multi-tenant conflicts

---

## 🚀 Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/llm-ebpf-tracer.git
cd llm-ebpf-tracer
```

### 2. Run a Sample Trace (Requires Root)

```bash
sudo bpftrace scripts/mlock.bt
```

### 3. Optional: Install Prometheus + eBPF Exporter

For live dashboards and metrics. (Coming soon.)

---

## 🙋‍♂️ Contributions Welcome

We're looking for collaborators who are into:

- PyTorch, HuggingFace, or vLLM internals
- Kernel memory and NUMA optimization
- eBPF, tracing, and observability pipelines

Open an issue, drop a PR, or just DM me on LinkedIn!

---

## 📬 Author

**Ravikumar Vallabhu**  
Linux Kernel + AI Infra Enthusiast  
[LinkedIn →](https://www.linkedin.com/in/ravikumar-vallabhu-62b7518/)  
[GitHub →](https://github.com/ravikumar1907/llm-ebpf-tracer)

---

## 📘 License

MIT License — Free to use, modify, and share.

---