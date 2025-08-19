package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq" // Importa o driver PostgreSQL
)

func main() {
	// Variáveis de ambiente (padrão docker-compose)
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	dbname := getEnv("DB_NAME", "projectlab")

	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	// Conexão com o banco
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Erro ao conectar no banco: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Erro ao pingar o banco: %v", err)
	}

	fmt.Println("✅ Conectado ao PostgreSQL com sucesso!")

	// Criar tabelas automaticamente
	createTables(db)

	// Rota de teste
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "🚀 Travel Planner API rodando!")
	})

	fmt.Println("Servidor rodando em http://localhost:8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

// Função para criar tabelas se não existirem
func createTables(db *sql.DB) {
	scripts := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			email VARCHAR(120) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS groups (
			id SERIAL PRIMARY KEY,
			name VARCHAR(150) NOT NULL,
			start_date DATE NOT NULL,
			end_date DATE NOT NULL,
			owner_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS group_members (
			user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			group_id INT NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
			role VARCHAR(50) DEFAULT 'member',
			joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (user_id, group_id)
		);`,
		`CREATE TABLE IF NOT EXISTS destinations (
			id SERIAL PRIMARY KEY,
			group_id INT NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
			name VARCHAR(150) NOT NULL,
			location VARCHAR(150) NOT NULL,
			description TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS votes (
			id SERIAL PRIMARY KEY,
			group_id INT NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
			title VARCHAR(150) NOT NULL,
			description TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			deadline TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS vote_options (
			id SERIAL PRIMARY KEY,
			vote_id INT NOT NULL REFERENCES votes(id) ON DELETE CASCADE,
			option_text VARCHAR(150) NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS vote_responses (
			id SERIAL PRIMARY KEY,
			vote_id INT NOT NULL REFERENCES votes(id) ON DELETE CASCADE,
			option_id INT NOT NULL REFERENCES vote_options(id) ON DELETE CASCADE,
			user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			voted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE (vote_id, user_id)
		);`,
		`CREATE TABLE IF NOT EXISTS expenses (
			id SERIAL PRIMARY KEY,
			group_id INT NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
			user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			description VARCHAR(200) NOT NULL,
			amount NUMERIC(10,2) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS tasks (
			id SERIAL PRIMARY KEY,
			group_id INT NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
			title VARCHAR(150) NOT NULL,
			description TEXT,
			assigned_to INT REFERENCES users(id) ON DELETE SET NULL,
			is_completed BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
	}

	for _, script := range scripts {
		_, err := db.Exec(script)
		if err != nil {
			log.Fatalf("Erro ao criar tabela: %v", err)
		}
	}

	fmt.Println("✅ Tabelas criadas ou já existentes!")
}

// Função auxiliar pra pegar variável de ambiente com fallback
func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}
