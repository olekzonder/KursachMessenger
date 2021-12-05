package api

import (
	"encoding/json"
	"net/http"
)

func MakeResponse(status bool, message string) map[string]interface{} { //Собрать сообщение в JSON: информацию об аккаунте, все сообщения и т.д
	return map[string]interface{}{"status": status, "message": message}
}

func Respond(w http.ResponseWriter,data map[string]interface{}){ //Отсылка JSON клиенту
	w.Header().Add("Content-Type","application/json")
	json.NewEncoder(w).Encode(data)
}