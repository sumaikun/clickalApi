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

//-------------------------------------- Patient Review functions ----------------------------------

func allPatientReviewEndPoint(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "application/json")

	patientReviews, err := dao.FindAllWithUsers("patientReviews")
	if err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, patientReviews)
}

func findPatientReviewByPatientEndPoint(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	w.Header().Set("Content-type", "application/json")

	patientReviews, err := dao.FindManyByKey("patientReviews", "patient", params["patient"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, patientReviews)

}

func createPatientReviewEndPoint(w http.ResponseWriter, r *http.Request) {

	//fmt.Print("here go the creation of patient review")

	user := context.Get(r, "user")

	userParsed := user.(bson.M)

	defer r.Body.Close()

	w.Header().Set("Content-type", "application/json")

	err, patientReview := patientReviewValidator(r)

	if len(err["validationError"].(url.Values)) > 0 {
		//fmt.Println(len(e))
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	patientReview.ID = bson.NewObjectId()
	patientReview.Date = time.Now().String()
	patientReview.UpdateDate = time.Now().String()
	patientReview.CreatedBy = userParsed["_id"].(bson.ObjectId).Hex()
	patientReview.UpdatedBy = userParsed["_id"].(bson.ObjectId).Hex()

	if err := dao.Insert("patientReviews", patientReview, nil); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusCreated, patientReview)

}

func findPatientReviewEndPoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	pet, err := dao.FindByID("patientReview", params["id"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Patient Review ID")
		return
	}
	Helpers.RespondWithJSON(w, http.StatusOK, pet)

}

func removePatientReviewEndPoint(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	err := dao.DeleteByID("patientReview", params["id"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Patient Review ID")
		return
	}
	Helpers.RespondWithJSON(w, http.StatusOK, nil)

}

func updatePatientReviewEndPoint(w http.ResponseWriter, r *http.Request) {

	user := context.Get(r, "user")

	userParsed := user.(bson.M)

	defer r.Body.Close()
	params := mux.Vars(r)

	w.Header().Set("Content-type", "application/json")

	err, patientReview := patientReviewValidator(r)

	if len(err["validationError"].(url.Values)) > 0 {
		//fmt.Println(len(e))
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	prevData, err2 := dao.FindByID("patientReviews", params["id"])
	if err2 != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Patient Review ID")
		return
	}

	parsedData := prevData.(bson.M)

	patientReview.ID = parsedData["_id"].(bson.ObjectId)

	patientReview.Date = parsedData["date"].(string)

	patientReview.CreatedBy = parsedData["createdBy"].(string)

	patientReview.UpdateDate = time.Now().String()

	patientReview.UpdatedBy = userParsed["_id"].(bson.ObjectId).Hex()

	if err := dao.Update("patientReviews", patientReview.ID, patientReview); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})

}
