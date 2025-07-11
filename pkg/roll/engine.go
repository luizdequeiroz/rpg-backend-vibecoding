package roll

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/luizdequeiroz/rpg-backend/internal/app/models"
)

// RollEngine gerencia rolagens de dados
type RollEngine struct {
	rand *rand.Rand
}

// NewRollEngine cria nova instância do motor de rolagem
func NewRollEngine() *RollEngine {
	return &RollEngine{
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// getNestedValue busca um valor em uma estrutura aninhada usando dot notation
func (re *RollEngine) getNestedValue(data models.PlayerSheetData, path string) (interface{}, error) {
	// Converter para map[string]interface{} para navegação dinâmica
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar dados da ficha: %w", err)
	}

	var dataMap map[string]interface{}
	err = json.Unmarshal(dataBytes, &dataMap)
	if err != nil {
		return nil, fmt.Errorf("erro ao deserializar dados da ficha: %w", err)
	}

	parts := strings.Split(path, ".")
	current := dataMap

	for i, part := range parts {
		if i == len(parts)-1 {
			// Último nível, retornar o valor
			value, exists := current[part]
			if !exists {
				return nil, fmt.Errorf("campo '%s' não encontrado na ficha", path)
			}
			return value, nil
		}

		// Nível intermediário, deve ser um mapa
		next, exists := current[part]
		if !exists {
			return nil, fmt.Errorf("campo '%s' não encontrado na ficha", path)
		}

		nextMap, ok := next.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("campo '%s' não é um objeto, não é possível navegar mais", part)
		}

		current = nextMap
	}

	return nil, fmt.Errorf("campo '%s' não encontrado", path)
}

// DiceExpression representa uma expressão de dados
type DiceExpression struct {
	Count    int // Número de dados
	Sides    int // Número de lados
	Modifier int // Modificador (+/-)
}

// ParseExpression analisa expressão de dados (e.g., "1d20+3", "2d6-1")
func (re *RollEngine) ParseExpression(expression string) (*DiceExpression, error) {
	// Regex para capturar XdY+Z ou XdY-Z
	re_dice := regexp.MustCompile(`^(\d+)d(\d+)([+-]\d+)?$`)

	// Limpar espaços
	expression = strings.ReplaceAll(expression, " ", "")
	expression = strings.ToLower(expression)

	matches := re_dice.FindStringSubmatch(expression)
	if len(matches) < 3 {
		return nil, fmt.Errorf("expressão inválida: %s (use formato XdY+Z)", expression)
	}

	count, err := strconv.Atoi(matches[1])
	if err != nil || count < 1 || count > 100 {
		return nil, fmt.Errorf("número de dados inválido: %s (1-100)", matches[1])
	}

	sides, err := strconv.Atoi(matches[2])
	if err != nil || sides < 2 || sides > 1000 {
		return nil, fmt.Errorf("número de lados inválido: %s (2-1000)", matches[2])
	}

	modifier := 0
	if len(matches) > 3 && matches[3] != "" {
		modifier, err = strconv.Atoi(matches[3])
		if err != nil {
			return nil, fmt.Errorf("modificador inválido: %s", matches[3])
		}
	}

	return &DiceExpression{
		Count:    count,
		Sides:    sides,
		Modifier: modifier,
	}, nil
}

// Roll executa rolagem baseada na expressão
func (re *RollEngine) Roll(expression string) (*models.RollDetails, error) {
	dice_expr, err := re.ParseExpression(expression)
	if err != nil {
		return nil, err
	}

	// Rolar dados
	dice_results := make([]int, dice_expr.Count)
	total := 0

	for i := 0; i < dice_expr.Count; i++ {
		roll := re.rand.Intn(dice_expr.Sides) + 1
		dice_results[i] = roll
		total += roll
	}

	// Aplicar modificador
	final_total := total + dice_expr.Modifier

	// Verificar críticos e fumbles (apenas para d20)
	critical := false
	fumble := false

	if dice_expr.Sides == 20 && dice_expr.Count == 1 {
		if dice_results[0] == 20 {
			critical = true
		} else if dice_results[0] == 1 {
			fumble = true
		}
	}

	return &models.RollDetails{
		Dice:     dice_results,
		Modifier: dice_expr.Modifier,
		Total:    final_total,
		Critical: critical,
		Fumble:   fumble,
	}, nil
}

// RollFromField extrai valor de campo da ficha e executa rolagem
func (re *RollEngine) RollFromField(sheetData models.PlayerSheetData, fieldName string) (*models.RollDetails, string, error) {
	// Buscar campo na ficha (com suporte a campos aninhados)
	value, err := re.getNestedValue(sheetData, fieldName)
	if err != nil {
		return nil, "", err
	}

	// Converter para string
	var expression string
	switch v := value.(type) {
	case string:
		expression = v
	case int:
		// Se for número, assumir como modificador para d20
		expression = fmt.Sprintf("1d20%+d", v)
	case float64:
		// JSON numbers são float64
		expression = fmt.Sprintf("1d20%+d", int(v))
	default:
		return nil, "", fmt.Errorf("campo '%s' não é um valor de rolagem válido", fieldName)
	}

	// Verificar se é expressão ou valor numérico
	if !strings.Contains(expression, "d") {
		// É só um número, assumir d20
		modifier, err := strconv.Atoi(expression)
		if err == nil {
			expression = fmt.Sprintf("1d20%+d", modifier)
		}
	}

	result, err := re.Roll(expression)
	return result, expression, err
}

// EvaluateSuccess avalia se rolagem foi bem-sucedida baseada em dificuldade
func (re *RollEngine) EvaluateSuccess(result *models.RollDetails, difficulty int) bool {
	if difficulty <= 0 {
		return true // Sem dificuldade definida
	}
	return result.Total >= difficulty
}
