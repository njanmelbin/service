package web

import (
	"context"
	"fmt"
	"net/http"
	"os"
)

type HandlerFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

type App struct {
	*http.ServeMux
	shutdown chan os.Signal
}

func New(shutdown chan os.Signal) *App {
	return &App{
		ServeMux: http.NewServeMux(),
		shutdown: shutdown,
	}
}

func (a *App) HandleFunc(pattern string, handler HandlerFunc) {

	h := func(w http.ResponseWriter, r *http.Request) {
		if err := handler(r.Context(), w, r); err != nil {
			fmt.Println(err)
			return
		}
	}
	a.ServeMux.HandleFunc(pattern, h)
}
