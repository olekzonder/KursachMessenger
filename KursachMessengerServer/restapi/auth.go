package restapi

import (
	"KursachMessenger/api"
	"KursachMessenger/database"
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"strings"
)

var Auth = func(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
		authNotRequired := []string{"/api/user/create","/api/user/login","/"} //Обработка тех запросов, где не требуется аутентификация: логин и регистрация
		curRequest := r.URL.Path
		for _, value := range authNotRequired{
			if value == curRequest {
				next.ServeHTTP(w,r)
				return
			}
		}

		response := make(map[string]interface{})
		tokenHeader := r.Header.Get("Authorization")
		if tokenHeader == "" {
			response = api.MakeResponse(false, "NO_TOKEN")
			w.WriteHeader(http.StatusForbidden)
			api.Respond(w,response)
			return
		}
		receivedToken := strings.Split(tokenHeader," ") //Проверка на битый токен, формат "Bearer <токен>"
		if len(receivedToken) != 2 {
			response = api.MakeResponse(false, "TOKEN_CORRUPT")
			w.WriteHeader(http.StatusForbidden)
			api.Respond(w, response)
			return
		}
		tokenpart := receivedToken[1]
		tokendb := &database.Token{}
		token, err := jwt.ParseWithClaims(tokenpart,tokendb,func(token *jwt.Token) (interface{},error){
			return []byte(os.Getenv("token_password")),nil
		})
		if err != nil {
			response = api.MakeResponse(false,"TOKEN_CORRUPT")
			w.WriteHeader(http.StatusForbidden)
			api.Respond(w,response)
			return
		}
		if !token.Valid{
			response = api.MakeResponse(false,"TOKEN_NOT_VALID")
			w.WriteHeader(http.StatusForbidden)
			api.Respond(w,response)
			return
		}
		fmt.Printf("REQUEST:%s (ID:%d) \n", database.GetUser(tokendb.UserId).Email,tokendb.UserId)
		ctx := context.WithValue(r.Context(),"user",tokendb.UserId)
		r = r.WithContext(ctx)
		next.ServeHTTP(w,r)
	})
}