package main

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"

	Helpers "github.com/sumaikun/clickal-rest-api/helpers"
)

//-------------------------------- PatientFiles functions ----------------------------------

func allPatientFilesEndPoint(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "application/json")

	patientFiles, err := dao.FindAllWithPatients("patientFiles")
	if err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, patientFiles)
}

func findPatientFilesByPatientEndPoint(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	w.Header().Set("Content-type", "application/json")

	//fmt.Println("patient log" + params["patient"])

	patientFiles, err := dao.FindManyByKey("patientFiles", "patient", params["patient"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, patientFiles)

}

func createPatientFilesEndPoint(w http.ResponseWriter, r *http.Request) {

	user := context.Get(r, "user")

	userParsed := user.(bson.M)

	defer r.Body.Close()
	w.Header().Set("Content-type", "application/json")

	err, patientsFiles := patientsFilesValidator(r)

	if len(err["validationError"].(url.Values)) > 0 {
		//fmt.Println(len(e))
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	patientsFiles.ID = bson.NewObjectId()
	patientsFiles.Date = time.Now().String()
	patientsFiles.UpdateDate = time.Now().String()
	patientsFiles.CreatedBy = userParsed["_id"].(bson.ObjectId).Hex()
	patientsFiles.UpdatedBy = userParsed["_id"].(bson.ObjectId).Hex()

	if err := dao.Insert("patientFiles", patientsFiles, nil); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusCreated, patientsFiles)

}

func findPatientFilesEndPoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	pet, err := dao.FindByID("patientFiles", params["id"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid PatientsFile ID")
		return
	}
	Helpers.RespondWithJSON(w, http.StatusOK, pet)

}

func removePatientFilesEndPoint(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	err := dao.DeleteByID("patientFiles", params["id"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid PatientsFile ID")
		return
	}
	Helpers.RespondWithJSON(w, http.StatusOK, nil)

}

func updatePatientFilesEndPoint(w http.ResponseWriter, r *http.Request) {

	fmt.Println("update log")

	user := context.Get(r, "user")

	userParsed := user.(bson.M)

	defer r.Body.Close()
	params := mux.Vars(r)

	w.Header().Set("Content-type", "application/json")

	err, patientsFiles := patientsFilesValidator(r)

	if len(err["validationError"].(url.Values)) > 0 {
		//fmt.Println(len(e))
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	prevData, err2 := dao.FindByID("patientFiles", params["id"])
	if err2 != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Patients File ID")
		return
	}

	parsedData := prevData.(bson.M)

	patientsFiles.ID = parsedData["_id"].(bson.ObjectId)

	patientsFiles.Date = parsedData["date"].(string)

	patientsFiles.UpdateDate = time.Now().String()

	patientsFiles.CreatedBy = parsedData["createdBy"].(string)

	patientsFiles.UpdatedBy = userParsed["_id"].(bson.ObjectId).Hex()

	if err := dao.Update("patientFiles", patientsFiles.ID, patientsFiles); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})

}
