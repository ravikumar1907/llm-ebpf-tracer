# 📌 LLM eBPF Tracer – Roadmap

This document outlines the planned features and enhancements to evolve the project into a full end-to-end system.

---

## Phase 1: Backend Completion

- [x] BPF programs for `mmap`, `mlock`, `openat`
- [x] Go-based loader and perf event reader
- [x] Raw decoding of events with structured output
- [ ] Export metrics via Prometheus (`/api/metrics.go`)
- [ ] Add CLI config or `.env` to enable/disable specific probes
- [ ] Optional: Add logging in structured format (e.g., JSONL)
- [ ] Optional: RBAC or API token for shared deployment

---

## 🖥️ Phase 2: Dashboard & UI

### 1. Grafana Integration (via Prometheus)
- [ ] Create `/dashboards/` folder with Grafana JSON panels
- [ ] Syscall count per `comm` (e.g., torch, docker, etc.)
- [ ] Alerts for excessive `mlock`, `mmap`
- [ ] Visualization of scheduler switches (future)

### 2. Optional Web UI (React/Next.js)
- [ ] Display live syscall stream
- [ ] Filter by PID / `comm`
- [ ] Timeline of workload activity

---

## 🚀 Phase 3: Packaging & Usage

- [ ] Dockerfile with backend + metrics export
- [ ] `make run` target: compile, load, start backend
- [ ] Script: `simulate_llm.py` to generate test events (torch-based)
- [ ] CLI flag `--simulate` to trigger mock data

---

## 🔮 Future Enhancements

- [ ] Uprobes on common LLM inference frameworks
- [ ] CPU Thread Scheduling Graph (from `sched_switch`)
- [ ] GPU access trace on `/dev/nvidia*`
- [ ] Support multi-host trace aggregation (via NATS or gRPC)

---

## 📁 Folder Structure Overview

```
/bpf/                  # Compiled BPF object files and C sources
/internal/tracer/     # Loader, Reader
/internal/events/     # Go structs for event decoding
/api/                 # Prometheus HTTP export
/dashboards/          # Grafana dashboard JSONs (WIP)
/docs/roadmap.md      # This file
```

---

## 🛠 In Progress

Track open items with GitHub issues using labels like `backend`, `dashboard`, `enhancement`, `help-wanted`.

