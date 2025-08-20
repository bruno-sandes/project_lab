package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Conectar no banco
	db := connectDB()
	defer db.Close()

	// Garantir tabelas
	createTables(db)

	// Rota de teste
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "🚀 Project Lab API rodando!")
	})

	fmt.Println("Servidor rodando em http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
