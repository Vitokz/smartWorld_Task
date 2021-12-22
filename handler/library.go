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

	bookID := c.Param("bookId")
	if bookID == "" {
		return c.JSON(http.StatusBadRequest, NewErrorResponse("book id is empty"))
	}

	bookIDInt, err := strconv.Atoi(bookID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(err.Error()))
	}

	userReservation, err := h.takeBook(*claims, bookIDInt)
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

func (h *Handler) takeBook(claims myJwt.Claims, bookID int) (models.ReservationBook, error) {
	book, err := h.Repository.TakeBookByID(context.Background(), bookID)
	if err != nil {
		return models.ReservationBook{}, err
	}

	if book.CountInLibrary <= 0 {
		return models.ReservationBook{}, errors.New("this book ended in library")
	}

	userID, err := strconv.Atoi(claims.Id)
	if err != nil {
		return models.ReservationBook{}, err
	}

	userResAcc, err := h.Repository.TakeUserReservationByID(context.Background(), userID)
	if err != nil {
		return models.ReservationBook{}, err
	}

	if len(userResAcc.ReservationBooks) >= 5 {
		return models.ReservationBook{}, errors.New("you have reached your book limit")
	}

	for _, v := range userResAcc.ReservationBooks {
		if v == bookID {
			return models.ReservationBook{}, errors.New("you have already reserve this book")
		}
	}

	userResAcc.ReservationBooks = append(userResAcc.ReservationBooks, bookID)

	err = h.Repository.UpdateUserReservationAcc(context.Background(), userResAcc)
	if err != nil {
		return models.ReservationBook{}, err
	}

	book.CountInLibrary--
	err = h.Repository.UpdateBook(context.Background(), book)

	if err != nil {
		return models.ReservationBook{}, err
	}

	err = h.Repository.InsertStoryField(context.Background(), models.StoryReservation{
		UserID: userID,
		BookID: bookID,
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

	userID := c.Param("userId")
	if userID == "" {
		return c.JSON(http.StatusBadRequest, NewErrorResponse("userId id is empty"))
	}

	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(err.Error()))
	}

	user, err := h.Repository.TakeUserByID(context.Background(), userIDInt)
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
