package handler

import "errors"

type registrationRequest struct {
	Name     string `json:"name"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func (r *registrationRequest) validate() error {
	switch {
	case r.Name == "":
		return errors.New("name is empty")
	case r.Login == "":
		return errors.New("login is empty")
	case r.Password == "":
		return errors.New("password is empty")
	case r.Role == "":
		return errors.New("role is empty")
	}

	return nil
}

type loginUserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (l *loginUserRequest) validate() error {
	switch {
	case l.Login == "":
		return errors.New("login is empty")
	case l.Password == "":
		return errors.New("password is empty")
	}

	return nil
}

type takeBookResponse struct {
	UserID           int   `json:"user_id"`
	ReservationBooks []int `json:"reservation_books"`
}

type tokensResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type ratingResponse struct {
	Name  string `json:"bookName"`
	Count int    `json:"bookCount"`
}

type ErrorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

func NewErrorResponse(mes string) ErrorResponse {
	return ErrorResponse{
		ErrorMessage: mes,
	}
}
