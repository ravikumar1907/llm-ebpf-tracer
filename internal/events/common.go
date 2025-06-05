package events

type SyscallEvent struct {
	Pid     uint32
	Comm    [16]byte
	Syscall [6]byte
}
