package db

import (
	"fmt"
	"os"

	"github.com/pressly/goose/v3"
)

// MigrationManager gerencia as operações de migração
type MigrationManager struct {
	db             *DB
	migrationsPath string
}

// NewMigrationManager cria um novo gerenciador de migrações
func NewMigrationManager(db *DB, migrationsPath string) *MigrationManager {
	if migrationsPath == "" {
		migrationsPath = "./migrations"
	}

	return &MigrationManager{
		db:             db,
		migrationsPath: migrationsPath,
	}
}

// CreateMigration cria uma nova migração com o nome especificado
func (mm *MigrationManager) CreateMigration(name string) error {
	// Criar diretório de migrações se não existir
	if err := os.MkdirAll(mm.migrationsPath, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório de migrações: %w", err)
	}

	// Configurar goose
	driver := determineDriver(mm.db.dsn)
	if err := goose.SetDialect(driver); err != nil {
		return fmt.Errorf("erro ao configurar dialeto do goose: %w", err)
	}

	// Criar migração
	if err := goose.Create(mm.db.DB.DB, mm.migrationsPath, name, "sql"); err != nil {
		return fmt.Errorf("erro ao criar migração: %w", err)
	}

	fmt.Printf("Migração criada com sucesso em: %s\n", mm.migrationsPath)
	return nil
}

// Up executa todas as migrações pendentes
func (mm *MigrationManager) Up() error {
	driver := determineDriver(mm.db.dsn)
	if err := goose.SetDialect(driver); err != nil {
		return fmt.Errorf("erro ao configurar dialeto do goose: %w", err)
	}

	if err := goose.Up(mm.db.DB.DB, mm.migrationsPath); err != nil {
		return fmt.Errorf("erro ao executar migrações up: %w", err)
	}

	return nil
}

// Down desfaz a última migração
func (mm *MigrationManager) Down() error {
	driver := determineDriver(mm.db.dsn)
	if err := goose.SetDialect(driver); err != nil {
		return fmt.Errorf("erro ao configurar dialeto do goose: %w", err)
	}

	if err := goose.Down(mm.db.DB.DB, mm.migrationsPath); err != nil {
		return fmt.Errorf("erro ao executar migração down: %w", err)
	}

	return nil
}

// Status mostra o status das migrações
func (mm *MigrationManager) Status() error {
	driver := determineDriver(mm.db.dsn)
	if err := goose.SetDialect(driver); err != nil {
		return fmt.Errorf("erro ao configurar dialeto do goose: %w", err)
	}

	if err := goose.Status(mm.db.DB.DB, mm.migrationsPath); err != nil {
		return fmt.Errorf("erro ao verificar status das migrações: %w", err)
	}

	return nil
}

// Reset remove todas as migrações aplicadas
func (mm *MigrationManager) Reset() error {
	driver := determineDriver(mm.db.dsn)
	if err := goose.SetDialect(driver); err != nil {
		return fmt.Errorf("erro ao configurar dialeto do goose: %w", err)
	}

	if err := goose.Reset(mm.db.DB.DB, mm.migrationsPath); err != nil {
		return fmt.Errorf("erro ao resetar migrações: %w", err)
	}

	return nil
}
