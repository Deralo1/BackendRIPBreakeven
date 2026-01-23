package repository

import (
	"Backeven/internal/app/ds"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
	// 🔹 Фильтрация по роли и userID
	if role == ds.RoleCreator && userID != 0 {
		query = query.Where(`"Creator_ID" = ?`, userID)
	}

	// 🔥 Модератор видит все заявки, кроме черновиков и удалённых
	// поэтому НИЧЕГО не добавляем

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
	err = r.db.Model(&ds.ExpenseForRequest{}).Where(`"BreakevenRequestID" = ?`, id).Count(&count).Error
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

	log.Println("=== START ProccessBreakEvenRequest ===")
	log.Printf("Input: id=%d, moderatorID=%d, action=%s\n", id, moderatorID, action)

	var req ds.BreakevenRequest

	// Загружаем заявку
	err := r.db.Where(`"BreakevenRequestID" = ?`, id).
		Preload("Creator").
		Preload("RequestExpense").
		Preload("RequestExpense.Expense").
		First(&req).Error

	if err != nil {
		log.Printf("ERROR: cannot load request: %v\n", err)
		return nil, err
	}

	log.Printf("Loaded request: ID=%d, Status=%s, AmountProduct=%d\n",
		req.BreakevenRequestID, req.BreakevenRequestStatus, req.AmountProduct)

	// Проверяем статус
	if req.BreakevenRequestStatus != "сформирован" {
		log.Printf("ERROR: wrong status: %s\n", req.BreakevenRequestStatus)
		return nil, fmt.Errorf("ожидается статус 'сформирован'")
	}

	// Обновляем модератора
	req.ModeratorID = sql.NullInt64{Int64: int64(moderatorID), Valid: true}

	switch action {

	case "complete":
		log.Println("Action: complete → preparing JSON for Rust")

		// Меняем статус на "в расчёте"
		req.BreakevenRequestStatus = "в расчёте"
		req.FormatedAt = sql.NullTime{Time: time.Now(), Valid: true}

		// Готовим JSON для Rust
		type RustExpense struct {
			TypeSpend     int `json:"type_spend"`
			AmountService int `json:"amount_service"`
			Expense       struct {
				Price int `json:"price"`
			} `json:"expense"`
		}

		rustExpenses := make([]RustExpense, 0)

		log.Println("Building expenses for Rust:")
		for _, e := range req.RequestExpense {
			log.Printf("  ExpenseID=%d, Type=%d, Amount=%d, Price=%d\n",
				e.ExpenseID, e.TypeSpend, e.AmountService, e.Expense.Price)

			re := RustExpense{
				TypeSpend:     e.TypeSpend,
				AmountService: e.AmountService,
			}
			re.Expense.Price = e.Expense.Price
			rustExpenses = append(rustExpenses, re)
		}

		rustReq := map[string]interface{}{
			"id":              req.BreakevenRequestID,
			"auth_token":      "secret123",
			"amount_product":  req.AmountProduct,
			"request_expense": rustExpenses,
		}

		jsonValue, _ := json.MarshalIndent(rustReq, "", "  ")
		log.Println("JSON sent to Rust:")
		log.Println(string(jsonValue))

		// Отправляем в Rust асинхронно
		go func() {
			resp, err := http.Post(
				"http://10.205.157.61:8083/calculateBreakeven",
				"application/json",
				bytes.NewBuffer(jsonValue),
			)
			if err != nil {
				log.Printf("ERROR sending to Rust: %v\n", err)
				return
			}
			log.Printf("Rust response status: %s\n", resp.Status)
		}()

	case "reject":
		log.Println("Action: reject")
		req.BreakevenRequestStatus = "отклонён"

	default:
		log.Printf("ERROR: invalid action: %s\n", action)
		return nil, fmt.Errorf("недопустимое действие")
	}

	// Сохраняем изменения
	if err := r.db.Save(&req).Error; err != nil {
		log.Printf("ERROR saving request: %v\n", err)
		return nil, err
	}

	log.Printf("Saved request: ID=%d, NewStatus=%s\n",
		req.BreakevenRequestID, req.BreakevenRequestStatus)
	// После сохранения — повторно загружаем заявку с расходами
	var updated ds.BreakevenRequest
	err = r.db.Where(`"BreakevenRequestID" = ?`, id).
		Preload("Creator").
		Preload("Moderator").
		Preload("RequestExpense").
		Preload("RequestExpense.Expense").
		First(&updated).Error
	if err != nil {
		return nil, err
	}

	// DTO
	dto := ds.BreakevenRequestDTO{
		BreakevenRequestID:     req.BreakevenRequestID,
		BreakevenRequestStatus: req.BreakevenRequestStatus,
		CreationDate:           req.CreationDate,
		CreatorLogin:           req.Creator.Login,
		AmountProduct:          req.AmountProduct,
		CalcAnswer:             req.CalcAnswer,
	}

	if req.FormatedAt.Valid {
		dto.FormatedAt = &req.FormatedAt.Time
	}
	if req.CompletedAt.Valid {
		dto.CompletedAt = &req.CompletedAt.Time
	}

	log.Println("=== END ProccessBreakEvenRequest ===")

	return &dto, nil
}

func (r *Repository) SaveBreakevenResult(id int, breakeven int) error {
	var req ds.BreakevenRequest

	if err := r.db.
		Where(`"BreakevenRequestID" = ?`, id).
		First(&req).Error; err != nil {
		return err
	}

	req.CalcAnswer = breakeven
	req.BreakevenRequestStatus = "завершён"
	req.CompletedAt = sql.NullTime{Time: time.Now(), Valid: true}

	return r.db.Save(&req).Error
}
