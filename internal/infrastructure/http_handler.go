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

func (h *BoardHandler) CreateCardEndpoint(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var cardData struct {
		ID          string `json:"id"`
		ColumnID    string `json:"column_id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Priority    string `json:"priority"`
	}

	if err := json.NewDecoder(r.Body).Decode(&cardData); err != nil {
		http.Error(w, "Entrada inválida", http.StatusBadRequest)
		return
	}

	// --- VALIDACIÓN DE LÍMITE WIP ---
	// Definimos un límite máximo de 3 tarjetas por columna para este ejemplo
	const maxWIPLimit = 3
	currentCount, err := h.repo.CountCardsByColumn(cardData.ColumnID)
	if err != nil {
		http.Error(w, "Error al validar el límite WIP", http.StatusInternalServerError)
		return
	}

	if currentCount >= maxWIPLimit {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict) // Código HTTP 409 Conflict
		w.Write([]byte(`{"error": "Límite WIP excedido: la columna ya alcanzó el máximo de tarjetas permitidas"}`))
		return
	}
	// ---------------------------------

	err = h.repo.SaveCard(cardData.ID, cardData.ColumnID, cardData.Title, cardData.Description, cardData.Priority)
	if err != nil {
		http.Error(w, "Error al guardar la tarjeta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "Tarjeta creada exitosamente bajo las reglas WIP"}`))
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
