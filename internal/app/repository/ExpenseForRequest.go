package repository

import "Backeven/internal/app/ds"

func (r *Repository) UpdateExpenseForReq(userid int, expenseid int, AmountService int, TypeSpend int) error {
	Calc, err := r.GetBreakevenCalc(userid)

	if err != nil {
		return err
	}

	var er ds.ExpenseForRequest

	err = r.db.Where(`"BreakevenRequestID" = ? AND "ExpenseID" = ?`, Calc.BreakevenRequestID, expenseid).First(&er).Error

	if err != nil {
		return err
	}

	er.AmountService = AmountService
	er.TypeSpend = TypeSpend

	return r.db.Save(&er).Error
}

func (r *Repository) DeleteFromCalc(expenseid int, userid int) error {
	Calc, err := r.GetBreakevenCalc(userid)
	if err != nil {
		return err
	}

	return r.db.Where(`"BreakevenRequestID" = ? AND "ExpenseID" = ?`, Calc.BreakevenRequestID, expenseid).Delete(&ds.ExpenseForRequest{}).Error
}
