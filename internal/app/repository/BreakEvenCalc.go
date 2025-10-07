package repository

import (
	"Backeven/internal/app/ds"
	"fmt"
	"time"
)

func (r *Repository) GetBreakevenCalc(userid int) (*ds.BreakevenRequest, error) {
	var calc ds.BreakevenRequest
	err := r.db.Where(`"Creator_ID" = ? AND "BreakEvenStatus" = 'черновик'`, userid).
		Preload("RequestExpense").
		Preload("RequestExpense.Expense").
		First(&calc).Error
	if err != nil {
		if err.Error() == "record not found" {
			newBreakEvenRequest := &ds.BreakevenRequest{
				CreatorID:    userid,
				CreationDate: time.Now(),
			}
			err = r.db.Create(newBreakEvenRequest).Error
			if err != nil {
				return nil, err
			}
			// Загружаем созданную заявку
			err = r.db.Where(`"Creator_ID" = ? AND "BreakEvenStatus" = 'черновик'`, userid).
				Preload("RequestExpense").
				Preload("RequestExpense.Expense").
				First(&calc).Error
			if err != nil {
				return nil, err
			}
			return &calc, nil
		}
		return nil, err
	}
	return &calc, err
}
func (r *Repository) GetBreakevenCalcByID(CalcID uint) (*ds.BreakevenRequest, error) {
	var calc ds.BreakevenRequest
	err := r.db.Where(`"BreakevenRequestID" = ? AND "BreakEvenStatus" != 'удалён'`, CalcID).
		Preload("RequestExpense").
		Preload("RequestExpense.Expense").
		First(&calc).Error
	if err != nil {
		return nil, err
	}
	return &calc, err
}
func (r *Repository) GetExpensesInCalcCount(UserId int) int {
	calc, err := r.GetBreakevenCalc(UserId)
	if err != nil || calc == nil || calc.BreakevenRequestID == 0 {
		return 0
	}
	var count int64
	err = r.db.Model(&ds.ExpenseForRequest{}).
		Where(`"BreakevenRequestID" = ?`, calc.BreakevenRequestID).Count(&count).Error
	if err != nil {
		return 0
	}
	return int(count)
}
func (r *Repository) AddExpenseToCalc(UserID, expenseID int) error {
	// Получаем текущую заявку или создаем новую, если ее нет
	calc, err := r.GetBreakevenCalc(UserID)
	if err != nil {
		return err
	}
	//Проверяем нет ли в корзине
	var count int64
	err = r.db.Model(&ds.ExpenseForRequest{}).
		Where(`"BreakevenRequestID" = ? AND "ExpenseID = ?"`, calc.BreakevenRequestID, expenseID).Count(&count).Error
	if err != nil {
		return err
	}
	// если уже добавлен возвращаем ошибку
	if count > 0 {
		return fmt.Errorf("трата уже добавлена в корзину")
	}
	// Добавляем в корзину
	item := ds.ExpenseForRequest{
		BreakevenRequestID: calc.BreakevenRequestID,
		ExpenseID:          expenseID,
		AmountService:      0,
		TypeSpend:          1,
	}
	return r.db.Create(&item).Error
}
