package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"kanban-dashboard/internal/infrastructure"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Cadena de conexión para el contenedor MySQL que encendimos en Docker
	dsn := "kanban_user:kanban_password@tcp(127.0.0.1:3306)/kanban_db?parseTime=true"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error al conectar con la base de datos: %v", err)
	}
	defer db.Close()

	// Probamos que la conexión responda
	err = db.Ping()
	if err != nil {
		log.Fatalf("No se pudo conectar a MySQL en Docker: %v", err)
	}

	fmt.Println("¡Conexión exitosa a MySQL 8 en Docker!")

	// Inicializamos el repositorio
	repo := infrastructure.NewMySQLBoardRepository(db)
	_ = repo

	// Endpoint de prueba rápido
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("API Kanban funcionando al 100% 🚀"))
	})

	fmt.Println("Servidor corriendo en el puerto 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
