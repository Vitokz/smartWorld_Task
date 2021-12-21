package repository

import (
	"context"
	"github.com/Vitokz/smartWorld_Task/internal/models"
	"github.com/lib/pq"
	"time"
)

func (r *repository) CreateUser(ctx context.Context, user models.User) (int, error) {
	insertQuery := `INSERT INTO users (name, login, password, role, created_at) VALUES ($1, $2, $3, $4,$5)`

	_, err := r.PostgresDB.ExecContext(ctx, insertQuery, user.Name, user.Login, user.Password, user.Role, time.Now().Unix())
	if err != nil {
		return 0, err
	}

	newUser, err := r.TakeUserByLogin(ctx, user.Login)
	if err != nil {
		return 0, err
	}

	return newUser.ID, nil
}

func (r *repository) TakeUserByID(ctx context.Context, id int) (models.User, error) {
	var res models.User
	if err := r.PostgresDB.GetContext(ctx, &res,
		`SELECT * FROM users WHERE id = $1`, id); err != nil {
		return models.User{}, err
	}

	return res, nil
}

func (r *repository) TakeUserByLogin(ctx context.Context, login string) (models.User, error) {
	var res models.User
	if err := r.PostgresDB.GetContext(ctx, &res,
		`SELECT * FROM users WHERE login = $1`, login); err != nil {
		return models.User{}, err
	}

	return res, nil
}

func (r *repository) UpdateUser(ctx context.Context, user models.User) error {
	_, err := r.PostgresDB.ExecContext(ctx,
		`UPDATE users 
    		SET name = $1, login = $2, password = $3, role = $4, is_blocked = $5
			WHERE id = $6`, user.Name, user.Login, user.Password, user.Role, user.IsBlocked, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) DeleteUserByID(ctx context.Context, id int) error {
	tx, err := r.PostgresDB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM users_jwt_tokens WHERE id_user = $1", id)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM reservation_books WHERE id_user = $1", id)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return nil
}

func (r *repository) InsertJWTToUser(ctx context.Context, tokens models.UserJwtTokens) error {
	insertQuery := `INSERT INTO users_jwt_tokens (id_user, access_token, refresh_token, updated_at) 
					VALUES ($1,$2,$3,$4)`

	_, err := r.PostgresDB.ExecContext(ctx, insertQuery, tokens.UserID, tokens.AccessToken, tokens.RefreshToken, time.Now().Unix())
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) UpdateJWT(ctx context.Context, tokens models.UserJwtTokens) error {
	_, err := r.PostgresDB.ExecContext(ctx,
		`UPDATE users_jwt_tokens 
    		SET access_token = $1, refresh_token = $2, updated_at = $3
			WHERE id_user = $4`, tokens.AccessToken, tokens.RefreshToken, time.Now().Unix(), tokens.UserID)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) TakeUserJWTByID(ctx context.Context, userID int) (models.UserJwtTokens, error) {
	var res models.UserJwtTokens
	if err := r.PostgresDB.GetContext(ctx, &res,
		`SELECT * FROM users_jwt_tokens WHERE id_user = $1`, userID); err != nil {
		return models.UserJwtTokens{}, err
	}

	return res, nil
}

func (r *repository) DeleteUserJWT(ctx context.Context, userID int) error {
	_, err := r.PostgresDB.ExecContext(ctx, "DELETE FROM users_jwt_tokens WHERE id_user = $1", userID)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) InsertUserReservationAcc(ctx context.Context, reservation models.ReservationBook) error {
	insertQuery := `INSERT INTO reservation_books (id_user, reservation_books) 
					VALUES ($1,$2)`

	_, err := r.PostgresDB.ExecContext(ctx, insertQuery, reservation.UserID, pq.Array(reservation.ReservationBooks))
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) UpdateUserReservationAcc(ctx context.Context, reservation models.ReservationBook) error {
	_, err := r.PostgresDB.ExecContext(ctx,
		`UPDATE reservation_books 
    		SET reservation_books = $1
			WHERE id_user = $2`, pq.Array(reservation.ReservationBooks), reservation.UserID)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) TakeUserReservationByID(ctx context.Context, userID int) (models.ReservationBook, error) {
	var res models.ReservationBook
	if err := r.PostgresDB.GetContext(ctx, &res,
		`SELECT * FROM reservation_books WHERE id_user = $1`, userID); err != nil {
		return models.ReservationBook{}, err
	}

	return res, nil
}

func (r *repository) DeleteUserReservationAcc(ctx context.Context, userID int) error {
	_, err := r.PostgresDB.ExecContext(ctx, "DELETE FROM reservation_books WHERE id_user = $1", userID)
	if err != nil {
		return err
	}

	return nil
}
