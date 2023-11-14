package main

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"service-weaver-test/internal/books"
	"service-weaver-test/models"

	"github.com/ServiceWeaver/weaver"
	"github.com/google/uuid"
)

//go:generate weaver generate ./...

func main() {
	if err := weaver.Run(context.Background(), server); err != nil {
		log.Fatal(err)
	}
}

type app struct {
	weaver.Implements[weaver.Main]
	bookService weaver.Ref[books.BookService]
	lis         weaver.Listener `weaver:"books"`
}

func server(ctx context.Context, a *app) error {
	a.Logger(ctx).InfoContext(ctx, "Starting server...")
	mux := http.NewServeMux()

	mux.Handle("/books", weaver.InstrumentHandlerFunc("books", func(w http.ResponseWriter, r *http.Request) {
		bookservice := a.bookService.Get()
		switch r.Method {
		case http.MethodGet:
			// Get the book with the given ID.
			id := r.URL.Query().Get("id")
			uuid, err := uuid.Parse(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			book, err := bookservice.Get(ctx, uuid)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			json, err := json.Marshal(book)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write(json)
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")

		case http.MethodPost:
			// Create a new book.
			var book models.Book
			err := json.NewDecoder(r.Body).Decode(&book)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			err = bookservice.Create(ctx, book)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusCreated)
			w.Header().Set("location", "/books?id="+book.ID.String())
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.HandleFunc(weaver.HealthzURL, weaver.HealthzHandler)

	a.Logger(ctx).InfoContext(ctx, "Listening on ", slog.String("on ", a.lis.Addr().String()))
	return http.Serve(a.lis, mux)
}
