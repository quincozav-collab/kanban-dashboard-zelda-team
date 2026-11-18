package infrastructure

import (
	"encoding/json"
	"kanban-dashboard/internal/domain/board"
	"net/http"
)

type BoardHandler struct {
	repo *MySQLBoardRepository
}

func NewBoardHandler(repo *MySQLBoardRepository) *BoardHandler {
	return &BoardHandler{repo: repo}
}

// CreateBoardEndpoint maneja la creación de un nuevo tablero vía POST
func (h *BoardHandler) CreateBoardEndpoint(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ID    string `json:"id"`
		Title string `json:"title"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Error al leer los datos", http.StatusBadRequest)
		return
	}

	b := &board.Board{
		ID:    req.ID,
		Title: req.Title,
	}

	err = h.repo.Save(b)
	if err != nil {
		http.Error(w, "Error al guardar en la base de datos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Tablero creado exitosamente 🚀"})
}

// CreateCardEndpoint maneja la creación de una tarjeta vía POST
func (h *BoardHandler) CreateCardEndpoint(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ID          string `json:"id"`
		ColumnID    string `json:"column_id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Priority    string `json:"priority"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Error al leer los datos de la tarjeta", http.StatusBadRequest)
		return
	}

	err = h.repo.SaveCard(req.ID, req.ColumnID, req.Title, req.Description, req.Priority)
	if err != nil {
		http.Error(w, "Error al guardar la tarjeta en la base de datos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Tarjeta creada exitosamente 🚀"})
}

// Endpoint para listar todos los tableros
func (h *BoardHandler) GetBoardsEndpoint(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	boards, err := h.repo.GetAllBoards()
	if err != nil {
		http.Error(w, "Error al obtener los tableros", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(boards)
}

// Endpoint para listar todas las tarjetas
func (h *BoardHandler) GetCardsEndpoint(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	cards, err := h.repo.GetAllCards()
	if err != nil {
		http.Error(w, "Error al obtener las tarjetas", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cards)
}
