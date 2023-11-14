package books

import (
	"context"
	"service-weaver-test/models"

	"github.com/ServiceWeaver/weaver"
	"github.com/google/uuid"
)

var bookStore = make(map[string]models.Book)

type BookService interface {
	// Get returns the requested book with the given ID.
	Get(ctx context.Context, id uuid.UUID) (models.Book, error)
	// Create saves a new book in the storage.
	Create(ctx context.Context, b models.Book) error
}

type bookService struct {
	weaver.Implements[BookService]
}

// Get returns the requested book with the given ID.
func (b bookService) Get(ctx context.Context, id uuid.UUID) (models.Book, error) {
	b.Logger(ctx).InfoContext(ctx, "Getting book with ID: ", id.String(), "...")
	return bookStore[id.String()], nil
}

// Create saves a new book in the storage.
func (b bookService) Create(ctx context.Context, book models.Book) error {
	b.Logger(ctx).InfoContext(ctx, "Creating book: ", book.Title, "...")
	uuid, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	book.ID = uuid
	bookStore[book.ID.String()] = book
	return nil
}
