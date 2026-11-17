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

// SaveCard persiste una tarjeta en la base de datos MySQL
func (r *MySQLBoardRepository) SaveCard(id, columnID, title, description, priority string) error {
	query := "INSERT INTO cards (id, column_id, title, description, priority) VALUES (?, ?, ?, ?, ?)"
	_, err := r.db.Exec(query, id, columnID, title, description, priority)
	if err != nil {
		return fmt.Errorf("error saving card: %w", err)
	}
	return nil
}
