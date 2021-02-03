package main

import (
	"net/http"
	"net/url"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"

	Models "github.com/sumaikun/clickal-rest-api/models"

	Helpers "github.com/sumaikun/clickal-rest-api/helpers"
)

//-----------------------------  Doctors functions --------------------------------------------------

func allDoctorsEndPoint(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "application/json")

	doctors, err := dao.FindAllWithCities("doctors")
	if err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, doctors)
}

func createDoctorsEndPoint(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	user := context.Get(r, "user")

	userParsed := user.(bson.M)

	w.Header().Set("Content-type", "application/json")

	err, doctor := doctorValidator(r)

	if len(err["validationError"].(url.Values)) > 0 {
		//fmt.Println(len(e))
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	doctor.ID = bson.NewObjectId()
	doctor.Date = time.Now().String()
	doctor.UpdateDate = time.Now().String()
	doctor.CreatedBy = userParsed["_id"].(bson.ObjectId).Hex()
	doctor.UpdatedBy = userParsed["_id"].(bson.ObjectId).Hex()
	doctor.State = "CHANGE_PASSWORD"

	if err := dao.Insert("doctors", doctor, []string{"email"}); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Models.TypeClaims{
		Username: doctor.Email,
		Type:     "doctor",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, _ := token.SignedString(jwtKey)

	go sendResetPasswordEmail(tokenString, doctor.Email)

	go Helpers.RespondWithJSON(w, http.StatusCreated, doctor)

}

func findDoctorEndPoint(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	doctor, err := dao.FindByID("doctors", params["id"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Doctor ID")
		return
	}
	Helpers.RespondWithJSON(w, http.StatusOK, doctor)

}

func removeDoctorEndPoint(w http.ResponseWriter, r *http.Request) {

	/*params := mux.Vars(r)
	err := dao.DeleteByID("doctors", params["id"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Doctor ID")
		return
	}
	Helpers.RespondWithJSON(w, http.StatusOK, nil)*/

	params := mux.Vars(r)

	_, err2 := dao.FindByID("doctors", params["id"])
	if err2 != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Doctor ID")
		return
	}

	err := dao.PartialUpdate("doctors", params["id"], bson.M{"state": "INACTIVE"})
	if err != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	Helpers.RespondWithJSON(w, http.StatusOK, nil)
}

func updateDoctorEndPoint(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	params := mux.Vars(r)

	user := context.Get(r, "user")

	userParsed := user.(bson.M)

	w.Header().Set("Content-type", "application/json")

	err, doctor := doctorValidator(r)

	if len(err["validationError"].(url.Values)) > 0 {
		//fmt.Println(len(e))
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	prevDoctor, err2 := dao.FindByID("doctors", params["id"])
	if err2 != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Doctor ID")
		return
	}

	parsedData := prevDoctor.(bson.M)

	doctor.ID = parsedData["_id"].(bson.ObjectId)

	doctor.State = parsedData["state"].(string)

	doctor.Date = parsedData["date"].(string)

	doctor.UpdateDate = time.Now().String()

	if parsedData["createdBy"] == nil {
		doctor.CreatedBy = userParsed["_id"].(bson.ObjectId).Hex()
	} else {
		doctor.CreatedBy = parsedData["createdBy"].(string)
	}

	doctor.UpdatedBy = userParsed["_id"].(bson.ObjectId).Hex()

	if len(doctor.Password) == 0 {
		doctor.Password = parsedData["password"].(string)
	} else {
		doctor.Password, _ = Helpers.HashPassword(doctor.Password)
	}

	if err := dao.Update("doctors", doctor.ID, doctor); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})

}
