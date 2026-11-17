package application

import (
	"errors"
	"kanban-dashboard/internal/domain/board"
)

type BoardService struct {
	// Aquí inyectaremos después los repositorios de base de datos
}

func NewBoardService() *BoardService {
	return &BoardService{}
}

// MoveCardUseCase ejecuta la regla de negocio para mover una tarjeta de columna
func (s *BoardService) MoveCardUseCase(b *board.Board, targetColumnID string) error {
	err := b.MoveCard(targetColumnID)
	if err != nil {
		if errors.Is(err, board.ErrWipLimitExceeded) {
			return err
		}
		return err
	}
	return nil
}
