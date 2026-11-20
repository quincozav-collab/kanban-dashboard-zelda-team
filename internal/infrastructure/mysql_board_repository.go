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

type BoardMetrics struct {
	AvgLeadTimeHours  float64        `json:"avg_lead_time_hours"`
	AvgCycleTimeHours float64        `json:"avg_cycle_time_hours"`
	CompletedCards    int            `json:"completed_cards"`
	CardsByColumn     map[string]int `json:"cards_by_column"`
}

// MoveCard actualiza la columna asignada a una tarjeta existente y registra marcas de tiempo
func (r *MySQLBoardRepository) MoveCard(cardID string, targetColumnID string) error {
	query := `
		UPDATE cards
		SET column_id = ?,
			started_at = CASE
				WHEN ? IN ('col-2', 'col-3', 'col-4') AND started_at IS NULL THEN NOW()
				ELSE started_at
			END,
			completed_at = CASE
				WHEN ? IN ('col-3', 'col-4') AND completed_at IS NULL THEN NOW()
				ELSE completed_at
			END
		WHERE id = ?
	`
	result, err := r.db.Exec(query, targetColumnID, targetColumnID, targetColumnID, cardID)
	if err != nil {
		return fmt.Errorf("error moving card: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("card not found")
	}

	return nil
}

// GetMetrics calcula el conteo por columna y el promedio de Lead Time y Cycle Time en horas
func (r *MySQLBoardRepository) GetMetrics() (BoardMetrics, error) {
	metrics := BoardMetrics{
		CardsByColumn: make(map[string]int),
	}

	rows, err := r.db.Query(`SELECT column_id, COUNT(*) FROM cards GROUP BY column_id`)
	if err != nil {
		return metrics, fmt.Errorf("error counting cards by column: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var colID string
		var count int
		if err := rows.Scan(&colID, &count); err != nil {
			return metrics, err
		}
		metrics.CardsByColumn[colID] = count
	}
	if err := rows.Err(); err != nil {
		return metrics, fmt.Errorf("error iterating cards by column: %w", err)
	}

	metricQuery := `
		SELECT
			COUNT(*),
			COALESCE(AVG(TIMESTAMPDIFF(MINUTE, created_at, completed_at) / 60.0), 0),
			COALESCE(AVG(TIMESTAMPDIFF(MINUTE, started_at, completed_at) / 60.0), 0)
		FROM cards
		WHERE completed_at IS NOT NULL
	`
	if err := r.db.QueryRow(metricQuery).Scan(
		&metrics.CompletedCards,
		&metrics.AvgLeadTimeHours,
		&metrics.AvgCycleTimeHours,
	); err != nil {
		return metrics, fmt.Errorf("error calculating metrics: %w", err)
	}

	return metrics, nil
}
