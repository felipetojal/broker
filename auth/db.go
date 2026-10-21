package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// Database encapsula uma conexão o banco de dados.
type Database struct {
	db *sql.DB
}

var (
	ErrorConnectingToDatabase = errors.New("erro ao conectar ao banco de dados")
)

const (
	MaxOpenConns    = 40
	MaxIdleConns    = 40
	MaxConnLifetime = 3 * time.Minute
)

// NewDatabase cria uma nova instância de Database e se conecta ao banco de dados.
func NewDatabase(name, password, user, host, port string) (*Database, error) {
	// Monte a string de conexão (ajuste o host, porta e usuário conforme seu ambiente)
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, name)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir conexão com o banco: %w", err)
	}

	// Configuração do pool de conexões (essencial para produção)
	db.SetMaxOpenConns(MaxOpenConns)
	db.SetMaxIdleConns(MaxIdleConns)
	db.SetConnMaxLifetime(MaxConnLifetime)

	// Valida se a conexão com o banco realmente funciona
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("%w: %w", ErrorConnectingToDatabase, err)
	}

	return &Database{db: db}, nil
}

// Close fecha a conexão com o banco de dados.
func (d *Database) Close() error {
	return d.db.Close()
}
