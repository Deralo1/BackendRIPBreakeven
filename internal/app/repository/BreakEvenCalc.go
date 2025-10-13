package repository

import (
	"Backeven/internal/app/ds"
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
			// создаем новую заявку
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
func (r *Repository) DeleteBreakEvenCalc(CalcID uint) error {
	return r.db.Exec(`Update Breakeven_Requests SET "BreakEvenStatus"= 'удалён' WHERE "BreakevenRequestID" = ?`, CalcID).Error
}
