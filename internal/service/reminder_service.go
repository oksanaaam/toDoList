package service

import (
	"log"
	"time"
)

type Reminder struct {
	ID           string
	ReminderTime time.Time
	TaskName     string
}

type ReminderService struct {
	reminderChannel     chan Reminder // Channel for new reminders
	stopChannel         chan struct{} // Channel for stopping goroutines
	notificationChannel chan string   // Channel for notifications
	reminders           []Reminder    // Slice for storing all reminders
}

// NewReminderService creates a new service for working with reminders
func NewReminderService(notificationChannel chan string) *ReminderService {
	return &ReminderService{
		reminderChannel:     make(chan Reminder),
		stopChannel:         make(chan struct{}),
		notificationChannel: notificationChannel,
		reminders:           []Reminder{},
	}
}

func (rs *ReminderService) StartWorker() {
	go func() {
		ticker := time.NewTicker(1 * time.Second) // Check every second
		defer ticker.Stop()

		for {
			select {
			case reminder := <-rs.reminderChannel:
				rs.reminders = append(rs.reminders, reminder)
				log.Printf("Reminder received: %v for task '%s'\n", reminder.ReminderTime, reminder.TaskName)
			case <-ticker.C:
				for _, reminder := range rs.reminders {
					if time.Now().After(reminder.ReminderTime) {
						rs.notificationChannel <- "You need to do this task: " + reminder.TaskName
					}
				}
			case <-rs.stopChannel:
				log.Println("Reminder worker is stopping...")
				return
			}
		}
	}()
}

func (rs *ReminderService) AddReminder(reminder Reminder) {
	log.Printf("Adding reminder for task id '%s' with time '%v'\n", reminder.ID, reminder.ReminderTime)
	rs.reminderChannel <- reminder
}

func (rs *ReminderService) StopWorker() {
	close(rs.stopChannel)
}
