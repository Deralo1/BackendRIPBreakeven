package ds

type Expense struct { // услуга
	ExpenseID        int    `gorm:"primaryKey;column:ExpenseID"`
	Title            string `gorm:"type varchar(100);column:ExpenseTitle"`
	ImageURL         string `gorm:"type varchar(255);column:ImageUrl"`
	Price            int    `gorm:"column:Price"`
	ShortDescription string `gorm:"type text;column:ShortDesc"`
	Description      string `gorm:"type text;column:Description"`
	IsDeleted        bool   `gorm:"type boolean;column:IsDeleted"`
}
