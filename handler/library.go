package handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Vitokz/smartWorld_Task/internal/models"
	myJwt "github.com/Vitokz/smartWorld_Task/internal/services/jwt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func (h *Handler) TakeBook(c echo.Context) error {
	claims := new(myJwt.Claims)
	err := json.Unmarshal(c.Get("claims").([]byte), claims)
	if err != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(err.Error()))
	}

	bookId := c.Param("bookId")
	if bookId == "" {
		return c.JSON(http.StatusBadRequest, NewErrorResponse("book id is empty"))
	}

	bookIdInt, err := strconv.Atoi(bookId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(err.Error()))

	}

	userReservation, err := h.takeBook(*claims, bookIdInt)
	if err != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(err.Error()))
	}

	newAccessToken := c.Get("NewAccessToken")
	if newAccessToken != nil {
		resp := c.Response().Header()
		resp.Set("NewAccessToken", newAccessToken.(string))
		resp.Set("NewRefreshToken", c.Get("NewRefreshToken").(string))
	}

	return c.JSON(http.StatusOK, takeBookResponse{userReservation.UserID, userReservation.ReservationBooks})
}

func (h *Handler) takeBook(claims myJwt.Claims, bookId int) (models.ReservationBook, error) {
	book, err := h.Repository.TakeBookByID(context.Background(), bookId)
	if err != nil {
		return models.ReservationBook{}, err
	}

	if book.CountInLibrary <= 0 {
		return models.ReservationBook{}, errors.New("this book ended in library")
	}

	userId, err := strconv.Atoi(claims.Id)
	if err != nil {
		return models.ReservationBook{}, err
	}

	userResAcc, err := h.Repository.TakeUserReservationByID(context.Background(), userId)
	if err != nil {
		return models.ReservationBook{}, err
	}
	if len(userResAcc.ReservationBooks) >= 5 {
		return models.ReservationBook{}, errors.New("you have reached your book limit")
	}

	for _, v := range userResAcc.ReservationBooks {
		if v == bookId {
			return models.ReservationBook{}, errors.New("you have already reserve this book")
		}
	}

	userResAcc.ReservationBooks = append(userResAcc.ReservationBooks, bookId)

	err = h.Repository.UpdateUserReservationAcc(context.Background(), userResAcc)
	if err != nil {
		return models.ReservationBook{}, err
	}

	book.CountInLibrary--
	err = h.Repository.UpdateBook(context.Background(), book)

	err = h.Repository.InsertStoryField(context.Background(), models.StoryReservation{
		UserID: userId,
		BookID: bookId,
	})
	if err != nil {
		return models.ReservationBook{}, err
	}

	return userResAcc, nil
}

func (h *Handler) BlockUser(c echo.Context) error {
	claims := new(myJwt.Claims)
	err := json.Unmarshal(c.Get("claims").([]byte), claims)
	if err != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(err.Error()))
	}

	if claims.Role != "admin" {
		return c.JSON(http.StatusBadRequest, NewErrorResponse("you are not admin"))
	}

	userId := c.Param("userId")
	if userId == "" {
		return c.JSON(http.StatusBadRequest, NewErrorResponse("userId id is empty"))
	}

	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(err.Error()))
	}

	user, err := h.Repository.TakeUserByID(context.Background(), userIdInt)
	if err != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(err.Error()))
	}
	user.IsBlocked = true

	err = h.Repository.UpdateUser(context.Background(), user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(err.Error()))
	}

	newAccessToken := c.Get("NewAccessToken")
	if newAccessToken != nil {
		resp := c.Response().Header()
		resp.Set("NewAccessToken", newAccessToken.(string))
		resp.Set("NewRefreshToken", c.Get("NewRefreshToken").(string))
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) RatingAllTime(c echo.Context) error {
	resp, err := h.Repository.RatingAllTime(context.Background())
	if err != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(err.Error()))

	}

	result := make([]ratingResponse, len(resp))
	for i, v := range resp {
		result[i] = ratingResponse(v)
	}

	return c.JSON(http.StatusOK, result)
}

func (h *Handler) RatingReserved(c echo.Context) error {
	resp, err := h.Repository.RatingReserved(context.Background())
	if err != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(err.Error()))

	}

	result := make([]ratingResponse, len(resp))
	for i, v := range resp {
		result[i] = ratingResponse(v)
	}

	return c.JSON(http.StatusOK, result)
}
