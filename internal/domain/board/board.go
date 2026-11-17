package board

import "errors"

// ErrWipLimitExceeded se lanza cuando una columna ya alcanzó su límite máximo de tarjetas.
var ErrWipLimitExceeded = errors.New("WIP limit exceeded for this column")

type Column struct {
	ID            string
	Title         string
	WipLimit      int
	CurrentCards  int
}

// CanAddCard valida si la columna puede recibir otra tarjeta según su límite WIP.
func (c *Column) CanAddCard() bool {
	if c.WipLimit <= 0 {
		return true // Sin límite configurado
	}
	return c.CurrentCards < c.WipLimit
}