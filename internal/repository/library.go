package repository

import (
	"context"
	"github.com/Vitokz/smartWorld_Task/internal/models"
	"time"
)

func (r *repository) CreateBook(ctx context.Context, book models.Book) error {
	insertQuery := `INSERT INTO books (author, name, count_in_library, updated_at) VALUES ($1, $2, $3, $4)`

	_, err := r.PostgresDB.ExecContext(ctx, insertQuery, book.Author, book.Name, book.CountInLibrary, time.Now().Unix())
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) DeleteBook(ctx context.Context, bookID int) error {
	_, err := r.PostgresDB.ExecContext(ctx, "DELETE FROM books WHERE id = $1", bookID)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) TakeBookByID(ctx context.Context, bookID int) (models.Book, error) {
	var res models.Book
	if err := r.PostgresDB.GetContext(ctx, &res,
		`SELECT * FROM books WHERE id = $1`, bookID); err != nil {
		return models.Book{}, err
	}

	return res, nil
}

func (r *repository) UpdateBook(ctx context.Context, book models.Book) error {
	_, err := r.PostgresDB.ExecContext(ctx,
		`UPDATE books 
    		SET author = $1, name = $2, count_in_library =$3, updated_at = $4
			WHERE id = $5`, book.Author, book.Name, book.CountInLibrary, time.Now().Unix(), book.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) InsertStoryField(ctx context.Context, story models.StoryReservation) error {
	insertQuery := `INSERT INTO books_reservation_story (id_book, id_user,reserved_at, returned_at) VALUES ($1, $2, $3, $4)`

	_, err := r.PostgresDB.ExecContext(ctx, insertQuery, story.BookID, story.UserID, time.Now().Unix(), nil)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) DeleteBookReserveStory(ctx context.Context, bookID int) error {
	_, err := r.PostgresDB.ExecContext(ctx, "DELETE FROM books_reservation_story WHERE id_book = $1", bookID)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) UpdateReturnedParamStory(ctx context.Context, bookID, userID int) error {
	_, err := r.PostgresDB.ExecContext(ctx,
		`UPDATE books_reservation_story 
    		SET returned_at = $1
			WHERE id_book = $2 AND id_user = $3 AND returned_at IS NULL`, time.Now().Unix(), bookID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) TakeBookReservedStory(ctx context.Context, bookID int) ([]models.StoryReservation, error) {
	var res = make([]models.StoryReservation, 0)
	if err := r.PostgresDB.SelectContext(ctx, &res,
		`SELECT * FROM books_reservation_story WHERE id_book = $1`, bookID); err != nil {
		return nil, err
	}

	return res, nil
}

func (r *repository) TakeUserReservedStory(ctx context.Context, userID int) ([]models.StoryReservation, error) {
	var res = make([]models.StoryReservation, 0)
	if err := r.PostgresDB.SelectContext(ctx, &res,
		`SELECT * FROM books_reservation_story WHERE id_user = $1`, userID); err != nil {
		return nil, err
	}

	return res, nil
}

func (r *repository) RatingAllTime(ctx context.Context) ([]models.Rating, error) {
	result := make([]models.Rating, 0)

	if err := r.PostgresDB.SelectContext(ctx, &result,
		`select b.name as book_name, COUNT(brs.id_book) as book_count from books_reservation_story brs
				left join books b on b.id = brs.id_book
				group by (b.name,brs.id_book) order by COUNT(brs.id_book) DESC limit 10`); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *repository) RatingReserved(ctx context.Context) ([]models.Rating, error) {
	result := make([]models.Rating, 0)

	if err := r.PostgresDB.SelectContext(ctx, &result,
		`select b.name as book_name, COUNT(brs.id_book) as book_count from books_reservation_story brs
				left join books b on b.id = brs.id_book
				where brs.returned_at is null
				group by (b.name,brs.id_book) order by COUNT(brs.id_book) desc limit 10`); err != nil {
		return nil, err
	}

	return result, nil
}
