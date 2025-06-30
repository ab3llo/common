package common

import "net/http"

func GetQueryParam(r *http.Request, key string) string {
	value := r.URL.Query().Get(key)
	if value == "" {
		return ""
	}
	return value
}
