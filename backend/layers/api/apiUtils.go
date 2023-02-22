package api

import (
	"bjssStoreGo/backend/utils"
	"encoding/json"
	"net/http"
)

func Error(w http.ResponseWriter, r *http.Request, status int, form string, msg string) {
	w.WriteHeader(http.StatusUnauthorized)

	json.NewEncoder(w).Encode(utils.ApiErrorResponse{
		Error: form,
		Msg:   msg,
	})
}
