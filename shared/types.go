package shared

// Применяется при запросе к оркестратору со строкой выражения
// /api/v1/calculate
type CalculateRequest struct {
	Expression string `json:"expression"`
}

type CalculateResponse struct {
	ID int64 `json:"id"`
}

type TaskResult struct {
	ID     int64   `json:"id"`
	Result float64 `json:"result"`
}

// Применяется при запросе к оркестратору
// /api/v1/expression/[:id]
type ExpressionRequest struct {
	Expression string `json:"expression"`
}

// Универсальный тип выражения
type Expression struct {
	ID     int64  `json:"id"`
	Status bool   `json:"status"`
	Result string `json:"result"`
}

// Список из выражений, выдающийся оркестратором при запросе
// /api/v1/expressions
type ExpressionList struct {
	Expressions []Expression `json:"expressions"`
}

type Task struct {
	ID             int64   `json:"id"`
	FirstArgument  string  `json:"arg1"`
	SecondArgument string  `json:"arg2"`
	Operator       rune    `json:"operator"`
	OperationTime  float64 `json:"operation_time"`
	Status         bool    `json:"status"`
	Result         float64 `json:"result"`
}
