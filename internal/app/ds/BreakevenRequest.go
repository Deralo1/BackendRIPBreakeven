package ds

import (
	"database/sql"
	"time"
)

type BreakevenRequest struct { //заявка
	BreakevenRequestID     int           `gorm:"PrimaryKey;column:BreakevenRequestID"`
	BreakevenRequestStatus string        `gorm:"type varchar(100);not null;default:'черновик';column:BreakEvenStatus"`
	CreationDate           time.Time     `gorm:"not null;column:CreatedAt"`
	CreatorID              int           `gorm:"not null;column:Creator_ID"`
	FormatedAt             sql.NullTime  `gorm:"column:FormatedAt"`
	CompletedAt            sql.NullTime  `gorm:"column:CompletedAt"`
	ModeratorID            sql.NullInt64 `gorm:"column:Moderator_Id"`
	AmountProduct          int           `gorm:"column:AmountProduct"`
	CalcAnswer             int           `gorm:"column:CalcAnswer"`

	RequestExpense []ExpenseForRequest `gorm:"foreignKey:BreakevenRequestID"`
	Moderator      User                `gorm:"foreignKey:ModeratorID;references:UserID"`
	Creator        User                `gorm:"foreignKey:CreatorID;references:UserID"`
}
