package restapi

import (
	"KursachMessenger/api"
	"KursachMessenger/database"
	"encoding/json"
	"fmt"
	"net/http"
)

var CreateAccount = func(w http.ResponseWriter, r *http.Request) {
	account := &database.Account{}
	err := json.NewDecoder(r.Body).Decode(account)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		api.Respond(w, api.MakeResponse(false, "ERR11"))
		return
	}
	response, ok := account.Create()
	if ok == true{
		fmt.Printf("Новый пользователь: %s \n", account.Email)
		api.Respond(w,response)
	} else{
		w.WriteHeader(http.StatusBadRequest)
		api.Respond(w,response)
	}

	return
}
var Login = func(w http.ResponseWriter, r *http.Request){
	account := &database.Account{}
	err := json.NewDecoder(r.Body).Decode(account)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		api.Respond(w,api.MakeResponse(false, "ERR12"))
		return
	}
	response, ok := database.Login(account.Email,account.Password)
	if ok == true {
		api.Respond(w, response)
		fmt.Printf("Токен получен пользователем: %s \n", account.Email)
		return
	} else{
		w.WriteHeader(http.StatusForbidden)
		api.Respond(w,response)
		fmt.Printf("Неудачная попытка входа: %s \n", account.Email)
		return
	}
}

var GetInfo = func (w http.ResponseWriter, r* http.Request){
	user := r.Context().Value("user").(uint)
	account := database.GetUser(user)
	response := api.MakeResponse(true, "OK")
	response["account"] = account
	api.Respond(w,response)
	return
}