package db

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite" // Driver SQLite puro Go (sem CGO)
)

// DB encapsula a conexão com o banco de dados
type DB struct {
	*sqlx.DB
	dsn string
}

// Config contém as configurações de conexão com o banco
type Config struct {
	DSN            string
	MigrationsPath string
}

// NewDB cria uma nova conexão com o banco de dados
// Lê a DSN da variável de ambiente DATABASE_URL ou usa SQLite local como padrão
func NewDB() (*DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// Criar diretório data se não existir
		dataDir := "./data"
		if err := os.MkdirAll(dataDir, 0755); err != nil {
			return nil, fmt.Errorf("erro ao criar diretório de dados: %w", err)
		}
		dsn = "file:./data/rpg.db?cache=shared&mode=rwc"
	}

	return NewDBWithDSN(dsn)
}

// NewDBWithDSN cria uma nova conexão com o banco usando uma DSN específica
func NewDBWithDSN(dsn string) (*DB, error) {
	// Determinar o driver baseado na DSN
	driver := determineDriver(dsn)

	sqlxDB, err := sqlx.Connect(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar com o banco de dados: %w", err)
	}

	// Configurar pool de conexões
	sqlxDB.SetMaxOpenConns(25)
	sqlxDB.SetMaxIdleConns(10)

	db := &DB{
		DB:  sqlxDB,
		dsn: dsn,
	}

	return db, nil
}

// determineDriver determina o driver baseado na DSN
func determineDriver(dsn string) string {
	// Lógica simples para determinar o driver
	// Pode ser expandida para suportar outros bancos
	if len(dsn) > 0 && (dsn[0:5] == "file:" || dsn[0:1] == "/" || dsn[0:2] == "./" || dsn[0:3] == "../") {
		return "sqlite" // modernc.org/sqlite usa "sqlite" como nome do driver
	}

	// Por enquanto, assumir SQLite como padrão
	// Futuro: adicionar suporte para PostgreSQL, MySQL, etc.
	return "sqlite"
}

// RunMigrations executa as migrações automaticamente na inicialização
func (db *DB) RunMigrations() error {
	migrationsPath := "./migrations"

	// Verificar se o diretório de migrações existe
	if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
		// Criar diretório de migrações se não existir
		if err := os.MkdirAll(migrationsPath, 0755); err != nil {
			return fmt.Errorf("erro ao criar diretório de migrações: %w", err)
		}
		return nil // Sem migrações para executar
	}

	// Configurar goose
	driver := determineDriver(db.dsn)
	if err := goose.SetDialect(driver); err != nil {
		return fmt.Errorf("erro ao configurar dialeto do goose: %w", err)
	}

	// Executar migrações
	if err := goose.Up(db.DB.DB, migrationsPath); err != nil {
		return fmt.Errorf("erro ao executar migrações: %w", err)
	}

	return nil
}

// Close fecha a conexão com o banco de dados
func (db *DB) Close() error {
	return db.DB.Close()
}

// Health verifica se a conexão com o banco está saudável
func (db *DB) Health() error {
	return db.DB.Ping()
}

// GetDSN retorna a DSN sendo usada pela conexão
func (db *DB) GetDSN() string {
	return db.dsn
}
