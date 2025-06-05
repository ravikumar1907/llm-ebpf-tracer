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

	tp, err := link.Tracepoint("syscalls", "sys_enter_mmap", objs.TraceMmap, nil)
	if err != nil {
		fmt.Printf("linking tracepoint: %v", err)
	}
	defer tp.Close()

	log.Println("BPF program loaded and tracepoint linked successfully.")

	go tracer.ReadEvents(objs.Events)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	fmt.Println("Listening for mmap events... Press Ctrl+C to stop.")
	<-ctx.Done()
}
