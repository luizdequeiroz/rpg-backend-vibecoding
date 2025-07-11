package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/luizdequeiroz/rpg-backend/pkg/db"
)

func main() {
	var (
		action = flag.String("action", "", "Ação a ser executada: create, up, down, status, reset")
		name   = flag.String("name", "", "Nome da migração (usado com create)")
	)
	flag.Parse()

	if *action == "" {
		fmt.Println("Uso: go run cmd/migrate/main.go -action=<create|up|down|status|reset> [-name=<nome_da_migracao>]")
		fmt.Println("\nAções disponíveis:")
		fmt.Println("  create - Cria uma nova migração (requer -name)")
		fmt.Println("  up     - Executa todas as migrações pendentes")
		fmt.Println("  down   - Desfaz a última migração")
		fmt.Println("  status - Mostra o status das migrações")
		fmt.Println("  reset  - Remove todas as migrações aplicadas")
		os.Exit(1)
	}

	// Conectar ao banco de dados
	database, err := db.NewDB()
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	defer database.Close()

	// Criar gerenciador de migrações
	migrationManager := db.NewMigrationManager(database, "./migrations")

	// Executar ação solicitada
	switch *action {
	case "create":
		if *name == "" {
			log.Fatal("Nome da migração é obrigatório para criar uma nova migração")
		}
		if err := migrationManager.CreateMigration(*name); err != nil {
			log.Fatalf("Erro ao criar migração: %v", err)
		}
		fmt.Printf("Migração '%s' criada com sucesso!\n", *name)

	case "up":
		if err := migrationManager.Up(); err != nil {
			log.Fatalf("Erro ao executar migrações up: %v", err)
		}
		fmt.Println("Migrações executadas com sucesso!")

	case "down":
		if err := migrationManager.Down(); err != nil {
			log.Fatalf("Erro ao executar migração down: %v", err)
		}
		fmt.Println("Migração desfeita com sucesso!")

	case "status":
		if err := migrationManager.Status(); err != nil {
			log.Fatalf("Erro ao verificar status das migrações: %v", err)
		}

	case "reset":
		if err := migrationManager.Reset(); err != nil {
			log.Fatalf("Erro ao resetar migrações: %v", err)
		}
		fmt.Println("Todas as migrações foram removidas!")

	default:
		log.Fatalf("Ação desconhecida: %s", *action)
	}
}
