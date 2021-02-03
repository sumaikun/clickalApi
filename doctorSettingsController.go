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

//-----------------------------  Doctor Settings functions --------------------------------------------------

func allDoctorSettingsEndPoint(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "application/json")

	doctorsSettings, err := dao.FindAll("doctorSettings")
	if err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, doctorsSettings)
}

func createDoctorSettingEndPoint(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	user := context.Get(r, "user")

	userParsed := user.(bson.M)

	w.Header().Set("Content-type", "application/json")

	err, doctorSettings := doctorSettingsValidator(r)

	if len(err["validationError"].(url.Values)) > 0 {
		//fmt.Println(len(e))
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	doctorSettings.ID = bson.NewObjectId()
	doctorSettings.Date = time.Now().String()
	doctorSettings.UpdateDate = time.Now().String()
	doctorSettings.CreatedBy = userParsed["_id"].(bson.ObjectId).Hex()
	doctorSettings.UpdatedBy = userParsed["_id"].(bson.ObjectId).Hex()

	if err := dao.Insert("doctorSettings", doctorSettings, []string{"doctor"}); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusCreated, doctorSettings)

}

func findDoctorSettingsEndPoint(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	doctorSettings, err := dao.FindByID("doctorSettings", params["id"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid DoctorSetting ID")
		return
	}
	Helpers.RespondWithJSON(w, http.StatusOK, doctorSettings)

}

func findDoctorSettingsByDoctorEndPoint(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	w.Header().Set("Content-type", "application/json")

	//fmt.Println("patient log" + params["patient"])

	doctorSettings, err := dao.FindOneByKEY("doctorSettings", "doctor", params["doctor"])
	if err != nil {
		//Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		Helpers.RespondWithJSON(w, http.StatusOK, nil)
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, doctorSettings)

}

func removeDoctorSettingsEndPoint(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	err := dao.DeleteByID("doctorSettings", params["id"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid DoctorSetting ID")
		return
	}
	Helpers.RespondWithJSON(w, http.StatusOK, nil)

}

func updateDoctorSettingsEndPoint(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	params := mux.Vars(r)

	usera := context.Get(r, "user")

	userParsed := usera.(bson.M)

	w.Header().Set("Content-type", "application/json")

	err, doctorSettings := doctorSettingsValidator(r)

	if len(err["validationError"].(url.Values)) > 0 {
		//fmt.Println(len(e))
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	prevDoctorSetting, err2 := dao.FindByID("doctorSettings", params["id"])
	if err2 != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid DoctorSettings ID")
		return
	}

	parsedData := prevDoctorSetting.(bson.M)

	doctorSettings.ID = parsedData["_id"].(bson.ObjectId)

	doctorSettings.Date = parsedData["date"].(string)

	doctorSettings.UpdateDate = time.Now().String()

	if parsedData["createdBy"] == nil {
		doctorSettings.CreatedBy = userParsed["_id"].(bson.ObjectId).Hex()
	} else {
		doctorSettings.CreatedBy = parsedData["createdBy"].(string)
	}

	doctorSettings.UpdatedBy = userParsed["_id"].(bson.ObjectId).Hex()

	if err := dao.Update("doctorSettings", doctorSettings.ID, doctorSettings); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})

}
