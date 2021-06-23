package goueue

//Command is a simple interface that implements command pattern
type Command interface {
	Handle() error
}
