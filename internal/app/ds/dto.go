package ds

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
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
	UserId int      `json:"UserId"`
	Login  string   `json:"Login"`
	Role   UserRole `json:"Role"`
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

type UserRole string

const (
	RoleGuest     UserRole = "guest"
	RoleCreator   UserRole = "creator" // обычный пользователь
	RoleModerator UserRole = "moderator"
)

// JWT claims определяет данные которые мы храним в токене
type JWTClaims struct {
	jwt.RegisteredClaims
	UserID int      `json:"user_id"`
	Role   UserRole `json:"role"`
}
type AuthResponseDTO struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"` // Unix timestamp истечения
}
type BreakevenResult struct {
	ID           int    `json:"id"`
	Breakeven    int    `json:"breakeven"`
	BackendToken string `json:"backend_token"`
}
