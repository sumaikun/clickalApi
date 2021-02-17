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

func allAgendaAnnotationsEndPoint(w http.ResponseWriter, r *http.Request) {

	userType := context.Get(r, "userType")

	user := context.Get(r, "user")

	userParsed := user.(bson.M)

	w.Header().Set("Content-type", "application/json")

	if userType.(int) == 1 {
		agendaAnnotations, err := dao.FindAllWithPatients("agendaAnnotations")
		if err != nil {
			Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		Helpers.RespondWithJSON(w, http.StatusOK, agendaAnnotations)
	}

	if userType.(int) == 2 {
		agendaAnnotations, err := dao.FindManyByKeyWithPatiens("agendaAnnotations", "doctor", userParsed["_id"].(bson.ObjectId).Hex())
		if err != nil {
			Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		Helpers.RespondWithJSON(w, http.StatusOK, agendaAnnotations)
	}
}

func findAgendaAnnotationsByPatientEndPoint(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	w.Header().Set("Content-type", "application/json")

	//fmt.Println("patient log" + params["patient"])

	agendaAnnotations, err := dao.FindManyByKey("agendaAnnotations", "patient", params["patient"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, agendaAnnotations)

}

func createAgendaAnnotationEndPoint(w http.ResponseWriter, r *http.Request) {

	user := context.Get(r, "user")

	userParsed := user.(bson.M)

	userType := context.Get(r, "userType")

	defer r.Body.Close()
	w.Header().Set("Content-type", "application/json")

	err, agendaAnnotation := agendaAnnotationValidator(r)

	if len(err["validationError"].(url.Values)) > 0 {
		//fmt.Println(len(e))
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	agendaAnnotation.ID = bson.NewObjectId()
	agendaAnnotation.Date = time.Now().String()
	agendaAnnotation.UpdateDate = time.Now().String()
	agendaAnnotation.CreatedBy = userParsed["_id"].(bson.ObjectId).Hex()
	agendaAnnotation.UpdatedBy = userParsed["_id"].(bson.ObjectId).Hex()

	if userType.(int) == 2 {
		agendaAnnotation.Doctor = userParsed["_id"].(bson.ObjectId).Hex()
	}

	if err := dao.Insert("agendaAnnotations", agendaAnnotation, nil); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusCreated, agendaAnnotation)

}

func findAgendaAnnotationEndPoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	agendaAnnotation, err := dao.FindByID("agendaAnnotations", params["id"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid AgendaAnnotation ID")
		return
	}
	Helpers.RespondWithJSON(w, http.StatusOK, agendaAnnotation)

}

func removeAgendaAnnotationEndPoint(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	err := dao.DeleteByID("agendaAnnotations", params["id"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid AgendaAnnotation ID")
		return
	}
	Helpers.RespondWithJSON(w, http.StatusOK, nil)

}

func updateAgendaAnnotationEndPoint(w http.ResponseWriter, r *http.Request) {

	//fmt.Printf("agenda update end point")

	user := context.Get(r, "user")

	userParsed := user.(bson.M)

	defer r.Body.Close()
	params := mux.Vars(r)

	w.Header().Set("Content-type", "application/json")

	err, agendaAnnotation := agendaAnnotationValidator(r)

	if len(err["validationError"].(url.Values)) > 0 {
		//fmt.Println(len(e))
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	prevData, err2 := dao.FindByID("agendaAnnotations", params["id"])
	if err2 != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid AgendaAnnotation ID")
		return
	}

	parsedData := prevData.(bson.M)

	agendaAnnotation.ID = parsedData["_id"].(bson.ObjectId)

	agendaAnnotation.Date = parsedData["date"].(string)

	agendaAnnotation.UpdateDate = time.Now().String()

	agendaAnnotation.CreatedBy = parsedData["createdBy"].(string)

	agendaAnnotation.UpdatedBy = userParsed["_id"].(bson.ObjectId).Hex()

	if err := dao.Update("agendaAnnotations", agendaAnnotation.ID, agendaAnnotation); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})

}
