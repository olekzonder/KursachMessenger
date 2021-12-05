package main
import (
 "KursachMessenger/restapi"
 "fmt"
 "github.com/gorilla/mux"
 "net/http"
 "os"
)
func main() {
 router := mux.NewRouter()
 router.HandleFunc("/api/user/create", restapi.CreateAccount).Methods("POST")
 router.HandleFunc("/api/user/login", restapi.Login).Methods("POST")
 router.HandleFunc("/api/user", restapi.GetInfo).Methods("GET")
 router.HandleFunc("/api/messages/send", restapi.SendMessage).Methods("POST")
 router.HandleFunc("/api/messages/get", restapi.GetAllMessages).Methods("GET")
 router.HandleFunc("/api/messages/get/unread", restapi.GetUnread).Methods("GET")
 router.HandleFunc("/", restapi.IsAlive).Methods("GET")
 router.Use(restapi.Auth)
 port := os.Getenv("port")
 if port == "" {
  port = "8080"
 }
 fmt.Println("Используемый порт: ", port)

  fmt.Println("Готов к работе.")
  err := http.ListenAndServe(":"+port,router)
  if err != nil {
   fmt.Println(err)
  }
 }