package events

type MlockEvent struct {
	Pid  uint32
	Comm [16]byte
}
