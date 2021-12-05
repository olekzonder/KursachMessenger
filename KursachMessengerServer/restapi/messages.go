package restapi

import (
	"KursachMessenger/api"
	"KursachMessenger/database"
	"encoding/json"
	"fmt"
	"net/http"
)

var SendMessage = func(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(uint)
	message := &database.Message{}
	account := database.GetUser(user)
	err := json.NewDecoder(r.Body).Decode(message)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		api.Respond(w, api.MakeResponse(false, "ERR11"))
		return
	}
	message.Name = account.Email
	response,ok := message.Send()
	if ok == true{
		fmt.Printf("%s : %s \n",account.Email,message.Message)
		api.Respond(w, response)
	} else{
		w.WriteHeader(http.StatusBadRequest)
		api.Respond(w, response)
	}
}

var GetAllMessages = func(w http.ResponseWriter, r *http.Request){
	user := r.Context().Value("user").(uint)
	account := database.GetUser(user)
	limit := database.GetLastMessageID()
	res, err := database.GetMessages(0,limit)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		api.Respond(w,api.MakeResponse(false,"ERR12"))
		return
	}
	response := api.MakeResponse(true,fmt.Sprintf("%d",limit))
	for i,item := range *res{
		response[fmt.Sprintf("%d", i)] = item
	}
	account.LastRead = limit
	database.GetDB().Save(&account)
	api.Respond(w, response)
}

var  GetUnread = func(w http.ResponseWriter, r *http.Request){
	user := r.Context().Value("user").(uint)
	account := database.GetUser(user)
	offset := account.LastRead
	limit := database.GetLastMessageID()
	if limit==offset {
		api.Respond(w,api.MakeResponse(false,"NO_NEW_MSG"))
		return
	}
	res, err := database.GetMessages(offset,limit)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		api.Respond(w,api.MakeResponse(false,"ERR13"))
		return
	}
	response := api.MakeResponse(true,fmt.Sprintf("%d",limit-offset))
	for i,item := range *res{
		response[fmt.Sprintf("%d", i)] = item
	}
	account.LastRead = limit
	database.GetDB().Save(&account)
	api.Respond(w, response)
}