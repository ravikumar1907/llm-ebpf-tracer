package events

type SyscallEvent struct {
	Pid      uint32
	Flags    uint32
	Mode     uint16
	Retval   int32
	Comm     [16]byte
	Filename [128]byte
}
