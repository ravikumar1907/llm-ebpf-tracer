package events

type MmapEvent struct {
	Pid  uint32
	Comm [16]byte
}
