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

//-------------------------------------- MedicalCenters Functions ----------------------------------

func allMedicalCentersEndPoint(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "application/json")

	medicalCenters, err := dao.FindAll("medicalCenters")
	if err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, medicalCenters)
}

func createMedicalCenterEndPoint(w http.ResponseWriter, r *http.Request) {

	user := context.Get(r, "user")

	userParsed := user.(bson.M)

	defer r.Body.Close()
	w.Header().Set("Content-type", "application/json")

	err, medicalCenter := medicalCenterValidator(r)

	if len(err["validationError"].(url.Values)) > 0 {
		//fmt.Println(len(e))
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	medicalCenter.ID = bson.NewObjectId()
	medicalCenter.Date = time.Now().String()
	medicalCenter.UpdateDate = time.Now().String()
	medicalCenter.CreatedBy = userParsed["_id"].(bson.ObjectId).Hex()
	medicalCenter.UpdatedBy = userParsed["_id"].(bson.ObjectId).Hex()

	if err := dao.Insert("medicalCenters", medicalCenter, []string{"name"}); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusCreated, medicalCenter)

}

func findMedicalCenterEndPoint(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	medicalCenter, err := dao.FindByID("medicalCenters", params["id"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid MedicalCenter ID")
		return
	}
	Helpers.RespondWithJSON(w, http.StatusOK, medicalCenter)

}

func removeMedicalCenterEndPoint(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	err := dao.DeleteByID("medicalCenters", params["id"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid MedicalCenter ID")
		return
	}
	Helpers.RespondWithJSON(w, http.StatusOK, nil)

}

func updateMedicalCenterEndPoint(w http.ResponseWriter, r *http.Request) {

	user := context.Get(r, "user")

	userParsed := user.(bson.M)

	defer r.Body.Close()
	params := mux.Vars(r)

	w.Header().Set("Content-type", "application/json")

	err, medicalCenter := medicalCenterValidator(r)

	if len(err["validationError"].(url.Values)) > 0 {
		//fmt.Println(len(e))
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	prevData, err2 := dao.FindByID("medicalCenters", params["id"])
	if err2 != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid MedicalCenter ID")
		return
	}

	parsedData := prevData.(bson.M)

	medicalCenter.ID = parsedData["_id"].(bson.ObjectId)

	medicalCenter.Date = parsedData["date"].(string)

	medicalCenter.UpdateDate = time.Now().String()

	if parsedData["createdBy"] == nil {
		medicalCenter.CreatedBy = userParsed["_id"].(bson.ObjectId).Hex()
	} else {
		medicalCenter.CreatedBy = parsedData["createdBy"].(string)
	}

	medicalCenter.UpdatedBy = userParsed["_id"].(bson.ObjectId).Hex()

	if err := dao.Update("medicalCenters", medicalCenter.ID, medicalCenter); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})

}
