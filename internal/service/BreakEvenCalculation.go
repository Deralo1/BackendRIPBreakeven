package service

import (
	"Backeven/internal/app/ds"
	"fmt"
)

func CalculateAnswer(req *ds.BreakevenRequest) (int, error) {
	var fixedCosts int
	var variableCosts int

	//Считаем
	for _, exp := range req.RequestExpense {
		switch exp.TypeSpend {
		case 1: // постоянные
			fixedCosts += exp.Expense.Price * exp.AmountService
		case 2: // переменные
			variableCosts += exp.Expense.Price * exp.AmountService
		}
	}
	pricePerUnit := req.AmountProduct
	if pricePerUnit-variableCosts == 0 {
		return 0, fmt.Errorf("деление на ноль, так как цена за единицу продукции равно переменным затратам")
	}
	breakeven := fixedCosts / (pricePerUnit - variableCosts)
	return breakeven, nil
}
