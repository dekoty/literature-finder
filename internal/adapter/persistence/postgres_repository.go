package persistence

import (
	"database/sql"
	"fmt"
	"literature-finder/internal/module/literature"
	"log"
	"strings"
)

type userBookDB struct {
	BookID    string `db:"book_id"`
	Title     string `db:"title"`
	Authors   string `db:"authors"`
	Thumbnail string `db:"thumbnail"`
	Link      string `db:"link"`
	Status    string `db:"status"`
	Year      string `db:"year"`
}

type PostgresRepository struct {
	DB *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{DB: db}
}

func (r *PostgresRepository) SaveBook(userID string, book literature.Literature) error {
	authorsStr := strings.Join(book.Authors, "; ")

	query := `INSERT INTO user_books (
				user_id, book_id, title, authors, thumbnail, link, status, year
			  ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			  ON CONFLICT (user_id, book_id, status) DO NOTHING;`

	_, err := r.DB.Exec(
		query,
		userID,
		book.ID,
		book.Title,
		authorsStr,
		book.Thumbnail,
		book.Link,
		book.Status,
		book.Year,
	)

	if err != nil {
		return fmt.Errorf("ошибка при сохранении книги (ID: %s) для пользователя %s: %w", book.ID, userID, err)
	}

	return nil
}

func (r *PostgresRepository) GetBooksByUserID(userID string, status string) ([]literature.Literature, error) {
	query := `SELECT book_id, title, authors, thumbnail, link, status, year
			  FROM user_books
			  WHERE user_id = $1 AND status = $2
			  ORDER BY created_at DESC`

	rows, err := r.DB.Query(query, userID, status)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса GetBooksByUserID для пользователя %s: %w", userID, err)
	}
	defer rows.Close()

	var books []literature.Literature

	for rows.Next() {
		var dbBook userBookDB

		err := rows.Scan(
			&dbBook.BookID,
			&dbBook.Title,
			&dbBook.Authors,
			&dbBook.Thumbnail,
			&dbBook.Link,
			&dbBook.Status,
			&dbBook.Year,
		)

		if err != nil {
			log.Printf("Внимание: ошибка сканирования строки БД для пользователя %s: %v", userID, err)
			return nil, fmt.Errorf("ошибка чтения данных из БД: %w", err)
		}

		book := literature.Literature{
			ID:        dbBook.BookID,
			Title:     dbBook.Title,
			Authors:   strings.Split(dbBook.Authors, "; "),
			Thumbnail: dbBook.Thumbnail,
			Link:      dbBook.Link,
			Status:    dbBook.Status,
			Year:      dbBook.Year,
		}
		books = append(books, book)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при итерации строк таблицы user_books: %w", err)
	}

	return books, nil
}

func (r *PostgresRepository) DeleteBook(userID string, bookID string) error {
	query := `
        DELETE FROM user_books
        WHERE user_id = $1 AND book_id = $2;
    `

	result, err := r.DB.Exec(query, userID, bookID)
	if err != nil {
		return fmt.Errorf("ошибка выполнения запроса удаления книги %s: %w", bookID, err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		log.Printf("Попытка удалить несуществующую книгу %s у пользователя %s", bookID, userID)
	} else {
		log.Printf("Книга %s успешно удалена для пользователя %s", bookID, userID)
	}

	return nil
}

func (r *PostgresRepository) ClearFavorites(userID string) error {
	query := `
        DELETE FROM user_books
        WHERE user_id = $1;
    `

	_, err := r.DB.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("ошибка очистки избранного для пользователя %s: %w", userID, err)
	}
	log.Printf("Список избранного успешно очищен для пользователя %s", userID)
	return nil
}
