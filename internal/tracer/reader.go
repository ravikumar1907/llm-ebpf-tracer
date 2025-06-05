package tracer

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/perf"

	"github.com/llmebpftracer/api"
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

		var evt events.SyscallEvent
		if err := binary.Read(bytes.NewBuffer(record.RawSample), binary.LittleEndian, &evt); err != nil {
			log.Printf("decode error: %v", err)
			continue
		}

		syscall := strings.TrimRight(string(evt.Syscall[:]), "\x00")
		comm := strings.TrimRight(string(evt.Comm[:]), "\x00")
		fmt.Printf("%s: PID=%d, Comm=%s\n", syscall, evt.Pid, comm)

		api.Increment(syscall, evt.Comm)

	}
}
