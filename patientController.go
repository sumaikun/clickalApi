package main

import (
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"

	Helpers "github.com/sumaikun/clickal-rest-api/helpers"
)

//-----------------------------  Patients functions --------------------------------------------------

func allPatientsEndPoint(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	userType := context.Get(r, "userType")

	//fmt.Println("userType", userType)

	w.Header().Set("Content-type", "application/json")

	if userType.(int) == 1 {
		patients, err := dao.FindAllWithCities("patients")
		if err != nil {
			Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		Helpers.RespondWithJSON(w, http.StatusOK, patients)
	}

	if userType.(int) == 2 {
		user := context.Get(r, "user")

		userParsed := user.(bson.M)

		//fmt.Println("userParsed", userParsed)

		patients, err := dao.FindInArrayKey("patients", "doctors", userParsed["_id"].(bson.ObjectId).Hex())
		if err != nil {
			Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		Helpers.RespondWithJSON(w, http.StatusOK, patients)

	}

}

func createPatientsEndPoint(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	userType := context.Get(r, "userType")

	user := context.Get(r, "user")

	userParsed := user.(bson.M)

	w.Header().Set("Content-type", "application/json")

	err, patient := patientValidator(r)

	if len(err["validationError"].(url.Values)) > 0 {
		//fmt.Println(len(e))
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	var doctorsArray []string

	patient.ID = bson.NewObjectId()
	patient.Date = time.Now().String()
	patient.UpdateDate = time.Now().String()
	patient.CreatedBy = userParsed["_id"].(bson.ObjectId).Hex()
	patient.UpdatedBy = userParsed["_id"].(bson.ObjectId).Hex()

	if userType.(int) == 2 {
		doctorsArray = append(doctorsArray, userParsed["_id"].(bson.ObjectId).Hex())

		patient.Doctors = doctorsArray
	}

	if len(patient.Password) != 0 {
		patient.Password, _ = Helpers.HashPassword(patient.Password)
	}

	if err := dao.Insert("patients", patient, []string{"email"}); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusCreated, patient)

}

func findPatientEndPoint(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	patient, err := dao.FindByID("patients", params["id"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Patient ID")
		return
	}
	Helpers.RespondWithJSON(w, http.StatusOK, patient)

}

func removePatientEndPoint(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	err := dao.DeleteByID("patients", params["id"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Patient ID")
		return
	}
	Helpers.RespondWithJSON(w, http.StatusOK, nil)

}

func updatePatientEndPoint(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	params := mux.Vars(r)

	usera := context.Get(r, "user")

	userParsed := usera.(bson.M)

	w.Header().Set("Content-type", "application/json")

	err, patient := patientValidator(r)

	if len(err["validationError"].(url.Values)) > 0 {
		//fmt.Println(len(e))
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	prevPatient, err2 := dao.FindByID("patients", params["id"])
	if err2 != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Patient ID")
		return
	}

	parsedData := prevPatient.(bson.M)

	patient.ID = parsedData["_id"].(bson.ObjectId)

	patient.State = parsedData["state"].(string)

	patient.Date = parsedData["date"].(string)

	patient.UpdateDate = time.Now().String()

	if parsedData["createdBy"] == nil {
		patient.CreatedBy = userParsed["_id"].(bson.ObjectId).Hex()
	} else {
		patient.CreatedBy = parsedData["createdBy"].(string)
	}

	patient.UpdatedBy = userParsed["_id"].(bson.ObjectId).Hex()

	if len(patient.Password) == 0 {
		patient.Password = parsedData["password"].(string)
	} else {
		patient.Password, _ = Helpers.HashPassword(patient.Password)
	}

	if err := dao.Update("patients", patient.ID, patient); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})

}
