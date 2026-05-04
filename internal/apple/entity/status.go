package entity

type Status int

const (
	Unknown Status = iota
	Success
	Error
)

func NewStatus(s string) Status {
	switch s {
	case "success":
		return Success
	case "error":
		return Error
	default:
		return Unknown
	}
}

func (s Status) String() string {
	switch s {
	case Success:
		return "success"
	case Error:
		return "error"
	default:
		return "unknown"
	}
}
