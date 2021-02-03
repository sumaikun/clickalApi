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

//--------------------------------physiological Constants functions ----------------------------------

func allPhysiologicalConstantsEndPoint(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "application/json")

	physiologicalConstant, err := dao.FindAllWithUsers("physiologicalConstants")
	if err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, physiologicalConstant)
}

func findPhysiologicalConstantsByPatientEndPoint(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	w.Header().Set("Content-type", "application/json")

	physiologicalConstant, err := dao.FindManyByKey("physiologicalConstants", "patient", params["patient"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, physiologicalConstant)

}

func createPhysiologicalConstantsEndPoint(w http.ResponseWriter, r *http.Request) {

	user := context.Get(r, "user")

	userParsed := user.(bson.M)

	defer r.Body.Close()
	w.Header().Set("Content-type", "application/json")

	err, physiologicalConstant := physiologicalConstantsValidator(r)

	if len(err["validationError"].(url.Values)) > 0 {
		//fmt.Println(len(e))
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	physiologicalConstant.ID = bson.NewObjectId()
	physiologicalConstant.Date = time.Now().String()
	physiologicalConstant.UpdateDate = time.Now().String()
	physiologicalConstant.CreatedBy = userParsed["_id"].(bson.ObjectId).Hex()
	physiologicalConstant.UpdatedBy = userParsed["_id"].(bson.ObjectId).Hex()

	if err := dao.Insert("physiologicalConstants", physiologicalConstant, nil); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusCreated, physiologicalConstant)

}

func findPhysiologicalConstantsEndPoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	pet, err := dao.FindByID("physiologicalConstants", params["id"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Constant ID")
		return
	}
	Helpers.RespondWithJSON(w, http.StatusOK, pet)

}

func removePhysiologicalConstantsEndPoint(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	err := dao.DeleteByID("physiologicalConstants", params["id"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Constant ID")
		return
	}
	Helpers.RespondWithJSON(w, http.StatusOK, nil)

}

func updatePhysiologicalConstantsEndPoint(w http.ResponseWriter, r *http.Request) {

	user := context.Get(r, "user")

	userParsed := user.(bson.M)

	defer r.Body.Close()
	params := mux.Vars(r)

	w.Header().Set("Content-type", "application/json")

	err, physiologicalConstant := physiologicalConstantsValidator(r)

	if len(err["validationError"].(url.Values)) > 0 {
		//fmt.Println(len(e))
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	prevData, err2 := dao.FindByID("physiologicalConstants", params["id"])
	if err2 != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Constant ID")
		return
	}

	parsedData := prevData.(bson.M)

	physiologicalConstant.ID = parsedData["_id"].(bson.ObjectId)

	physiologicalConstant.Date = parsedData["date"].(string)

	physiologicalConstant.UpdateDate = time.Now().String()

	physiologicalConstant.CreatedBy = parsedData["createdBy"].(string)

	physiologicalConstant.UpdatedBy = userParsed["_id"].(bson.ObjectId).Hex()

	if err := dao.Update("physiologicalConstants", physiologicalConstant.ID, physiologicalConstant); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})

}
