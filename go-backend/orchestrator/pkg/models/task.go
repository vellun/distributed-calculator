package models

type Task struct {
	Id         int    `json:"id"`
	Exp_id     int    `json:"exp_id"`
	Operand1   string `json:"operand1"`  // Член подвыражения число
	Operand2   string `json:"operand2"`
	Operation  string `json:"operation"`
	Task_id1   int    `json:"task_id1"`  // Член выражения ссылка на другое подвыражение
	Task_id2   int    `json:"task_id2"`
	Status     string `json:"status"` // process/complete
	Seq_number int    `json:"seq_number"`
	Duration   int    `json:"duration"`
	Result     int    `json:"result"`
}
