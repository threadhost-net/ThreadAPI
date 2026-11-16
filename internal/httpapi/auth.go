package httpapi

import (
	"encoding/json"
	"net"
	"net/http"
)

func register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Некорректный запрос", http.StatusBadRequest)
		return
	}

	result, err := toolbox.Register(r.Context(), req.Username, req.Password, req.Email, req.HCaptchaResponse)
	if err != nil {
		http.Error(w, err.Error(), result)
		return
	}

	w.WriteHeader(result)
}

func login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Некорректный запрос", http.StatusBadRequest)
		return
	}

	address, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		address = r.RemoteAddr
	}

	userAgent := parseUserAgent(r.Header.Get("User-Agent"))

	result, cookieToken, err := toolbox.Login(r.Context(), req.Username, req.Password, req.HCaptchaResponse, address, userAgent)
	if err != nil {
		http.Error(w, err.Error(), result)
		return
	}

	setSessionCookie(w, cookieToken)
	w.WriteHeader(result)
}
