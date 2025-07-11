package main

import (
	"fmt"
	"log"

	"github.com/luizdequeiroz/rpg-backend/pkg/db"
)

func main() {
	// Conectar ao banco diretamente
	database, err := db.NewDB()
	if err != nil {
		log.Fatal("Erro ao conectar ao banco:", err)
	}
	defer database.Close()

	// Verificar se a tabela existe
	var count int
	err = database.Get(&count, "SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='sheet_templates'")
	if err != nil {
		log.Fatal("Erro ao verificar tabela:", err)
	}

	fmt.Printf("Tabela sheet_templates existe: %v\n", count > 0)

	if count > 0 {
		// Obter schema da tabela
		var schema string
		err = database.Get(&schema, "SELECT sql FROM sqlite_master WHERE type='table' AND name='sheet_templates'")
		if err != nil {
			log.Fatal("Erro ao obter schema:", err)
		}
		fmt.Printf("Schema da tabela:\n%s\n", schema)

		// Contar registros
		var total int
		err = database.Get(&total, "SELECT COUNT(*) FROM sheet_templates")
		if err != nil {
			log.Fatal("Erro ao contar registros:", err)
		}
		fmt.Printf("Total de registros: %d\n", total)
	}
}
