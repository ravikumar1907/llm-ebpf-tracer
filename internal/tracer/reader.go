package tracer

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/perf"

	"github.com/llmebpftracer/internal/events"
)

func ReadEvents(eventsMap *ebpf.Map) error {
	rd, err := perf.NewReader(eventsMap, os.Getpagesize())
	if err != nil {
		return fmt.Errorf("create perf reader: %w", err)
	}
	defer rd.Close()

	fmt.Println("Listening for mlock/mmap/openat events...")

	for {
		record, err := rd.Read()
		if err != nil {
			fmt.Printf("read error: %v", err)
			continue
		}

		var evt events.MlockEvent // You can switch to MmapEvent or OpenEvent as needed
		if err := binary.Read(bytes.NewBuffer(record.RawSample), binary.LittleEndian, &evt); err != nil {
			fmt.Printf("decode error: %v", err)
			continue
		}

		fmt.Printf("mlock: PID=%d, Comm=%s\n", evt.Pid, string(bytes.TrimRight(evt.Comm[:], "\x00")))
		var eventsMap events.MmapEvent // You can switch to MmapEvent or OpenEvent as needed
		if err := binary.Read(bytes.NewBuffer(record.RawSample), binary.LittleEndian, &eventsMap); err != nil {
			fmt.Printf("decode error: %v", err)
			continue
		}
		fmt.Printf("mmap: PID=%d, Comm=%s\n", eventsMap.Pid, string(bytes.TrimRight(evt.Comm[:], "\x00")))

	}
}
