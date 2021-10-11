package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/egaevan/merchant/model"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	DB *sql.DB
}

type ErrorResponse struct {
	Err string
}

func NewUserRepository() UserRepository {
	return &User{}
}

func (u *User) Login(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	user := &model.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		var resp = map[string]interface{}{"status": false, "message": "Invalid request"}
		json.NewEncoder(w).Encode(resp)
		return
	}
	resp := u.FindOne(ctx, user.Email, user.Password)
	json.NewEncoder(w).Encode(resp)
}

func (u *User) FindOne(ctx context.Context, email, password string) map[string]interface{} {
	user := &model.User{}

	query := `
			SELECT 
				ID,
				name,
				email,
				gender
			FROM 
				user
			WHERE
				email = ? AND flag_aktif = 1`

	err := u.DB.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Name, &user.Email, &user.Gender)
	if err != nil {
		fmt.Println(err)
	}

	if err != nil {
		var resp = map[string]interface{}{"status": false, "message": "Email address not found"}
		return resp
	}

	expiresAt := time.Now().Add(time.Minute * 100000).Unix()

	errf := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if errf != nil && errf == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		var resp = map[string]interface{}{"status": false, "message": "Invalid login credentials. Please try again"}
		return resp
	}

	tk := &model.Token{
		UserID: user.ID,
		Name:   user.Name,
		Email:  user.Email,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)

	tokenString, error := token.SignedString([]byte("secret"))
	if error != nil {
		fmt.Println(error)
	}

	var resp = map[string]interface{}{"status": false, "message": "logged in"}
	resp["token"] = tokenString //Store the token in the response
	resp["user"] = user
	return resp
}

// func FetchUsers(w http.ResponseWriter, r *http.Request) {
// 	var users []model.User
// 	db.Preload("auths").Find(&users)

// 	json.NewEncoder(w).Encode(users)
// }

// func CreateUser(w http.ResponseWriter, r *http.Request) {

// 	user := &model.User{}
// 	json.NewDecoder(r.Body).Decode(user)

// 	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
// 	if err != nil {
// 		fmt.Println(err)
// 		err := ErrorResponse{
// 			Err: "Password Encryption  failed",
// 		}
// 		json.NewEncoder(w).Encode(err)
// 	}

// 	user.Password = string(pass)

// 	createdUser := db.Create(user)
// 	var errMessage = createdUser.Error

// 	if createdUser.Error != nil {
// 		fmt.Println(errMessage)
// 	}
// 	json.NewEncoder(w).Encode(createdUser)
// }

// func FetchUsers(w http.ResponseWriter, r *http.Request) {
// 	var users []model.User
// 	db.Preload("auths").Find(&users)

// 	json.NewEncoder(w).Encode(users)
// }

// func UpdateUser(w http.ResponseWriter, r *http.Request) {
// 	user := &model.User{}
// 	params := mux.Vars(r)
// 	var id = params["id"]
// 	db.First(&user, id)
// 	json.NewDecoder(r.Body).Decode(user)
// 	db.Save(&user)
// 	json.NewEncoder(w).Encode(&user)
// }

// func DeleteUser(w http.ResponseWriter, r *http.Request) {
// 	params := mux.Vars(r)
// 	var id = params["id"]
// 	var user model.User
// 	db.First(&user, id)
// 	db.Delete(&user)
// 	json.NewEncoder(w).Encode("User deleted")
// }

// func GetUser(w http.ResponseWriter, r *http.Request) {
// 	params := mux.Vars(r)
// 	var id = params["id"]
// 	var user model.User
// 	db.First(&user, id)
// 	json.NewEncoder(w).Encode(&user)
// }
