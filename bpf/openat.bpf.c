#include "vmlinux.h"
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_tracing.h>
#include <bpf/bpf_core_read.h>

char LICENSE[] SEC("license") = "GPL";

struct syscall_event {
    __u32 pid;
    __u32 flags;       // O_RDONLY, O_CREAT, etc.
    __u16 mode;        // File permissions (if O_CREAT)
    int retval;        // Return value (success/failure)
    char comm[16];     // Process name
    char filename[128]; // File being opened
};

struct {
    __uint(type, BPF_MAP_TYPE_PERF_EVENT_ARRAY);
    __uint(key_size, sizeof(__u32));
    __uint(value_size, sizeof(__u32));
} events SEC(".maps");

struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(key_size, sizeof(__u32));
    __uint(value_size, sizeof(struct syscall_event));
    __uint(max_entries, 1024);
} start_events SEC(".maps");

// Kprobe: enter openat2
SEC("kprobe/do_sys_openat2")
int BPF_KPROBE(do_sys_openat2, int dfd, const char *filename, struct open_how *how) {
    struct syscall_event evt = {};
    evt.pid = bpf_get_current_pid_tgid() >> 32;
    bpf_get_current_comm(&evt.comm, sizeof(evt.comm));

    // Read open_how struct (flags, mode)
    if (how) {
        evt.flags = BPF_CORE_READ(how, flags);
        evt.mode  = BPF_CORE_READ(how, mode);
    }

    // Safely read filename (user-space pointer)
    bpf_probe_read_user_str(evt.filename, sizeof(evt.filename), filename);

    // Store in map to access in kretprobe
    bpf_map_update_elem(&start_events, &evt.pid, &evt, BPF_ANY);
    return 0;
}

// Kretprobe: exit openat2 (capture return value)
SEC("kretprobe/do_sys_openat2")
int BPF_KRETPROBE(do_sys_openat2_exit, int ret) {
    __u32 pid = bpf_get_current_pid_tgid() >> 32;
    struct syscall_event *evt = bpf_map_lookup_elem(&start_events, &pid);
    if (!evt) return 0;

    evt->retval = ret;
    bpf_perf_event_output(ctx, &events, BPF_F_CURRENT_CPU, evt, sizeof(*evt));
    bpf_map_delete_elem(&start_events, &pid);
    return 0;
}