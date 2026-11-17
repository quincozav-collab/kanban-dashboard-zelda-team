package card

import "time"

type Card struct {
	ID          string
	ColumnID    string
	Title       string
	Description string
	Priority    string
	DueDate     time.Time
	CreatedAt   time.Time
}

// NewCard crea una nueva tarjeta con los datos básicos necesarios.
func NewCard(id, columnID, title, description, priority string, dueDate time.Time) *Card {
	return &Card{
		ID:          id,
		ColumnID:    columnID,
		Title:       title,
		Description: description,
		Priority:    priority,
		DueDate:     dueDate,
		CreatedAt:   time.Now(),
	}
}
