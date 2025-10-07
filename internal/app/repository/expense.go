package repository

import (
	"Backeven/internal/app/ds"
	"fmt"
)

func (r *Repository) GetAllExpense() ([]ds.Expense, error) {
	var expenses []ds.Expense
	err := r.db.Find(&expenses).Error
	if err != nil {
		return nil, err
	}
	if len(expenses) == 0 {
		return nil, fmt.Errorf("затрат не найдено")
	}
	return expenses, err
}
func (r *Repository) GetExpenseByID(id int) (ds.Expense, error) {
	expense := ds.Expense{}
	err := r.db.Where(`"ExpenseID" = ?`, id).First(&expense).Error
	if err != nil {
		return ds.Expense{}, err
	}
	return expense, err
}
func (r *Repository) GetExpenseByTitle(Title string) ([]ds.Expense, error) {
	var expenses []ds.Expense
	err := r.db.Where(`"Title" LIKE ?`, "%"+Title+"%").Find(&expenses).Error
	if err != nil {
		return nil, err
	}
	return expenses, nil
}
