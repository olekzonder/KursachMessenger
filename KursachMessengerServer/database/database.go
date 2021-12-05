package database

import (
	"KursachMessenger/api"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"strings"
)

type Token struct{
	UserId uint
	jwt.StandardClaims //JWT - создание токена доступа, де-факто стандарт для аутентификации пользователей
}

type Account struct{
	gorm.Model
	Email    string `json:"email"` //Регистрация по email, чтобы можно было потом добавить, например, поддержку Gravatar для аватарок
	Password string `json:"password"`
	Token    string `json:"token"`
	LastRead int `json:"last_read"`
}

type Message struct { //Структура с сообщениями
	gorm.Model
	ID      int    `json:"ID"`   //ID сообщения
	Name    string `json:"name"` //От кого
	Message string `json:"msg"`  //Текст сообщения
}

//ФУНКЦИИ, СВЯЗАННЫЕ С АККАУНТОМ

func (account *Account) Validate() (bool,map[string]interface{}){ //Подтвердить правильность данных, введённых пользователем
	if !strings.Contains(account.Email,"@"){
		return false,api.MakeResponse(false,"NOEMAIL")
	}
	if len(account.Password)<6{
		return false,api.MakeResponse(false,"LEN")
	}
	temp := GetDB().Where("Email = ?", account.Email).First(account)
	err := temp.Error
	if err != nil && err != gorm.ErrRecordNotFound{
		return false, api.MakeResponse(false, "ERR01")
	}
	if err == gorm.ErrRecordNotFound {
		return true, api.MakeResponse(true, "OK")
	} else {
		return false, api.MakeResponse(false, "ALREADY_USED")
	}
}


func (account *Account) Create() (map[string]interface{},bool) {
	if ok, response := account.Validate();!ok{
		return response,false
	}
	encryptpass,_ := bcrypt.GenerateFromPassword([]byte(account.Password),bcrypt.DefaultCost)
	account.Password = string(encryptpass)
	GetDB().Create(account)
	if account.ID <= 0 {
		return api.MakeResponse(false,"ERR02"),false
	}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"),&Token{UserId:account.ID})
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString
	account.LastRead = 0
	response := api.MakeResponse(true,"OK")
	response["account"] = account
	return response,true
}

func Login(email, password string) (map[string]interface{},bool) {
	account := &Account{}
	err := GetDB().Where("email =?", email).First(account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound{
			return api.MakeResponse(false,"USER_NOT_FOUND"),false
		}
		return api.MakeResponse(false,"ERR 03"),false
	}
	err = bcrypt.CompareHashAndPassword([]byte(account.Password),[]byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword{
		return api.MakeResponse(false,"INCORRECT_PASS"),false
	}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), &Token{UserId:account.ID})
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString
	response := api.MakeResponse(true,"OK")
	response["account"] = account
	return response,true
}

func GetUser(id uint) *Account {
	var acc *Account
	GetDB().Table("accounts").Where("id =?", id).Scan(&acc)
	if acc.Email == ""{
		return nil
	}
	return acc
}


//ФУНКЦИИ, СВЯЗАННЫЕ С СООБЩЕНИЯМИ

func (message *Message) Send() (map[string]interface{},bool) {
	if len(strings.TrimSpace(message.Message)) == 0{
		return api.MakeResponse(false,"EMPTY_MSG"),false
	}
	GetDB().Create(message)
	if message.ID <= 0 {
		return api.MakeResponse(false,"ERR04"),false
	}
	response := api.MakeResponse(true,"OK")
	response["message"] = message
	return response,true
}

func GetLastMessageID() (id int){
	var msg *Message
	GetDB().Last(&msg)
	id = msg.ID
	return id
}

func GetMessages (offset,limit int) (*[]Message,error){
	var err error
	list := make([]Message,0)
	query := GetDB().Table("messages").Limit(limit).Offset(offset).Order("created_at").Find(&list)
	err = query.Error
	if err == nil {
		return &list, nil
	} else {
		return &list,err
	}
}