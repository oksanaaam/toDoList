package model

type Status string

const (
	Created    Status = "created"
	InProgress Status = "in progress"
	Done       Status = "done"
)
