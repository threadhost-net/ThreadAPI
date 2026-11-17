package httpapi

import (
	"encoding/json"
	"net"
	"net/http"

	"github.com/threadhost-net/ThreadAPI/cookie"
	"github.com/threadhost-net/ThreadAPI/useragent"
)

func register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Некорректный запрос", http.StatusBadRequest)
		return
	}

	status, err := toolbox.Register(r.Context(), req.Username, req.Password, req.Email, req.HCaptchaResponse)
	if err != nil {
		http.Error(w, err.Error(), status)
		return
	}

	w.WriteHeader(status)
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

	userAgent := useragent.Parse(r.Header.Get("User-Agent"))

	status, cookieToken, err := toolbox.Login(r.Context(), req.Username, req.Password, req.HCaptchaResponse, address, userAgent)
	if err != nil {
		http.Error(w, err.Error(), status)
		return
	}

	cookie.SetSessionCookie(w, cookieToken)
	w.WriteHeader(status)
}

func logout(w http.ResponseWriter, r *http.Request) {
	cookieToken, err := r.Cookie("session")
	if err != nil {
		cookie.ClearSessionCookie(w)
		w.WriteHeader(http.StatusAccepted)
		return
	}

	status, err := toolbox.Logout(r.Context(), cookieToken.Value)
	if err != nil {
		http.Error(w, err.Error(), status)
	}

}

func emailVerify(w http.ResponseWriter, r *http.Request) {
	return
}
