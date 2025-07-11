package roll

import (
	"testing"

	"github.com/luizdequeiroz/rpg-backend/internal/app/models"
	"github.com/stretchr/testify/assert"
)

func TestParseExpression(t *testing.T) {
	engine := NewRollEngine()

	tests := []struct {
		name     string
		expr     string
		expected DiceExpression
		hasError bool
	}{
		{
			name: "1d20",
			expr: "1d20",
			expected: DiceExpression{
				Count:    1,
				Sides:    20,
				Modifier: 0,
			},
			hasError: false,
		},
		{
			name: "2d6+3",
			expr: "2d6+3",
			expected: DiceExpression{
				Count:    2,
				Sides:    6,
				Modifier: 3,
			},
			hasError: false,
		},
		{
			name: "3d8-1",
			expr: "3d8-1",
			expected: DiceExpression{
				Count:    3,
				Sides:    8,
				Modifier: -1,
			},
			hasError: false,
		},
		{
			name:     "Expressão inválida",
			expr:     "abc",
			expected: DiceExpression{},
			hasError: true,
		},
		{
			name:     "Dados zero",
			expr:     "0d20",
			expected: DiceExpression{},
			hasError: true,
		},
		{
			name:     "Faces inválidas",
			expr:     "1d1",
			expected: DiceExpression{},
			hasError: true,
		},
		{
			name:     "Muitos dados",
			expr:     "101d20",
			expected: DiceExpression{},
			hasError: true,
		},
		{
			name:     "Muitas faces",
			expr:     "1d1001",
			expected: DiceExpression{},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := engine.ParseExpression(tt.expr)

			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.Count, result.Count)
				assert.Equal(t, tt.expected.Sides, result.Sides)
				assert.Equal(t, tt.expected.Modifier, result.Modifier)
			}
		})
	}
}

func TestRoll(t *testing.T) {
	engine := NewRollEngine()

	tests := []struct {
		name      string
		expr      string
		minResult int
		maxResult int
	}{
		{
			name:      "1d20",
			expr:      "1d20",
			minResult: 1,
			maxResult: 20,
		},
		{
			name:      "2d6+3",
			expr:      "2d6+3",
			minResult: 5,  // 2*1 + 3
			maxResult: 15, // 2*6 + 3
		},
		{
			name:      "3d8-1",
			expr:      "3d8-1",
			minResult: 2,  // 3*1 - 1
			maxResult: 23, // 3*8 - 1
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := engine.Roll(tt.expr)
			assert.NoError(t, err)
			assert.GreaterOrEqual(t, result.Total, tt.minResult)
			assert.LessOrEqual(t, result.Total, tt.maxResult)
			assert.NotEmpty(t, result.Dice)
		})
	}
}

func TestRollCriticalAndFumble(t *testing.T) {
	engine := NewRollEngine()

	// Para testar critical/fumble de forma consistente, vamos simular casos específicos
	tests := []struct {
		name         string
		expr         string
		expectCrit   bool
		expectFumble bool
	}{
		{
			name:         "d20 pode ter critical ou fumble",
			expr:         "1d20",
			expectCrit:   false, // Não podemos garantir o resultado
			expectFumble: false, // Só verificamos que não gera erro
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := engine.Roll(tt.expr)
			assert.NoError(t, err)
			assert.NotNil(t, result)

			// Verificar estrutura básica
			assert.GreaterOrEqual(t, result.Total, 1)
			assert.LessOrEqual(t, result.Total, 20)

			// Critical e fumble são mutuamente exclusivos
			if result.Critical {
				assert.False(t, result.Fumble)
			}
			if result.Fumble {
				assert.False(t, result.Critical)
			}
		})
	}
}

func TestRollFromField(t *testing.T) {
	engine := NewRollEngine()

	// Dados de ficha de teste
	sheetData := models.PlayerSheetData{
		"attributes": map[string]interface{}{
			"strength":  15,
			"dexterity": 12,
		},
		"skills": map[string]interface{}{
			"arcana":  8,
			"stealth": "1d20+2",
		},
	}

	tests := []struct {
		name      string
		fieldName string
		hasError  bool
	}{
		{
			name:      "Campo numérico existente",
			fieldName: "attributes.strength",
			hasError:  false,
		},
		{
			name:      "Campo de expressão existente",
			fieldName: "skills.stealth",
			hasError:  false,
		},
		{
			name:      "Campo inexistente",
			fieldName: "attributes.nonexistent",
			hasError:  true,
		},
		{
			name:      "Caminho inválido",
			fieldName: "invalid.path.deep",
			hasError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, expression, err := engine.RollFromField(sheetData, tt.fieldName)

			if tt.hasError {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Empty(t, expression)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.NotEmpty(t, expression)
				assert.GreaterOrEqual(t, result.Total, 1)
			}
		})
	}
}

func TestEvaluateSuccess(t *testing.T) {
	engine := NewRollEngine()

	result := &models.RollDetails{
		Dice:     []int{15},
		Modifier: 3,
		Total:    18,
		Critical: false,
		Fumble:   false,
	}

	tests := []struct {
		name       string
		difficulty int
		expected   bool
	}{
		{
			name:       "Sucesso fácil",
			difficulty: 10,
			expected:   true,
		},
		{
			name:       "Sucesso exato",
			difficulty: 18,
			expected:   true,
		},
		{
			name:       "Falha",
			difficulty: 20,
			expected:   false,
		},
		{
			name:       "Sem dificuldade",
			difficulty: 0,
			expected:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			success := engine.EvaluateSuccess(result, tt.difficulty)
			assert.Equal(t, tt.expected, success)
		})
	}
}

// Benchmarks para testar performance
func BenchmarkParseExpression(b *testing.B) {
	engine := NewRollEngine()
	for i := 0; i < b.N; i++ {
		engine.ParseExpression("2d6+3")
	}
}

func BenchmarkRoll(b *testing.B) {
	engine := NewRollEngine()
	for i := 0; i < b.N; i++ {
		engine.Roll("1d20+5")
	}
}

func BenchmarkRollFromField(b *testing.B) {
	engine := NewRollEngine()
	sheetData := models.PlayerSheetData{
		"attributes": map[string]interface{}{
			"strength": 15,
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		engine.RollFromField(sheetData, "attributes.strength")
	}
}
