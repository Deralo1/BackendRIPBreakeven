package repository

import (
	"Backeven/internal/app/ds"
	"Backeven/internal/service"
	"database/sql"
	"fmt"
	"time"
)

func (r *Repository) GetCalcInfo(userid int) (currentCalcID int, count int64, err error) {
	var BreakEvenCalc ds.BreakevenRequest
	err = r.db.Where(`"Creator_ID" = ? AND "BreakEvenStatus" = 'черновик'`, userid).First(&BreakEvenCalc).Error
	if err != nil {
		// если черновика нет — возвращаем 0 без ошибки
		if err.Error() == "record not found" {
			return 0, 0, nil
		}
		// при реальной ошибке — возвращаем её
		return 0, 0, err
	}

	// считаем количество записей в корзине (ExpenseForRequest), а не количество заявок
	err = r.db.Model(&ds.ExpenseForRequest{}).Where(`"BreakevenRequestID" = ?`, BreakEvenCalc.BreakevenRequestID).Count(&count).Error
	if err != nil {
		return 0, 0, err
	}
	return BreakEvenCalc.BreakevenRequestID, count, nil

}

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
			err = r.db.Create(newBreakEvenRequest).Error
			if err != nil {
				return nil, err
			}
			return &calc, nil
		}
		return nil, err
	}
	return &calc, err
}

func (r *Repository) GetListCalcByDateAndStatus(
	userID int,
	role ds.UserRole,
	status string,
	startDate, endDate time.Time,
) ([]ds.BreakevenRequestDTO, error) {

	var calc []ds.BreakevenRequest

	// Базовый запрос: исключаем удалённые и черновики
	query := r.db.Where(`"BreakEvenStatus" != 'удалён' AND "BreakEvenStatus" != 'черновик'`)

	// 🔹 Фильтрация по роли и userID
	if role == ds.RoleCreator && userID != 0 {
		query = query.Where(`"Creator_ID" = ?`, userID)
	}
	if role == ds.RoleModerator && userID != 0 {
		query = query.Where(`"Moderator_Id" = ?`, userID)
	}

	// 🔹 Фильтрация по статусу
	if status != "" {
		query = query.Where(`"BreakEvenStatus" = ?`, status)
	}

	// 🔹 Фильтрация по датам
	if !startDate.IsZero() {
		query = query.Where(`"FormatedAt" >= ?`, startDate)
	}
	if !endDate.IsZero() {
		query = query.Where(`"CompletedAt" <= ?`, endDate)
	}

	// Выполняем запрос с подгрузкой связей
	err := query.
		Preload("Creator").
		Preload("Moderator").
		Preload("RequestExpense.Expense").
		Find(&calc).Error
	if err != nil {
		return nil, err
	}

	// Маппинг в DTO
	dtos := make([]ds.BreakevenRequestDTO, len(calc))
	for i, c := range calc {
		dto := ds.BreakevenRequestDTO{
			BreakevenRequestID:     c.BreakevenRequestID,
			BreakevenRequestStatus: c.BreakevenRequestStatus,
			CreationDate:           c.CreationDate,
			CreatorLogin:           c.Creator.Login,
			AmountProduct:          c.AmountProduct,
			CalcAnswer:             c.CalcAnswer,
		}
		if c.FormatedAt.Valid {
			dto.FormatedAt = &c.FormatedAt.Time
		}
		if c.CompletedAt.Valid {
			dto.CompletedAt = &c.CompletedAt.Time
		}
		if c.ModeratorID.Valid {
			moderatorLogin := c.Moderator.Login
			dto.ModeratorLogin = &moderatorLogin
		}
		for _, ce := range c.RequestExpense {
			dto.RequestExpense = append(dto.RequestExpense, ds.ExpenseForRequestDTO{
				ExpenseID:     ce.ExpenseID,
				Title:         ce.Expense.Title,
				ImageURL:      ce.Expense.ImageURL,
				AmountService: ce.AmountService,
				TypeSpend:     ce.TypeSpend,
			})
		}
		dtos[i] = dto
	}

	return dtos, nil
}

