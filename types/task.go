package types

const (
	TASK_TYPE_ACM_C = 0x0
)

type Task struct {
	Buf  *Buffer
	Type int
}
