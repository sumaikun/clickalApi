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

//-------------------------------- Medicines functions ----------------------------------

func findMedicinesByPatientEndPoint(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	w.Header().Set("Content-type", "application/json")

	medicines, err := dao.FindManyByKey("medicines", "patient", params["patient"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, medicines)

}

func findMedicinesByAppointmentEndPoint(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	w.Header().Set("Content-type", "application/json")

	medicines, err := dao.FindManyByKey("medicines", "appointment", params["appointment"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, medicines)

}

func createMedicinesEndPoint(w http.ResponseWriter, r *http.Request) {

	user := context.Get(r, "user")

	userParsed := user.(bson.M)

	defer r.Body.Close()
	w.Header().Set("Content-type", "application/json")

	err, medicine := medicinesValidator(r)

	if len(err["validationError"].(url.Values)) > 0 {
		//fmt.Println(len(e))
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	medicine.ID = bson.NewObjectId()
	medicine.Date = time.Now().String()
	medicine.UpdateDate = time.Now().String()
	medicine.CreatedBy = userParsed["_id"].(bson.ObjectId).Hex()
	medicine.UpdatedBy = userParsed["_id"].(bson.ObjectId).Hex()

	if err := dao.Insert("medicines", medicine, nil); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusCreated, medicine)

}

func findMedicinesEndPoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	medicine, err := dao.FindByID("medicines", params["id"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Medicine ID")
		return
	}
	Helpers.RespondWithJSON(w, http.StatusOK, medicine)

}

func removeMedicinesEndPoint(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	err := dao.DeleteByID("medicines", params["id"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Medicine ID")
		return
	}
	Helpers.RespondWithJSON(w, http.StatusOK, nil)

}

func updateMedicinesEndPoint(w http.ResponseWriter, r *http.Request) {

	user := context.Get(r, "user")

	userParsed := user.(bson.M)

	defer r.Body.Close()
	params := mux.Vars(r)

	w.Header().Set("Content-type", "application/json")

	err, medicine := medicinesValidator(r)

	if len(err["validationError"].(url.Values)) > 0 {
		//fmt.Println(len(e))
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	prevData, err2 := dao.FindByID("medicines", params["id"])
	if err2 != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Appointment ID")
		return
	}

	parsedData := prevData.(bson.M)

	medicine.ID = parsedData["_id"].(bson.ObjectId)

	medicine.Date = parsedData["date"].(string)

	medicine.UpdateDate = time.Now().String()

	medicine.CreatedBy = parsedData["createdBy"].(string)

	medicine.UpdatedBy = userParsed["_id"].(bson.ObjectId).Hex()

	if err := dao.Update("medicines", medicine.ID, medicine); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})

}