func (r *Repository) GetBreakevenCalcByID(CalcID int) (*ds.BreakevenRequestDTO, error) {
	var calc ds.BreakevenRequest
	err := r.db.Where(`"BreakevenRequestID" = ? AND "BreakEvenStatus" != 'удалён'`, CalcID).
		Preload("RequestExpense").
		Preload("RequestExpense.Expense").
		Preload("Creator").
		Preload("Moderator").
		First(&calc).Error
	if err != nil {
		return nil, err
	}
	dto := &ds.BreakevenRequestDTO{
		BreakevenRequestID:     calc.BreakevenRequestID,
		BreakevenRequestStatus: calc.BreakevenRequestStatus,
		CreationDate:           calc.CreationDate,
		CreatorLogin:           calc.Creator.Login,
		AmountProduct:          calc.AmountProduct,
		CalcAnswer:             calc.CalcAnswer,
	}
	if calc.FormatedAt.Valid {
		dto.FormatedAt = &calc.FormatedAt.Time
	}
	if calc.CompletedAt.Valid {
		dto.CompletedAt = &calc.CompletedAt.Time
	}
	if calc.ModeratorID.Valid {
		moderatorLogin := calc.Moderator.Login
		dto.ModeratorLogin = &moderatorLogin
	}
	for _, ce := range calc.RequestExpense {
		dto.RequestExpense = append(dto.RequestExpense, ds.ExpenseForRequestDTO{
			ExpenseID:     ce.ExpenseID,
			Title:         ce.Expense.Title,
			ImageURL:      ce.Expense.ImageURL,
			AmountService: ce.AmountService,
			TypeSpend:     ce.TypeSpend,
		})
	}
	return dto, err
}
func (r *Repository) UpdateBreakEvenCalc(id uint, breakevenRequestUpdate ds.UpdateBreakEvenCalcDTO) error {
	var BreakEvenCalc ds.BreakevenRequest
	err := r.db.Where(`"BreakevenRequestID" = ? AND "BreakEvenStatus" = 'черновик'`, id).First(&BreakEvenCalc).Error
	if err != nil {
		return err
	}
	if breakevenRequestUpdate.AmountProduct != 0 {
		BreakEvenCalc.AmountProduct = breakevenRequestUpdate.AmountProduct
	}

	return r.db.Save(&BreakEvenCalc).Error
}
func (r *Repository) FormByCreatorBreakEvenCalc(id uint) error {
	var BreakEvenCalc ds.BreakevenRequest
	err := r.db.Where(`"BreakevenRequestID" = ? AND "BreakEvenStatus" = 'черновик'`, id).First(&BreakEvenCalc).Error
	if err != nil {
		return nil
	}
	if BreakEvenCalc.AmountProduct == 0 {
		return fmt.Errorf("количество продуктов не может быть ноль")
	}
	var count int64
	err = r.db.Model(&ds.BreakevenRequest{}).Where(`"BreakevenRequestID" = ?`, id).Count(&count).Error
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("нельзя сформировать корзину без затрат")
	}
	BreakEvenCalc.BreakevenRequestStatus = "сформирован"
	BreakEvenCalc.FormatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	return r.db.Save(&BreakEvenCalc).Error
}
func (r *Repository) DeleteBreakEvenCalc(CalcID uint) error {
	return r.db.Exec(`Update Breakeven_Requests SET "BreakEvenStatus"= 'удалён' WHERE "BreakevenRequestID" = ?`, CalcID).Error
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

func (r *Repository) ProccessBreakEvenRequest(id uint, moderatorID int, action string) (*ds.BreakevenRequestDTO, error) {
	var CalcBreakEven ds.BreakevenRequest

	CalcBreakEven.ModeratorID = sql.NullInt64{Int64: int64(moderatorID), Valid: true}

	err := r.db.Where(`"BreakevenRequestID" = ?`, id).
		Preload("Moderator").
		Preload("Creator").
		Preload("RequestExpense").
		Preload("RequestExpense.Expense").
		First(&CalcBreakEven).
		Error
	if err != nil {
		return nil, err
	}

	if CalcBreakEven.BreakevenRequestStatus != "сформирован" {
		return nil, fmt.Errorf("заявка не может быть обработана, так как ожидается статус сформирован. Ee текущий статус: %s", CalcBreakEven.BreakevenRequestStatus)
	}
	CalcBreakEven.ModeratorID = sql.NullInt64{Int64: int64(moderatorID), Valid: true}

	switch action {
	case "complete":
		// Считаем ТБУ
		breakeven, err := service.CalculateAnswer(&CalcBreakEven)
		if err != nil {
			return nil, err
		}
		CalcBreakEven.CalcAnswer = breakeven
		CalcBreakEven.BreakevenRequestStatus = "завершён"
		CalcBreakEven.CompletedAt = sql.NullTime{Time: time.Now(), Valid: true}
	case "reject":
		CalcBreakEven.BreakevenRequestStatus = "отклонён"
	default:
		return nil, fmt.Errorf("действие %s недопустимо, допустимые действия 'complete' и 'reject'", action)
	}
	err = r.db.Save(&CalcBreakEven).Error
	if err != nil {
		return nil, err
	}
	var moderatorLogin string
	if CalcBreakEven.ModeratorID.Valid {
		// Прямой запрос логина пользователя по id
		err = r.db.Table("users").Where(`"UserID" = ?`, CalcBreakEven.ModeratorID.Int64).Select(`"Login"`).Scan(&moderatorLogin).Error
		if err != nil {
			moderatorLogin = ""
		}
	}
	dto := ds.BreakevenRequestDTO{
		BreakevenRequestID:     CalcBreakEven.BreakevenRequestID,
		BreakevenRequestStatus: CalcBreakEven.BreakevenRequestStatus,
		CreationDate:           CalcBreakEven.CreationDate,
		CreatorLogin:           CalcBreakEven.Creator.Login,
		AmountProduct:          CalcBreakEven.AmountProduct,
		CalcAnswer:             CalcBreakEven.CalcAnswer,
	}
	if CalcBreakEven.FormatedAt.Valid {
		dto.FormatedAt = &CalcBreakEven.FormatedAt.Time
	}
	if CalcBreakEven.CompletedAt.Valid {
		dto.CompletedAt = &CalcBreakEven.CompletedAt.Time
	}
	// Присваиваем логин
	if moderatorLogin == "" {
		dto.ModeratorLogin = &moderatorLogin
	}
	for _, ce := range CalcBreakEven.RequestExpense {
		dto.RequestExpense = append(dto.RequestExpense, ds.ExpenseForRequestDTO{
			ExpenseID:     ce.ExpenseID,
			Title:         ce.Expense.Title,
			ImageURL:      ce.Expense.ImageURL,
			AmountService: ce.AmountService,
			TypeSpend:     ce.TypeSpend,
		})
	}
	return &dto, nil
}
