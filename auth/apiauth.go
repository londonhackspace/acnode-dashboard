package auth

import (
	"net/http"
	"os"
)

func apiAuth(r *http.Request) bool {
	apikey := r.Header.Get("API-KEY")

	return apikey == os.Getenv("API_KEY")
}
