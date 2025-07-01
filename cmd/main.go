package main

import (
	"context"
	"fmt"
	"log"

	"os"
	"os/signal"
	"syscall"

	"github.com/llmebpftracer/internal/tracer"

	"github.com/cilium/ebpf/link"
)

func main() {
	objs, err := tracer.LoadBPFObjects(nil, nil)
	if err != nil {
		fmt.Printf("loading BPF objects: %v", err)
	}
	defer objs.TraceMmap.Close()
	defer objs.TraceOpen.Close()

	tp, err := link.Tracepoint("syscalls", "sys_enter_mmap", objs.TraceMmap, nil)
	if err != nil {
		fmt.Printf("linking tracepoint: %v", err)
	}
	defer tp.Close()

	// Attach kprobe and kretprobe for openat2
	kp, err := link.Kprobe("do_sys_openat2", objs.TraceOpen, nil)
	if err != nil {
		fmt.Printf("linking kprobe: %v", err)
	}
	defer kp.Close()

	krp, err := link.Kretprobe("do_sys_openat2", objs.TraceOpen, nil)
	if err != nil {
		fmt.Printf("linking kretprobe: %v", err)
	}
	defer krp.Close()

	log.Println("BPF programs loaded and tracepoint/kprobes linked successfully.")

	go tracer.ReadEvents(objs.Events)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	fmt.Println("Listening for mmap and openat events... Press Ctrl+C to stop.")
	<-ctx.Done()
}
