package handlers

import (
	"net/http"

	"github.com/erodrigufer/serenitynow/internal/views"
	"github.com/erodrigufer/serenitynow/internal/web"
)

func (h *Handlers) HandleGetHome() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := web.RenderComponent(r.Context(), w, views.Home())
		if err != nil {
			web.HandleServerError(w, r, err, h.errorLog)
		}
	}
}
