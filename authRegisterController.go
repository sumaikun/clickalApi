package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/dgrijalva/jwt-go"
	"gopkg.in/mgo.v2/bson"

	Models "github.com/sumaikun/clickal-rest-api/models"

	Helpers "github.com/sumaikun/clickal-rest-api/helpers"

	C "github.com/sumaikun/clickal-rest-api/config"
)

//-----------------------------  Auth functions --------------------------------------------------

func authentication(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	var userType int

	response := &Models.TokenResponse{Token: "", User: nil, UserType: 0}

	var creds Models.Credentials

	copyBody := r.Body

	// Get the JSON body and decode into credentials
	err := json.NewDecoder(copyBody).Decode(&creds)

	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the expected password from our in memory map
	expectedPassword, ok := Models.Users[creds.Username]

	// If a password exists for the given user
	// AND, if it is the same as the password we received, the we can move ahead
	// if NOT, then we return an "Unauthorized" status
	if !ok || !Helpers.CheckPasswordHash(creds.Password, expectedPassword) {

		fmt.Println("in condition")

		user, err := dao.FindOneByKEY("users", "email", creds.Username)

		//fmt.Println("user")
		//fmt.Println(user)

		if user == nil {

			fmt.Println("user not found trying doctor")

			user, err = dao.FindOneByKEY("doctors", "email", creds.Username)

			//fmt.Println("user", user)

			if user == nil {

				fmt.Println("user not found trying patient")

				user, err = dao.FindOneByKEY("patients", "email", creds.Username)
				if err != nil {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				if user != nil {
					userType = 3
				}

			} else {
				userType = 2
			}

			if err != nil {

				fmt.Println("err", err)

				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			//fmt.Println(user)
		} else {
			userType = 1
		}

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		fmt.Println(user)

		match := Helpers.CheckPasswordHash(creds.Password, user.(bson.M)["password"].(string))

		if !match {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if user.(bson.M)["state"].(string) != "ACTIVE" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		response.User = user.(bson.M)

	}

	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	expirationTime := time.Now().Add(8 * time.Hour)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Models.Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Finally, we set the client cookie for "token" as the JWT we just generated
	// we also set an expiry time which is the same as the token itself
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	w.Header().Set("Content-type", "application/json")

	//Generate json response for get the token
	response.Token = tokenString

	response.UserType = userType

	json.NewEncoder(w).Encode(response)
}

func exampleHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	io.WriteString(w, `{"status":"ok"}`)
}

func createInititalUser(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "application/json")

	users, err := dao.FindAll("users")
	if err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if users == nil {

		fmt.Println("is nil")

		var user Models.User

		err := json.NewDecoder(r.Body).Decode(&user)

		if err != nil {
			// If the structure of the body is wrong, return an HTTP error
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		user.ID = bson.NewObjectId()
		user.Date = time.Now().String()
		user.UpdateDate = time.Now().String()

		if len(user.Password) != 0 {
			user.Password, _ = Helpers.HashPassword(user.Password)
		}

		if err := dao.Insert("users", user, []string{"email"}); err != nil {
			Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		Helpers.RespondWithJSON(w, http.StatusCreated, user)

	} else {
		Helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "can not create initial users again"})
	}

}

func resetPassword(w http.ResponseWriter, r *http.Request) {

	var config = C.Config{}
	config.Read()

	w.Header().Set("Content-type", "application/json")

	err, reset := resetPasswordValidator(r)

	if len(err["validationError"].(url.Values)) > 0 {
		//fmt.Println(len(e))
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	claims := jwt.MapClaims{}
	_, err2 := jwt.ParseWithClaims(reset.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Jwtkey), nil
	})

	if err2 != nil {
		Helpers.RespondWithJSON(w, http.StatusForbidden, map[string]string{"result": "Error decoding jwt"})
		//log.Fatal("Error decoding jwt")
	}

	//fmt.Println(claims)

	//claims["username"].(string)

	if claims["type"].(string) == "forgot-password" {

		user, _ := dao.FindOneByKEY("users", "email", claims["username"].(string))

		doctor, _ := dao.FindOneByKEY("doctors", "email", claims["username"].(string))

		patient, _ := dao.FindOneByKEY("patients", "email", claims["username"].(string))

		if user != nil {
			parsedUser := user.(bson.M)
			parsedUser["password"], _ = Helpers.HashPassword(reset.Password)
			parsedUser["state"] = "ACTIVE"
			if err := dao.Update("users", parsedUser["_id"].(bson.ObjectId), parsedUser); err != nil {
				Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
				return
			}
		}

		if doctor != nil {
			parsedUser := doctor.(bson.M)
			if err := dao.Update("doctors", parsedUser["_id"].(bson.ObjectId), parsedUser); err != nil {
				Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
				return
			}
		}

		if patient != nil {
			parsedUser := patient.(bson.M)
			if err := dao.Update("patients", parsedUser["_id"].(bson.ObjectId), parsedUser); err != nil {
				Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
				return
			}
		}

		Helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "password reseted"})

	} else {
		user, err3 := dao.FindOneByKEY(claims["type"].(string)+"s", "email", claims["username"].(string))

		if err3 != nil {
			Helpers.RespondWithError(w, http.StatusInternalServerError, err3.Error())
			return
		}

		//fmt.Println(user)

		parsedUser := user.(bson.M)

		//fmt.Println(parsedUser["state"])

		if parsedUser["state"].(string) == "CHANGE_PASSWORD" {
			parsedUser["password"], _ = Helpers.HashPassword(reset.Password)
			parsedUser["state"] = "ACTIVE"
			//fmt.Println(parsedUser)
			if err := dao.Update(claims["type"].(string)+"s", parsedUser["_id"].(bson.ObjectId), parsedUser); err != nil {
				Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
				return
			}

		} else {
			Helpers.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"result": "can't change password of this account"})
			return
		}

		Helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "password reseted"})
	}

}

