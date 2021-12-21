package models

import (
	"fmt"
	"strconv"
	"strings"
)

type User struct {
	ID        int    `db:"id"`
	Name      string `db:"name"`
	Login     string `db:"login"`
	Password  string `db:"password"`
	Role      string `db:"role"`
	CreatedAt int    `db:"created_at"`
	IsBlocked bool   `db:"is_blocked"`
}

type Book struct {
	ID             int    `db:"id"`
	Author         string `db:"author"`
	Name           string `db:"name"`
	CountInLibrary int    `db:"count_in_library"`
	UpdatedAt      int    `db:"updated_at"`
}

type UserJwtTokens struct {
	UserID       int    `db:"id_user"`
	AccessToken  string `db:"access_token"`
	RefreshToken string `db:"refresh_token"`
	UpdatedAt    int    `db:"updated_at"`
}

type Rating struct {
	Name  string `db:"book_name"`
	Count int    `db:"book_count"`
}

type StoryReservation struct {
	BookID     int  `db:"id_book"`
	UserID     int  `db:"id_user"`
	ReservedAt int  `db:"reserved_at"`
	ReturnedAt *int `db:"returned_at"`
}

type ReservationBook struct {
	UserID           int              `db:"id_user"`
	ReservationBooks reservationArray `db:"reservation_books" json:"reservation_books"`
}

type reservationArray []int

func (j *reservationArray) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	str := fmt.Sprintf(string(src.([]byte)))
	str = strings.TrimPrefix(str, "{")
	str = strings.TrimSuffix(str, "}")

	if str == "" {
		return nil
	}

	ints := strings.Split(str, ",")
	if len(ints) == 0 {
		return nil
	}
	numbers := make([]int, 0)
	for _, v := range ints {
		number, err := strconv.Atoi(v)
		if err != nil {
			return err
		}
		numbers = append(numbers, number)
	}
	result := reservationArray(numbers)
	*j = result
	return nil
}
