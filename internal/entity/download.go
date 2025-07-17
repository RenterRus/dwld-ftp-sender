package entity

type Status int

const (
	TO_SEND Status = iota + 1
	SENDING
	DONE
)

var StatusMapping map[Status]string = map[Status]string{
	TO_SEND: "TO_SEND",
	SENDING: "SENDING",
	DONE:    "DONE",
}
