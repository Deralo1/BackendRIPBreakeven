package ds

type ExpenseForRequest struct {
	BreakevenRequestID int `gorm:"primaryKey;not null;column:BreakevenRequestID;uniqueIndex:idx_request_expense"`
	ExpenseID          int `gorm:"primaryKey;not null;column:ExpenseID;uniqueIndex:idx_request_expense"`
	AmountService      int `gorm:"column:AmountService"`
	TypeSpend          int `gorm:"column:TypeSpend"`

	BreakevenRequest BreakevenRequest `gorm:"foreignKey:BreakevenRequestID"`
	Expense          Expense          `gorm:"foreignKey:ExpenseID"`
}
