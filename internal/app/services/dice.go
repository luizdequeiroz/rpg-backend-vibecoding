package services

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/luizdequeiroz/rpg-backend/internal/app/models"
	"github.com/luizdequeiroz/rpg-backend/internal/app/repositories"
)

type DiceService struct {
	rollRepo *repositories.RollRepository
}

func NewDiceService(rollRepo *repositories.RollRepository) *DiceService {
	return &DiceService{
		rollRepo: rollRepo,
	}
}

// DiceExpression representa uma expressão de dados parseada
type DiceExpression struct {
	Count    int
	Sides    int
	Modifier int
}

// ParseDiceExpression analisa uma expressão como "1d20+3" ou "2d6-1"
func (s *DiceService) ParseDiceExpression(expression string) (*DiceExpression, error) {
	// Regex para capturar expressões como XdY+Z ou XdY-Z
	re := regexp.MustCompile(`^(\d+)d(\d+)([\+\-]\d+)?$`)
	matches := re.FindStringSubmatch(strings.ToLower(strings.TrimSpace(expression)))

	if len(matches) < 3 {
		return nil, fmt.Errorf("expressão inválida: %s. Use formato como '1d20+3'", expression)
	}

	count, err := strconv.Atoi(matches[1])
	if err != nil || count < 1 || count > 100 {
		return nil, fmt.Errorf("número de dados inválido: deve ser entre 1 e 100")
	}

	sides, err := strconv.Atoi(matches[2])
	if err != nil || sides < 2 || sides > 1000 {
		return nil, fmt.Errorf("número de lados inválido: deve ser entre 2 e 1000")
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

// RollDice executa uma rolagem de dados
func (s *DiceService) RollDice(expression string, userID int) (*models.DiceRollResponse, error) {
	dice, err := s.ParseDiceExpression(expression)
	if err != nil {
		return nil, err
	}

	// Executar rolagem
	rolls := make([]int, dice.Count)
	total := 0

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < dice.Count; i++ {
		roll := rand.Intn(dice.Sides) + 1
		rolls[i] = roll
		total += roll
	}

	finalResult := total + dice.Modifier

	// Verificar crítico/fumble (apenas para d20)
	isCritical := false
	isFumble := false
	if dice.Sides == 20 && dice.Count == 1 {
		if rolls[0] == 20 {
			isCritical = true
		} else if rolls[0] == 1 {
			isFumble = true
		}
	}

	// Criar detalhes da rolagem
	details := s.formatRollDetails(rolls, dice.Modifier, finalResult)

	// Salvar no banco
	roll := &models.Roll{
		Expression:    expression,
		ResultValue:   finalResult,
		ResultDetails: details,
		UserID:        userID,
		CreatedAt:     time.Now(),
	}

	err = s.rollRepo.Create(roll)
	if err != nil {
		return nil, fmt.Errorf("erro ao salvar rolagem: %v", err)
	}

	return &models.DiceRollResponse{
		ID:         roll.ID,
		Expression: roll.Expression,
		Result:     roll.ResultValue,
		Details:    roll.ResultDetails,
		IsCritical: isCritical,
		IsFumble:   isFumble,
		SheetID:    roll.SheetID,
		UserID:     roll.UserID,
		CreatedAt:  roll.CreatedAt,
	}, nil
}

// getNestedValue busca um valor em uma estrutura aninhada usando dot notation
func getNestedValue(data map[string]interface{}, path string) (interface{}, bool) {
	parts := strings.Split(path, ".")
	current := data

	for i, part := range parts {
		if i == len(parts)-1 {
			// Último nível, retornar o valor
			value, exists := current[part]
			return value, exists
		}

		// Nível intermediário, deve ser um mapa
		next, exists := current[part]
		if !exists {
			return nil, false
		}

		nextMap, ok := next.(map[string]interface{})
		if !ok {
			return nil, false
		}

		current = nextMap
	}

	return nil, false
}

// RollWithSheet executa rolagem com dados da ficha
func (s *DiceService) RollWithSheet(expression, attributeField string, sheet *models.PlayerSheetResponse) (*models.DiceRollResponse, error) {
	// Substituir placeholders na expressão
	finalExpression := expression

	if attributeField != "" && strings.Contains(expression, "{"+attributeField+"}") {
		value, exists := getNestedValue(sheet.Data, attributeField)
		if !exists {
			return nil, fmt.Errorf("campo '%s' não encontrado na ficha", attributeField)
		}

		// Converter valor para int
		var modifier int
		switch v := value.(type) {
		case float64:
			modifier = int(v)
		case int:
			modifier = v
		case string:
			var err error
			modifier, err = strconv.Atoi(v)
			if err != nil {
				return nil, fmt.Errorf("valor do atributo '%s' não é numérico", attributeField)
			}
		default:
			return nil, fmt.Errorf("tipo de valor inválido para atributo '%s'", attributeField)
		}

		finalExpression = strings.ReplaceAll(expression, "{"+attributeField+"}", strconv.Itoa(modifier))
	}

	// Executar rolagem normal
	result, err := s.RollDice(finalExpression, sheet.OwnerID)
	if err != nil {
		return nil, err
	}

	// Atualizar com ID da ficha
	if sheet.ID != "" {
		result.SheetID = &sheet.ID
		// Atualizar no banco também
		s.rollRepo.UpdateSheetID(result.ID, sheet.ID)
	}

	return result, nil
}

// GetUserHistory recupera histórico de rolagens do usuário
func (s *DiceService) GetUserHistory(userID, page, limit int) ([]models.DiceRollResponse, int, error) {
	offset := (page - 1) * limit

	rolls, err := s.rollRepo.GetByUserID(userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.rollRepo.CountByUserID(userID)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]models.DiceRollResponse, len(rolls))
	for i, roll := range rolls {
		responses[i] = models.DiceRollResponse{
			ID:         roll.ID,
			Expression: roll.Expression,
			Result:     roll.ResultValue,
			Details:    roll.ResultDetails,
			IsCritical: false, // Calcular baseado nos detalhes se necessário
			IsFumble:   false, // Calcular baseado nos detalhes se necessário
			SheetID:    roll.SheetID,
			UserID:     roll.UserID,
			CreatedAt:  roll.CreatedAt,
		}
	}

	return responses, total, nil
}

// formatRollDetails formata os detalhes da rolagem
func (s *DiceService) formatRollDetails(rolls []int, modifier, total int) string {
	if len(rolls) == 1 {
		if modifier == 0 {
			return fmt.Sprintf("[%d] = %d", rolls[0], total)
		}
		return fmt.Sprintf("[%d] %+d = %d", rolls[0], modifier, total)
	}

	rollsStr := make([]string, len(rolls))
	for i, roll := range rolls {
		rollsStr[i] = strconv.Itoa(roll)
	}

	if modifier == 0 {
		return fmt.Sprintf("[%s] = %d", strings.Join(rollsStr, ", "), total)
	}
	return fmt.Sprintf("[%s] %+d = %d", strings.Join(rollsStr, ", "), modifier, total)
}
