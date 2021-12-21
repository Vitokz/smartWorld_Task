package repository

import (
	"context"
	"fmt"
	"github.com/Vitokz/smartWorld_Task/config"
	"github.com/Vitokz/smartWorld_Task/internal/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type repository struct {
	PostgresDB *sqlx.DB
}

type Repository interface {
	CreateUser(ctx context.Context, user models.User) (int, error)
	TakeUserByID(ctx context.Context, id int) (models.User, error)
	TakeUserByLogin(ctx context.Context, login string) (models.User, error)
	UpdateUser(ctx context.Context, user models.User) error
	DeleteUserByID(ctx context.Context, id int) error

	InsertJWTToUser(ctx context.Context, tokens models.UserJwtTokens) error
	UpdateJWT(ctx context.Context, tokens models.UserJwtTokens) error
	TakeUserJWTByID(ctx context.Context, userID int) (models.UserJwtTokens, error)
	DeleteUserJWT(ctx context.Context, userID int) error

	InsertUserReservationAcc(ctx context.Context, reservation models.ReservationBook) error
	UpdateUserReservationAcc(ctx context.Context, reservation models.ReservationBook) error
	TakeUserReservationByID(ctx context.Context, userID int) (models.ReservationBook, error)
	DeleteUserReservationAcc(ctx context.Context, userID int) error

	CreateBook(ctx context.Context, book models.Book) error
	DeleteBook(ctx context.Context, bookID int) error
	TakeBookByID(ctx context.Context, bookID int) (models.Book, error)
	UpdateBook(ctx context.Context, book models.Book) error

	InsertStoryField(ctx context.Context, story models.StoryReservation) error
	DeleteBookReserveStory(ctx context.Context, bookID int) error
	UpdateReturnedParamStory(ctx context.Context, bookID, userID int) error
	TakeBookReservedStory(ctx context.Context, bookID int) ([]models.StoryReservation, error)
	TakeUserReservedStory(ctx context.Context, userID int) ([]models.StoryReservation, error)

	RatingAllTime(ctx context.Context) ([]models.Rating, error)
	RatingReserved(ctx context.Context) ([]models.Rating, error)
}

func New(postgresDB *sqlx.DB) *repository {
	return &repository{
		PostgresDB: postgresDB,
	}
}

func NewPgSQL(ctx context.Context, cfg *config.Config) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", cfg.Postgres.Dsn)
	if err != nil {
		return nil, fmt.Errorf("postgres: sqlx.Connect: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("postgres: sqlx.Ping: %w", err)
	}

	return db, nil
}
