package applications

import "net/http"

func (tb *ToolBox) Register(w http.ResponseWriter, r *http.Request)           { return }
func (tb *ToolBox) Login(w http.ResponseWriter, r *http.Request)              { return }
func (tb *ToolBox) Logout(w http.ResponseWriter, r *http.Request)             { return }
func (tb *ToolBox) ConfirmEmail(w http.ResponseWriter, r *http.Request)       { return }
func (tb *ToolBox) ResendVerification(w http.ResponseWriter, r *http.Request) { return }
