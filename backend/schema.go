package main

import (
	"database/sql"
	"fmt"
)

func createTables(db *sql.DB) {
	// Aqui você coloca todas as tabelas necessárias
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS projects (
		id SERIAL PRIMARY KEY,
		title VARCHAR(150) NOT NULL,
		description TEXT,
		owner_id INT REFERENCES users(id),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := db.Exec(query)
	if err != nil {
		panic(fmt.Sprintf("Erro ao criar tabelas: %v", err))
	}

	fmt.Println("Tabelas garantidas no banco!")
}
