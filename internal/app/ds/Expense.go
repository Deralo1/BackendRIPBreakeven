package ds

type Expense struct { // услуга
	ExpenseID        int    `gorm:"primaryKey;column:ExpenseID"`
	Title            string `gorm:"type varchar(100);column:ExpenseTitle;not null"`
	ImageURL         string `gorm:"type varchar(255);column:ImageUrl"`
	Price            int    `gorm:"column:Price;not null"`
	ShortDescription string `gorm:"type text;column:ShortDesc;not null"`
	Description      string `gorm:"type text;column:Description;not null"`
	IsDeleted        bool   `gorm:"type boolean;column:IsDeleted;not null"`
}
