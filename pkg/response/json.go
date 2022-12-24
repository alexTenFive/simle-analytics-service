package response

import (
	"encoding/json"
	"net/http"
	"testcase_v2/pkg/httpErrors"
)

func JsonOk(w http.ResponseWriter) {
	Json(w, map[string]interface{}{
		"status": "ok",
	})
}
func Json(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, httpErrors.ErrBadResult.Error(), http.StatusInternalServerError)
	}
}
