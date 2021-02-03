package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"

	Helpers "github.com/sumaikun/clickal-rest-api/helpers"
)

//-------------------------------------- Parameters Functions --------------------------------

func createParameterEndPoint(w http.ResponseWriter, r *http.Request) {

	entity := strings.Replace(r.URL.Path, "/", "", -1)

	defer r.Body.Close()
	w.Header().Set("Content-type", "application/json")

	err, parameter, uniqueKeys := validatorSelector(r, entity)

	//fmt.Println(parameter)

	if len(err["validationError"].(url.Values)) > 0 {
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	if err := dao.Insert(entity, parameter, uniqueKeys); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusCreated, parameter)

}

func allParametersEndPoint(w http.ResponseWriter, r *http.Request) {
	entity := strings.Replace(r.URL.Path, "/", "", -1)
	w.Header().Set("Content-type", "application/json")

	parameters, err := dao.FindAll(entity)
	if err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, parameters)

}

func findParameterEndPoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	entity := strings.Replace(r.URL.Path, "/"+params["id"], "", -1)
	entity = strings.Replace(entity, "/", "", -1)

	parameter, err := dao.FindByID(entity, params["id"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Parameter ID")
		return
	}
	Helpers.RespondWithJSON(w, http.StatusOK, parameter)
}

func deleteParameterEndPoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	entity := strings.Replace(r.URL.Path, "/"+params["id"], "", -1)
	entity = strings.Replace(entity, "/", "", -1)
	err := dao.DeleteByID(entity, params["id"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Parameter ID")
		return
	}
	Helpers.RespondWithJSON(w, http.StatusOK, nil)
}

func updateParameterEndPoint(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	params := mux.Vars(r)
	entity := strings.Replace(r.URL.Path, "/"+params["id"], "", -1)
	entity = strings.Replace(entity, "/", "", -1)
	w.Header().Set("Content-type", "application/json")

	prevData, err2 := dao.FindByID(entity, params["id"])

	if err2 != nil {
		fmt.Println(err2)
		fmt.Println(params["id"])
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Parameter ID")
		return
	}

	parsedData := prevData.(bson.M)

	err, data, dataID := validatorSelectorUpdate(r, entity, parsedData)

	if len(err["validationError"].(url.Values)) > 0 {
		//fmt.Println(len(e))
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	if err := dao.Update(entity, dataID, data); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, "invalid")
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})

}
