package applications

import "net/http"

func (tb *ToolBox) ForgotPassword(w http.ResponseWriter, r *http.Request)        { return }
func (tb *ToolBox) ConfirmForgotPassword(w http.ResponseWriter, r *http.Request) { return }
func (tb *ToolBox) ConfirmNewEmail(w http.ResponseWriter, r *http.Request)       { return }
func (tb *ToolBox) UserInfo(w http.ResponseWriter, r *http.Request)              { return }
func (tb *ToolBox) ChangeNewPassword(w http.ResponseWriter, r *http.Request)     { return }
func (tb *ToolBox) ChangeNewEmail(w http.ResponseWriter, r *http.Request)        { return }
func (tb *ToolBox) UserDelete(w http.ResponseWriter, r *http.Request)            { return }
