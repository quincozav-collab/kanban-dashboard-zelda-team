package board

type Board struct {
	ID      string
	Title   string
	Columns []Column
}

// MoveCard valida si se puede mover una tarjeta a una columna destino respetando el límite WIP.
func (b *Board) MoveCard(targetColumnID string) error {
	for i := range b.Columns {
		if b.Columns[i].ID == targetColumnID {
			// Aquí aplicamos la regla de oro del Kanban que exige tu profesor
			if !b.Columns[i].CanAddCard() {
				return ErrWipLimitExceeded
			}
			// Si pasa la validación, incrementamos temporalmente el contador de la columna
			b.Columns[i].CurrentCards++
			return nil
		}
	}
	return nil
}
