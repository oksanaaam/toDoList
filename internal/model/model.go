package model

type ToDo struct {
	ID     string `json:"id,omitempty" bson:"_id,omitempty"`
	Title  string `json:"title"`
	Status Status `json:"status"`
}

// IsValidStatus checks, if status is valid
func IsValidStatus(status Status) bool {
	switch status {
	case Created, InProgress, Done:
		return true
	}
	return false
}
