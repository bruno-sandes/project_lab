package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func connectDB() *sql.DB {
	// Pega configs do docker-compose (com fallback)
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "admin")
	password := getEnv("DB_PASSWORD", "admin123")
	dbname := getEnv("DB_NAME", "project_lab")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Erro ao abrir conexão com o banco:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Banco não respondeu:", err)
	}

	fmt.Println("Conexão bem sucedida com o Postgres!")
	return db
}

// Função auxiliar
func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}
