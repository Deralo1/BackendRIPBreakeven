package ds

import (
	"time"
)

type BreakevenRequestDTO struct { //заявка
	BreakevenRequestID     int                    `json:"BreakevenRequestID"`
	BreakevenRequestStatus string                 `json:"BreakevenRequestStatus"`
	CreationDate           time.Time              `json:"CreationDate"`
	CreatorLogin           string                 `json:"CreatorLogin"`
	FormatedAt             *time.Time             `json:"FormatedAt,omitempty"`
	CompletedAt            *time.Time             `json:"CompletedAt,omitempty"`
	ModeratorLogin         *string                `json:"ModeratorLogin,omitempty"`
	AmountProduct          int                    `json:"AmountProduct"`
	CalcAnswer             int                    `json:"CalcAnswer"`
	RequestExpense         []ExpenseForRequestDTO `json:"RequestExpense"`
}

type ExpenseDTO struct { // услуга
	ExpenseID        int    `json:"ExpenseID"`
	Title            string `json:"Title"`
	ImageURL         string `json:"ImageURL,omitempty"`
	Price            int    `json:"Price"`
	ShortDescription string `json:"ShortDescription"`
	Description      string `json:"Description"`
}

type UserDTO struct {
	UserId int    `json:"UserId"`
	Login  string `json:"Login"`
}
type ExpenseForRequestDTO struct {
	ExpenseID     int    `json:"ExpenseID"`
	Title         string `json:"Title"`
	ImageURL      string `json:"ImageURL"`
	AmountService int    `json:"AmountService"`
	TypeSpend     int    `json:"TypeSpend"`
}
type UpdateRequestExpenseDTO struct {
	AmountService int `json:"AmountService"`
	TypeSpend     int `json:"TypeSpend"`
}
type UpdateexpenseDTO struct { // post без изображения
	Title            string `json:"Title"`
	Price            int    `json:"Price"`
	ShortDescription string `json:"ShortDescription"`
	Description      string `json:"Description"`
}
type UpdateBreakEvenCalcDTO struct {
	AmountProduct int `json:"AmountProduct"`
	CalcAnswer    int `json:"CalcAnswer"`
}
type ChangeUserDTO struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
