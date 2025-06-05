package tracer

import (
	"bytes"
	"embed"
	"fmt"

	"github.com/cilium/ebpf"
)

//go:embed bpf/compiled/*.o
var bpfObjectsFS embed.FS

type BPFObjects struct {
	TraceMmap  *ebpf.Program
	TraceMlock *ebpf.Program
	TraceOpen  *ebpf.Program
	Events     *ebpf.Map
}

func LoadBPFObjects(_ *ebpf.CollectionSpec, _ *ebpf.CollectionOptions) (*BPFObjects, error) {
	objs := &BPFObjects{}

	// Load mmap.bpf.o
	mmapBytes, err := bpfObjectsFS.ReadFile("bpf/compiled/mmap.bpf.o")
	if err != nil {
		return nil, fmt.Errorf("read mmap.o: %w", err)
	}
	specMmap, err := ebpf.LoadCollectionSpecFromReader(bytes.NewReader(mmapBytes))
	if err != nil {
		return nil, fmt.Errorf("load mmap.o spec: %w", err)
	}
	collMmap, err := ebpf.NewCollection(specMmap)
	if err != nil {
		return nil, fmt.Errorf("load mmap.o collection: %w", err)
	}
	objs.TraceMmap = collMmap.Programs["trace_mmap"]
	objs.Events = collMmap.Maps["events"]

	// Load mlock.bpf.o
	mlockBytes, err := bpfObjectsFS.ReadFile("bpf/compiled/mlock.bpf.o")
	if err != nil {
		return nil, fmt.Errorf("read mlock.o: %w", err)
	}
	specMlock, err := ebpf.LoadCollectionSpecFromReader(bytes.NewReader(mlockBytes))
	if err != nil {
		return nil, fmt.Errorf("load mlock.o spec: %w", err)
	}
	collMlock, err := ebpf.NewCollection(specMlock)
	if err != nil {
		return nil, fmt.Errorf("load mlock.o collection: %w", err)
	}
	objs.TraceMlock = collMlock.Programs["trace_mlock"]

	// Load openat.bpf.o
	openBytes, err := bpfObjectsFS.ReadFile("bpf/compiled/openat.bpf.o")
	if err != nil {
		return nil, fmt.Errorf("read openat.o: %w", err)
	}
	specOpen, err := ebpf.LoadCollectionSpecFromReader(bytes.NewReader(openBytes))
	if err != nil {
		return nil, fmt.Errorf("load openat.o spec: %w", err)
	}
	collOpen, err := ebpf.NewCollection(specOpen)
	if err != nil {
		return nil, fmt.Errorf("load openat.o collection: %w", err)
	}
	objs.TraceOpen = collOpen.Programs["trace_openat"]

	return objs, nil
}
