package api

import (
	"database/sql"
	"strconv"
)

type ReminderStorer interface {
	GetReminders() ([]*Reminder, error)
	CreateReminder(*Reminder) (*Reminder, error)
	UpdateReminder(id string, reminder *Reminder) (*Reminder, error)
	DeleteReminder(id string) error
}

type MySQLService struct {
	DB *sql.DB
}

func NewMySQLService(DB *sql.DB) *MySQLService {
	return &MySQLService{
		DB: DB,
	}
}

func (s *MySQLService) GetReminders() ([]*Reminder, error) {
	rows, err := s.DB.Query("SELECT * FROM reminders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reminders []*Reminder

	for rows.Next() {
		var reminderID int64
		var reminderTitle string
		var reminderDescription string

		err = rows.Scan(&reminderID, &reminderTitle, &reminderDescription)
		if err != nil {
			return nil, err
		}

		reminders = append(reminders, &Reminder{
			ID:          reminderID,
			Title:       reminderTitle,
			Description: reminderDescription,
		})
	}

	return reminders, nil
}

func (s *MySQLService) CreateReminder(reminder *Reminder) (*Reminder, error) {
	const INSERT_REMINDER = "INSERT INTO reminders (reminder_title, reminder_description) VALUES (?, ?)"

	result, err := s.DB.Exec(INSERT_REMINDER, reminder.Title, reminder.Description)
	if err != nil {
		return nil, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	reminder.ID = lastId
	return reminder, nil
}

func (s *MySQLService) UpdateReminder(id string, reminder *Reminder) (*Reminder, error) {
	const UPDATE_REMINDER = "UPDATE reminders SET reminder_title = ?, reminder_description = ? WHERE reminder_id = ?"

	_, err := s.DB.Exec(UPDATE_REMINDER, reminder.Title, reminder.Description, id)
	if err != nil {
		return nil, err
	}

	num, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, err
	}

	reminder.ID = num
	return reminder, nil
}

func (s *MySQLService) DeleteReminder(id string) error {
	const UPDATE_REMINDER = "DELETE FROM reminders WHERE reminder_id = ?"
	_, err := s.DB.Exec(UPDATE_REMINDER, id)
	return err
}
