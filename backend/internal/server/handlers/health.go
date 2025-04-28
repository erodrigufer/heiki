package handlers

import (
	"net/http"

	"github.com/erodrigufer/serenitynow/internal/web"
)

func (h *Handlers) GetHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := map[string]string{
			"Status": "ok",
		}

		web.SendJSONResponse(w, http.StatusOK, response)
	}
}
