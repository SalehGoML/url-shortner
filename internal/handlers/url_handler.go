package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
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
	//userID := r.Context().Value(middleware.UserIDKey).(uint)
	userID := uint(1)
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

func (h *URLHandler) ListMyURLs(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uint)

	urls, err := h.urlService.ListByUser(userID)
	if err != nil {
		http.Error(w, "failed to list urls", http.StatusInternalServerError)

		return
	}

	json.NewEncoder(w).Encode(urls)
}

func (h *URLHandler) Deactivate(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "missing id", http.StatusBadRequest)

		return
	}

	urlID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)

		return
	}

	err = h.urlService.Deactivate(uint(urlID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("lind deactivated"))
}