func forgotPassword(w http.ResponseWriter, r *http.Request) {

	var config = C.Config{}
	config.Read()

	w.Header().Set("Content-type", "application/json")

	err, reset := forgotPasswordValidator(r)

	if len(err["validationError"].(url.Values)) > 0 {
		//fmt.Println(len(e))
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Models.TypeClaims{
		Username: reset.Email,
		Type:     "forgot-password",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, _ := token.SignedString(jwtKey)

	go sendForgotPasswordEmail(tokenString, reset.Email)
	go Helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func confirmAccount(w http.ResponseWriter, r *http.Request) {

	var config = C.Config{}
	config.Read()

	w.Header().Set("Content-type", "application/json")

	err, reset := confirmAccountValidator(r)

	if len(err["validationError"].(url.Values)) > 0 {
		//fmt.Println(len(e))
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	claims := jwt.MapClaims{}
	_, err2 := jwt.ParseWithClaims(reset.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Jwtkey), nil
	})

	if err2 != nil {
		Helpers.RespondWithJSON(w, http.StatusForbidden, map[string]string{"result": "Error decoding jwt"})
		//log.Fatal("Error decoding jwt")
	}

	fmt.Println(claims)

	//claims["username"].(string)

	user, err3 := dao.FindOneByKEY(claims["type"].(string)+"s", "email", claims["username"].(string))

	if err3 != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err3.Error())
		return
	}

	//fmt.Println(user)

	parsedUser := user.(bson.M)

	fmt.Println("parsedUser", parsedUser)

	parsedUser["state"] = "ACTIVE"

	//fmt.Println(parsedUser)
	if err := dao.PartialUpdate(claims["type"].(string)+"s", parsedUser["_id"].(bson.ObjectId).Hex(), bson.M{"state": "ACTIVE"}); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "Account confirmed"})

}

func registerDoctor(w http.ResponseWriter, r *http.Request) {

	var config = C.Config{}
	config.Read()

	w.Header().Set("Content-type", "application/json")

	err, user := userRegisterValidator(r)

	if len(err["validationError"].(url.Values)) > 0 {
		//fmt.Println(len(e))
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	user.ID = bson.NewObjectId()
	user.Date = time.Now().String()
	user.UpdateDate = time.Now().String()
	user.State = "INACTIVE"
	user.Password, _ = Helpers.HashPassword(user.Password)

	if err := dao.Insert("doctors", user, []string{"email"}); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Models.TypeClaims{
		Username: user.Email,
		Type:     "doctor",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, _ := token.SignedString(jwtKey)

	go sendConfirmationEmail(tokenString, user.Email)

	go Helpers.RespondWithJSON(w, http.StatusCreated, user)

}

func registerPatient(w http.ResponseWriter, r *http.Request) {

	var config = C.Config{}
	config.Read()

	w.Header().Set("Content-type", "application/json")

	err, user := userRegisterValidator(r)

	if len(err["validationError"].(url.Values)) > 0 {
		//fmt.Println(len(e))
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	user.ID = bson.NewObjectId()
	user.Date = time.Now().String()
	user.UpdateDate = time.Now().String()
	user.State = "INACTIVE"
	user.Password, _ = Helpers.HashPassword(user.Password)

	if err := dao.Insert("patients", user, []string{"email"}); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Models.TypeClaims{
		Username: user.Email,
		Type:     "patient",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, _ := token.SignedString(jwtKey)

	go sendConfirmationEmail(tokenString, user.Email)

	go Helpers.RespondWithJSON(w, http.StatusCreated, user)
}
