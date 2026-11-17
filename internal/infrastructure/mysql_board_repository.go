package infrastructure

import (
	"database/sql"
	"fmt"
	"kanban-dashboard/internal/domain/board"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLBoardRepository struct {
	db *sql.DB
}

func NewMySQLBoardRepository(db *sql.DB) *MySQLBoardRepository {
	return &MySQLBoardRepository{db: db}
}

// Save persiste un tablero en MySQL 8
func (r *MySQLBoardRepository) Save(b *board.Board) error {
	query := "INSERT INTO boards (id, title) VALUES (?, ?)"
	_, err := r.db.Exec(query, b.ID, b.Title)
	if err != nil {
		return fmt.Errorf("error saving board: %w", err)
	}
	return nil
}
