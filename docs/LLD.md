# Low-Level Design (LLD): LLM eBPF Tracer

## 🔧 Tracepoints

| Event              | Tracepoint             | Description                          |
|--------------------|------------------------|--------------------------------------|
| mmap               | sys_enter_mmap         | Monitor model weight loading         |
| mlock              | sys_enter_mlock        | Detect memory pinning                |
| openat             | sys_enter_openat       | Track GPU device file access         |
| scheduling         | sched:sched_switch     | Analyze inference thread behavior    |

## 🧪 eBPF Code (C)

Each probe is defined in `.bpf.c` files, compiled via clang/llvm:

```c
SEC("tracepoint/syscalls/sys_enter_mmap")
int trace_mmap(struct trace_event_raw_sys_enter *ctx) {
    u64 pid = bpf_get_current_pid_tgid();
    bpf_printk("mmap by PID %d\n", pid);
    return 0;
}
```

## ⚙ Go Components

### main.go

- Initialize BPF spec from embedded assets
- Attach to selected tracepoints
- Read events from perf buffer
- Serialize to stdout or metrics exporter

### events/

- Structs for each tracepoint (e.g., `MmapEvent`)
- JSON-serializable format
- Optionally enrich with `comm`, `pid`, `tid`

## 🗂 Folder Layout

```
/cmd
  - main.go

/internal/tracer
  - loader.go       # BPF loader/attacher
  - reader.go       # perf/ringbuffer reader

/internal/events
  - mmap.go
  - mlock.go
  - sched.go

/bpf
  - mmap.bpf.c
  - mlock.bpf.c
  - openat.bpf.c

/api (optional)
  - metrics.go      # Prometheus exporter

/scripts
  - simulate_inference.sh

/docs
  - HLD.md
  - LLD.md
```

## 📊 Output Format

```json
{
  "pid": 12345,
  "comm": "python3",
  "syscall": "mmap",
  "timestamp": 1682312320000000
}
```