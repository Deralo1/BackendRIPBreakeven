package repository

import (
	"Backeven/internal/app/ds"
	"fmt"
)

func (r *Repository) GetAllExpense() ([]ds.ExpenseDTO, error) {
	var expenses []ds.Expense
	err := r.db.Where(`"IsDeleted" = false`).Find(&expenses).Error
	if err != nil {
		return nil, err
	}
	if len(expenses) == 0 {
		return nil, fmt.Errorf("затрат не найдено")
	}
	expenseDTOs := make([]ds.ExpenseDTO, len(expenses))
	for i, e := range expenses {
		expenseDTOs[i] = ds.ExpenseDTO{
			ExpenseID:        e.ExpenseID,
			Title:            e.Title,
			ImageURL:         e.ImageURL,
			Price:            e.Price,
			ShortDescription: e.ShortDescription,
			Description:      e.Description,
		}
	}
	return expenseDTOs, err
}

func (r *Repository) GetExpenseByTitle(Title string) ([]ds.ExpenseDTO, error) {
	var expenses []ds.Expense
	err := r.db.Where(`"ExpenseTitle" ILIKE ? AND "IsDeleted" = false`, "%"+Title+"%").Find(&expenses).Error
	if err != nil {
		return nil, err
	}
	expenseDTOs := make([]ds.ExpenseDTO, len(expenses))
	for i, e := range expenses {
		expenseDTOs[i] = ds.ExpenseDTO{
			ExpenseID:        e.ExpenseID,
			Title:            e.Title,
			ImageURL:         e.ImageURL,
			Price:            e.Price,
			ShortDescription: e.ShortDescription,
			Description:      e.Description,
		}
	}
	return expenseDTOs, nil
}
func (r *Repository) GetExpenseByID(id int) (ds.ExpenseDTO, error) {
	expense := ds.Expense{}
	err := r.db.Where(`"ExpenseID" = ? and "IsDeleted" = false`, id).First(&expense).Error
	if err != nil {
		return ds.ExpenseDTO{}, err
	}
	return ds.ExpenseDTO{
		ExpenseID:        expense.ExpenseID,
		Title:            expense.Title,
		ImageURL:         expense.ImageURL,
		Price:            expense.Price,
		ShortDescription: expense.ShortDescription,
		Description:      expense.Description,
	}, nil
}
func (r *Repository) CreateExpense(expenseDTO ds.ExpenseDTO) (*ds.ExpenseDTO, error) {
	expense := ds.Expense{
		Title:            expenseDTO.Title,
		ImageURL:         expenseDTO.ImageURL,
		Price:            expenseDTO.Price,
		ShortDescription: expenseDTO.ShortDescription,
		Description:      expenseDTO.Description,
		IsDeleted:        false,
	}
	err := r.db.Create(&expense).Error
	if err != nil {
		return nil, err
	}
	return &ds.ExpenseDTO{
		ExpenseID:        expense.ExpenseID,
		Title:            expense.Title,
		ImageURL:         expense.ImageURL,
		Price:            expense.Price,
		ShortDescription: expense.ShortDescription,
		Description:      expense.Description,
	}, nil
}
func (r *Repository) Updateexpense(id int, expenseUpdate ds.UpdateexpenseDTO) (*ds.ExpenseDTO, error) {
	var expense ds.Expense
	err := r.db.Where(`"ExpenseID" = ? and "IsDeleted" =false`, id).First(&expense).Error
	if err != nil {
		return nil, err
	}
	if expenseUpdate.Title != "" {
		expense.Title = expenseUpdate.Title
	}
	if expenseUpdate.Price != 0 {
		expense.Price = expenseUpdate.Price
	}
	if expenseUpdate.Description != "" {
		expense.Description = expenseUpdate.Description
	}
	if expenseUpdate.ShortDescription != "" {
		expense.ShortDescription = expenseUpdate.ShortDescription
	}
	err = r.db.Save(&expense).Error
	if err != nil {
		return nil, err
	}
	return &ds.ExpenseDTO{
		ExpenseID:        expense.ExpenseID,
		Title:            expense.Title,
		ImageURL:         expense.ImageURL,
		Price:            expense.Price,
		ShortDescription: expense.ShortDescription,
		Description:      expense.Description,
	}, nil
}

func (r *Repository) DeleteExpense(id int) error {
	return r.db.Model(&ds.Expense{}).Where(`"ExpenseID" = ? AND "IsDeleted" =false`, id).Update("IsDeleted", true).Error
}

func (r *Repository) UpdateExpenseImage(id int, imageURL string) (*ds.ExpenseDTO, error) {
	var expense ds.Expense
	err := r.db.Where(`"ExpenseID" = ? and "IsDeleted" = false`, id).First(&expense).Error
	if err != nil {
		return nil, err
	}
	expense.ImageURL = imageURL
	err = r.db.Save(&expense).Error
	if err != nil {
		return nil, err
	}
	return &ds.ExpenseDTO{
		ExpenseID:        expense.ExpenseID,
		Title:            expense.Title,
		ImageURL:         expense.ImageURL,
		Price:            expense.Price,
		ShortDescription: expense.ShortDescription,
		Description:      expense.Description,
	}, nil
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
		Where(`"BreakevenRequestID" = ? AND "ExpenseID" = ?`, calc.BreakevenRequestID, expenseID).Count(&count).Error
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
