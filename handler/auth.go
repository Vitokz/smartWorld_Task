package handler

import (
	"context"
	"crypto/sha1"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Vitokz/smartWorld_Task/internal/models"
	myJwt "github.com/Vitokz/smartWorld_Task/internal/services/jwt"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

const salt = "sdsaflmkjdsf"

func (h *Handler) Registration(c echo.Context) error {
	regUser := new(registrationRequest)

	err := c.Bind(regUser)
	if err != nil {
		h.log.WithFields(logrus.Fields{
			"event: ": "registration user",
			"err: ":   err,
			"time: ":  time.Now(),
		})
		return c.JSON(http.StatusBadRequest, NewErrorResponse(err.Error()))
	}

	regUser.Password = Hash(regUser.Password)

	err = regUser.validate()
	if err != nil {
		h.log.WithFields(logrus.Fields{
			"event: ": "registration user",
			"err: ":   err,
			"time: ":  time.Now(),
		}).Error()
		return c.JSON(http.StatusBadRequest, NewErrorResponse(err.Error()))
	}

	accToken, refreshToken, err := h.registration(context.Background(), *regUser)
	if err != nil {
		h.log.WithFields(logrus.Fields{
			"event: ": "registration user",
			"err: ":   err,
			"time: ":  time.Now(),
		}).Error()
		return c.JSON(http.StatusBadRequest, NewErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, tokensResponse{accToken, refreshToken})
}

func (h *Handler) Login(c echo.Context) error {
	logUser := new(loginUserRequest)

	err := c.Bind(logUser)
	if err != nil {
		h.log.WithFields(logrus.Fields{
			"event: ": "login user",
			"err: ":   err,
			"time: ":  time.Now(),
		}).Error()
		return c.JSON(http.StatusBadRequest, NewErrorResponse(err.Error()))
	}

	logUser.Password = Hash(logUser.Password)

	err = logUser.validate()
	if err != nil {
		h.log.WithFields(logrus.Fields{
			"event: ": "login user",
			"err: ":   err,
			"time: ":  time.Now(),
		}).Error()
		return c.JSON(http.StatusBadRequest, NewErrorResponse(err.Error()))
	}

	accToken, refreshToken, err := h.login(context.Background(), *logUser)
	if err != nil {
		h.log.WithFields(logrus.Fields{
			"event: ": "login user",
			"err: ":   err,
			"time: ":  time.Now(),
		}).Error()
		return c.JSON(http.StatusBadRequest, NewErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, tokensResponse{accToken, refreshToken})
}

func (h *Handler) Logout(c echo.Context) error {
	claims := new(myJwt.Claims)
	err := json.Unmarshal(c.Get("claims").([]byte), claims)
	if err != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(err.Error()))
	}

	id, err := strconv.Atoi(claims.Id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(err.Error()))
	}

	err = h.Repository.UpdateJWT(context.Background(), models.UserJwtTokens{
		UserID:       id,
		AccessToken:  "",
		RefreshToken: "s",
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(err.Error()))
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) registration(ctx context.Context, user registrationRequest) (string, string, error) {
	id, err := h.Repository.CreateUser(ctx, models.User{
		Name:     user.Name,
		Login:    user.Login,
		Password: user.Password,
		Role:     user.Role,
	})
	if err != nil {
		return "", "", err
	}

	accToken, refrreshToken, err := h.createTokens(user.Login, user.Name, user.Role, id)
	if err != nil {
		return "", "", err
	}

	err = h.Repository.InsertJWTToUser(ctx, models.UserJwtTokens{
		UserID:       id,
		AccessToken:  accToken,
		RefreshToken: refrreshToken,
	})
	if err != nil {
		return "", "", err
	}

	err = h.Repository.InsertUserReservationAcc(ctx, models.ReservationBook{
		UserID:           id,
		ReservationBooks: make([]int, 0),
	})
	if err != nil {
		return "", "", err
	}

	return accToken, refrreshToken, err
}

func (h *Handler) login(ctx context.Context, user loginUserRequest) (string, string, error) {
	dbUser, err := h.Repository.TakeUserByLogin(ctx, user.Login)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", errors.New("not found user with this login")
		}
		return "", "", err
	}

	if dbUser.Password != user.Password {
		return "", "", errors.New("incorrect password")
	}

	if dbUser.IsBlocked == true {
		return "", "", errors.New("you are blocked")
	}

	accToken, refreshToken, err := h.createTokens(user.Login, dbUser.Name, dbUser.Role, dbUser.ID)
	if err != nil {
		return "", "", err
	}

	err = h.Repository.UpdateJWT(ctx, models.UserJwtTokens{
		UserID:       dbUser.ID,
		AccessToken:  accToken,
		RefreshToken: refreshToken,
	})
	if err != nil {
		return "", "", err
	}

	return accToken, refreshToken, nil
}

func (h *Handler) createTokens(login, name, role string, id int) (string, string, error) {
	accToken, err := h.Jwt.CreateToken(login, name, role, id, time.Duration(h.Config.Jwt.AccessTimeExpired)*time.Minute)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := h.Jwt.CreateToken(login, name, role, id, time.Duration(h.Config.Jwt.RefreshTimeExpires)*time.Minute)
	if err != nil {
		return "", "", err
	}

	return accToken, refreshToken, nil
}

func (h *Handler) RefreshTokens(claims myJwt.Claims) (string, string, error) {
	id, err := strconv.Atoi(claims.Id)
	if err != nil {
		return "", "", err
	}

	accToken, refreshToken, err := h.createTokens(claims.Login, claims.Name, claims.Role, id)
	if err != nil {
		return "", "", err
	}

	err = h.Repository.UpdateJWT(context.Background(), models.UserJwtTokens{
		UserID:       id,
		AccessToken:  accToken,
		RefreshToken: refreshToken,
	})
	if err != nil {
		return "", "", err
	}

	return accToken, refreshToken, nil
}

func Hash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
