package infrastructure

import (
	"database/sql"
	"fmt"
	"kanban-dashboard/internal/domain/board"
	"kanban-dashboard/internal/domain/card"

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

// GetAllBoards obtiene todos los tableros de la base de datos
func (r *MySQLBoardRepository) GetAllBoards() ([]board.Board, error) {
	rows, err := r.db.Query("SELECT id, title FROM boards")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var boards []board.Board
	for rows.Next() {
		var b board.Board
		if err := rows.Scan(&b.ID, &b.Title); err != nil {
			return nil, err
		}
		boards = append(boards, b)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return boards, nil
}

// GetAllCards obtiene todas las tarjetas de la base de datos
func (r *MySQLBoardRepository) GetAllCards() ([]card.Card, error) {
	rows, err := r.db.Query("SELECT id, column_id, title, description, priority, created_at FROM cards")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []card.Card
	for rows.Next() {
		var c card.Card
		if err := rows.Scan(&c.ID, &c.ColumnID, &c.Title, &c.Description, &c.Priority, &c.CreatedAt); err != nil {
			return nil, err
		}
		cards = append(cards, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return cards, nil
}

// CountCardsByColumn cuenta cuántas tarjetas existen en una columna específica
func (r *MySQLBoardRepository) CountCardsByColumn(columnID string) (int, error) {
	query := "SELECT COUNT(*) FROM cards WHERE column_id = ?"
	var count int
	err := r.db.QueryRow(query, columnID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error counting cards by column: %w", err)
	}
	return count, nil
}
