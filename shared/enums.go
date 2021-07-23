package shared

type Status string

const (
	Initialized Status = "initialized"
	Applied     Status = "applied"
	Destroyed   Status = "destroyed"
)
