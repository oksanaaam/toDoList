package model

type ToDo struct {
	ID           string `json:"id,omitempty" bson:"_id,omitempty"`
	Title        string `json:"title" bson:"title"`
	Status       Status `json:"status" bson:"status"`
	ImagePath    string `json:"image_path,omitempty" bson:"image_path,omitempty"`
	ReminderTime string `json:"reminder_time,omitempty" bson:"reminder_time,omitempty"`
}

// IsValidStatus checks, if status is valid
func IsValidStatus(status Status) bool {
	switch status {
	case Created, InProgress, Done:
		return true
	}
	return false
}
