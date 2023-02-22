package api

import (
	"bjssStoreGo/backend/utils"
	"encoding/json"
	"net/http"

	"github.com/gorilla/sessions"
)

func Error(w http.ResponseWriter, r *http.Request, status int, form string, msg string) {
	w.WriteHeader(http.StatusUnauthorized)

	json.NewEncoder(w).Encode(utils.ApiErrorResponse{
		Error: form,
		Msg:   msg,
	})
}

func GetSignedInUserId(r *http.Request, s *sessions.CookieStore) int {
	session, _ := s.Get(r, "session-name")
	customerId, ok := session.Values["customerId"]

	if !ok {
		return 0
	}

	return customerId.(int)
}
