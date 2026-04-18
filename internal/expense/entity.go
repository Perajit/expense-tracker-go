package expense

func GetModels() []any {
	return []any{&ExpenseEntity{}, &CategoryEntity{}, &TagEntity{}}
}
