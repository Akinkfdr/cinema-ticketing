package httpapi

import "net/http"

// health yalnız API'yi kontrol eder; PostgreSQL veya Redis bağlantısını sınamaz.
func health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("{\"status\":\"ok\"}\n"))
}
