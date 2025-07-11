package db

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDatabaseConnection(t *testing.T) {
	// Criar mock de banco
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	// Criar instância de DB com mock
	sqlxDB := sqlx.NewDb(mockDB, "sqlite3")
	db := &DB{DB: sqlxDB}

	// Configurar expectativa para o ping
	mock.ExpectPing()

	// Testar conexão
	err = db.Ping()
	assert.NoError(t, err)

	// Verificar se todas as expectativas foram atendidas
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestDatabaseHealthCheck(t *testing.T) {
	t.Run("Database healthy", func(t *testing.T) {
		// Criar mock de banco com ping monitoring habilitado
		mockDB, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
		require.NoError(t, err)
		defer mockDB.Close()

		sqlxDB := sqlx.NewDb(mockDB, "sqlite3")
		db := &DB{DB: sqlxDB}

		// Configurar expectativa para sucesso
		mock.ExpectPing()

		err = db.Health()
		assert.NoError(t, err)

		// Verificar expectativas
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("Database unhealthy", func(t *testing.T) {
		// Criar mock de banco com ping monitoring habilitado
		mockDB, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
		require.NoError(t, err)
		defer mockDB.Close()

		sqlxDB := sqlx.NewDb(mockDB, "sqlite3")
		db := &DB{DB: sqlxDB}

		// Configurar expectativa para erro
		mock.ExpectPing().WillReturnError(sqlmock.ErrCancelled)

		err = db.Health()
		assert.Error(t, err)

		// Verificar expectativas
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestDatabaseQueries(t *testing.T) {
	// Criar mock de banco
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	sqlxDB := sqlx.NewDb(mockDB, "sqlite3")
	db := &DB{DB: sqlxDB}

	t.Run("SELECT query", func(t *testing.T) {
		// Configurar expectativa para query SELECT
		rows := sqlmock.NewRows([]string{"id", "name"}).
			AddRow(1, "Test User")

		mock.ExpectQuery("SELECT id, name FROM users WHERE id = ?").
			WithArgs(1).
			WillReturnRows(rows)

		// Executar query
		var id int
		var name string
		err := db.Get(&struct {
			ID   int    `db:"id"`
			Name string `db:"name"`
		}{ID: id, Name: name}, "SELECT id, name FROM users WHERE id = ?", 1)

		assert.NoError(t, err)
	})

	t.Run("INSERT query", func(t *testing.T) {
		// Configurar expectativa para INSERT
		mock.ExpectExec("INSERT INTO users \\(name, email\\) VALUES \\(\\?, \\?\\)").
			WithArgs("John Doe", "john@example.com").
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Executar INSERT
		result, err := db.Exec("INSERT INTO users (name, email) VALUES (?, ?)", "John Doe", "john@example.com")
		assert.NoError(t, err)

		rowsAffected, err := result.RowsAffected()
		assert.NoError(t, err)
		assert.Equal(t, int64(1), rowsAffected)
	})

	t.Run("UPDATE query", func(t *testing.T) {
		// Configurar expectativa para UPDATE
		mock.ExpectExec("UPDATE users SET name = \\? WHERE id = \\?").
			WithArgs("Jane Doe", 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		// Executar UPDATE
		result, err := db.Exec("UPDATE users SET name = ? WHERE id = ?", "Jane Doe", 1)
		assert.NoError(t, err)

		rowsAffected, err := result.RowsAffected()
		assert.NoError(t, err)
		assert.Equal(t, int64(1), rowsAffected)
	})

	t.Run("DELETE query", func(t *testing.T) {
		// Configurar expectativa para DELETE
		mock.ExpectExec("DELETE FROM users WHERE id = \\?").
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		// Executar DELETE
		result, err := db.Exec("DELETE FROM users WHERE id = ?", 1)
		assert.NoError(t, err)

		rowsAffected, err := result.RowsAffected()
		assert.NoError(t, err)
		assert.Equal(t, int64(1), rowsAffected)
	})

	// Verificar expectativas
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestDatabaseTransactions(t *testing.T) {
	// Criar mock de banco
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	sqlxDB := sqlx.NewDb(mockDB, "sqlite3")
	db := &DB{DB: sqlxDB}

	t.Run("Successful transaction", func(t *testing.T) {
		// Configurar expectativas para transação
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO users").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		// Executar transação
		tx, err := db.Begin()
		assert.NoError(t, err)

		_, err = tx.Exec("INSERT INTO users (name) VALUES (?)", "Test")
		assert.NoError(t, err)

		err = tx.Commit()
		assert.NoError(t, err)
	})

	t.Run("Failed transaction with rollback", func(t *testing.T) {
		// Configurar expectativas para transação com rollback
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO users").WillReturnError(sqlmock.ErrCancelled)
		mock.ExpectRollback()

		// Executar transação com erro
		tx, err := db.Begin()
		assert.NoError(t, err)

		_, err = tx.Exec("INSERT INTO users (name) VALUES (?)", "Test")
		assert.Error(t, err)

		err = tx.Rollback()
		assert.NoError(t, err)
	})

	// Verificar expectativas
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

// Benchmark para testar performance de queries
func BenchmarkDatabaseQueries(b *testing.B) {
	// Criar mock de banco
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		b.Fatal(err)
	}
	defer mockDB.Close()

	sqlxDB := sqlx.NewDb(mockDB, "sqlite3")
	db := &DB{DB: sqlxDB}

	// Configurar expectativas para múltiplas queries
	for i := 0; i < b.N; i++ {
		rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
		mock.ExpectQuery("SELECT id FROM users LIMIT 1").WillReturnRows(rows)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var id int
		db.Get(&id, "SELECT id FROM users LIMIT 1")
	}
}
