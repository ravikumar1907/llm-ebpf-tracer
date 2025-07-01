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

		comm := strings.TrimRight(string(evt.Comm[:]), "\x00")
		filename := strings.TrimRight(string(evt.Filename[:]), "\x00")
		fmt.Printf("openat: PID=%d, Comm=%s, Flags=0x%x, Mode=0%o, Retval=%d, Filename=%s\n",
			evt.Pid, comm, evt.Flags, evt.Mode, evt.Retval, filename)

		api.Increment("openat", evt.Comm)
	}
}
