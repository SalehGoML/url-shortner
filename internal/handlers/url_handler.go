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
	//userID := uint(1)
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

	//json.NewEncoder(w).Encode(url)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	encoder.Encode(url)
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

	//json.NewEncoder(w).Encode(urls)
	endcoder := json.NewEncoder(w)
	endcoder.SetIndent("", "  ")
	endcoder.Encode(urls)
}

func (h *URLHandler) Deactivate(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uint)

	var req struct {
		URLID uint `json:"url_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)

		return
	}

	urls, err := h.urlService.ListByUser(userID)
	if err != nil {
		http.Error(w, "failed to list urls", http.StatusInternalServerError)

		return
	}

	authorized := false
	for _, url := range urls {
		if url.ID == req.URLID {
			authorized = true

			break
		}
	}
	if !authorized {
		http.Error(w, "not authorized to deactivate this URL", http.StatusForbidden)

		return
	}

	err = h.urlService.Deactivate(req.URLID)
	if err != nil {
		http.Error(w, "failed to deactivate URL", http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "URL successfully deactivated successfully",
	})
}
