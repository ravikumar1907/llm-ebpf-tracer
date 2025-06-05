package events

type SchedSwitchEvent struct {
	PrevComm [16]byte
	NextComm [16]byte
	PrevPid  uint32
	NextPid  uint32
}
