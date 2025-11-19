package logging

type Loggable interface {
	GetLoggableID() uint64
	GetLoggableType() string
}
