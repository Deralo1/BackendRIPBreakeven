package handler

import (
	"Backeven/internal/app/ds"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) RegisterUser(ctx *gin.Context) {
	var input ds.ChangeUserDTO

	if err := ctx.BindJSON(&input); err != nil {
		h.errorhandler(ctx, http.StatusBadRequest, err)
		return
	}

	if input.Login == "" || input.Password == "" {
		h.errorhandler(ctx, http.StatusBadRequest, fmt.Errorf("введите логин и пароль"))
		return
	}

	err := h.Repository.RegisterUser(input)
	if err != nil {
		h.errorhandler(ctx, http.StatusInternalServerError, err)
		return
	}
	h.successResponse(ctx, gin.H{})
}

func (h *Handler) LoginUser(ctx *gin.Context) {
	var input ds.ChangeUserDTO
	if err := ctx.BindJSON(&input); err != nil {
		h.errorhandler(ctx, http.StatusBadRequest, err)
		return
	}

	if input.Login == "" || input.Password == "" {
		h.errorhandler(ctx, http.StatusBadRequest, fmt.Errorf("введите логин и пароль"))
		return
	}

	userdto, err := h.Repository.LoginUser(input.Login, input.Password)
	if err != nil {
		h.errorhandler(ctx, http.StatusUnauthorized, err)
		return
	}
	h.successResponse(ctx, userdto)
}

func (h *Handler) GetProfile(ctx *gin.Context) {
	userID := h.GetCurrentUserId()
	userdto, err := h.Repository.GetUserByID(userID)
	if err != nil {
		h.errorhandler(ctx, http.StatusInternalServerError, err)
		return
	}
	h.successResponse(ctx, userdto)
}

func (h *Handler) UpdateUserProf(ctx *gin.Context) {

	userID := h.GetCurrentUserId()
	var userUpdate ds.ChangeUserDTO

	if err := ctx.BindJSON(&userUpdate); err != nil {
		h.errorhandler(ctx, http.StatusBadRequest, err)
		return
	}

	updateduserDto, err := h.Repository.UpdateUser(userID, userUpdate)
	if err != nil {
		h.errorhandler(ctx, http.StatusInternalServerError, err)
		return
	}

	h.successResponse(ctx, updateduserDto)
}

func (h *Handler) LogoutUser(ctx *gin.Context) {
	userid := h.GetCurrentUserId()

	err := h.Repository.LogoutUser(userid)
	if err != nil {
		h.errorhandler(ctx, http.StatusInternalServerError, err)
		return
	}

	h.successResponse(ctx, gin.H{
		"message": "Успешный выход из системы",
	})
}
