package request

import "encoding/json"

//GetRequest Структура тела запроса на получение данных
type GetRequest struct {
	HookURL      string                       `json:"hook_url"`
	Source       string                       `json:"source"`
	Requirements *RequirementExpression       `json:"requirements"`
	Modifiers    []map[string]json.RawMessage `json:"modifiers"`
	Fields       []Field                      `json:"fields"` // Data contains slice of fields names
}

type Modifier struct {
	// Doesn't implement yet
}

// RequirementExpression структура условия
type RequirementExpression struct {
	Type  string                   `json:"type"`
	OrAnd []*RequirementExpression `json:"or_and"`
	Req   Requirement              `json:"requirement"`
}

type Requirement struct {
	Operator string             `json:"operator"`
	Field    Field              `json:"field"`
	Operand  RequirementOperand `json:"operand"`
}

type Field struct {
	Child *Field `json:"child"`
	Name  string `json:"name"`
}

type RequirementOperand struct {
	Type string      `json:"type"`
	Spec interface{} `json:"spec"`
}
