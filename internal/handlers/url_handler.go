package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/SalehGoML/internal/middleware"
	"github.com/SalehGoML/internal/service"
)

type URLHandler struct {
	urlService service.URLService
}

func NewURLHandler(urlService service.URLService) *URLHandler {
	return &URLHandler{urlService: urlService}
}

func (h *URLHandler) Shorten(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uint)

	var req struct {
		LongURL string     `json:"long_url"`
		Expiry  *time.Time `json:"expiry,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)

		return
	}

	url, err := h.urlService.Shorten(req.LongURL, userID, req.Expiry)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	json.NewEncoder(w).Encode(url)
}

func (h *URLHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Path[1:]

	url, err := h.urlService.GetByCode(code)
	if err != nil {
		http.NotFound(w, r)

		return
	}

	_ = h.urlService.IncrementClicks(url.ID)

	http.Redirect(w, r, url.LongURL, http.StatusFound)
}
