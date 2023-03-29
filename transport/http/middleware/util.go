package middleware

import (
	"fmt"
	"net/http"
	"strings"
)

func getHeaders(h http.Header) string {
	str := ""
	for key, value := range h {
		str = fmt.Sprintln(str, key+": "+value[0])
	}

	return str
}

func getIP(r *http.Request) string {
	realIP := r.Header.Get("X-REAL-IP")
	if realIP != "" {
		return realIP
	}

	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}

	return strings.Split(r.RemoteAddr, ":")[0]
}
