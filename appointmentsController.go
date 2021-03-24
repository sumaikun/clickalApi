package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"

	Models "github.com/sumaikun/clickal-rest-api/models"

	Helpers "github.com/sumaikun/clickal-rest-api/helpers"
)

//--------------------------------Appointments functions ----------------------------------

func allAppointmentsEndPoint(w http.ResponseWriter, r *http.Request) {

	userType := context.Get(r, "userType")

	user := context.Get(r, "user")

	userParsed := user.(bson.M)

	if userType.(int) == 1 {
		w.Header().Set("Content-type", "application/json")

		appointments, err := dao.FindAllWithPatients("appointments")
		if err != nil {
			Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		Helpers.RespondWithJSON(w, http.StatusOK, appointments)
	}

	if userType.(int) == 2 {
		w.Header().Set("Content-type", "application/json")

		appointments, err := dao.FindManyByKeyWithPatiens("appointments", "doctor", userParsed["_id"].(bson.ObjectId).Hex())
		if err != nil {
			Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		Helpers.RespondWithJSON(w, http.StatusOK, appointments)
	}

}

func findAppointmentsByPatientEndPoint(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	w.Header().Set("Content-type", "application/json")

	appointments, err := dao.FindManyByKey("appointments", "patient", params["patient"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, appointments)

}

func appointmentsByPatientAndDateEndPoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	//fmt.Println("patient", params["patient"])

	//fmt.Println("date", params["date"])

	appointments, err := dao.FindAppointmentByDateAndPatient(params["patient"], params["date"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	fmt.Println("appointments", appointments)

	Helpers.RespondWithJSON(w, http.StatusOK, appointments)
}

func createAppointmentsEndPoint(w http.ResponseWriter, r *http.Request) {

	userType := context.Get(r, "userType")

	user := context.Get(r, "user")

	userParsed := user.(bson.M)

	w.Header().Set("Content-type", "application/json")

	// temporary buffer
	b := bytes.NewBuffer(make([]byte, 0))

	// TeeReader returns a Reader that writes to b what it reads from r.Body.
	reader := io.TeeReader(r.Body, b)

	var appointment Models.Appointments

	var err map[string]interface{}
	// Get the JSON body and decode into credentials
	err0 := json.NewDecoder(reader).Decode(&appointment)

	if err0 != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// we are done with body
	defer r.Body.Close()

	r.Body = ioutil.NopCloser(b)

	newID := bson.NewObjectId()

	if appointment.State == "PENDING" {
		fmt.Println("on pending")
		err, appointment = appointmentsScheduleValidator(r)

		fmt.Println("err", err, appointment)

		if len(err["validationError"].(url.Values)) > 0 {
			Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
			return
		} else {

			patient, _ := dao.FindByID("patients", appointment.Patient)

			parsedPatient := patient.(bson.M)

			expirationTime := time.Now().Add(24 * time.Hour)

			// set token for email
			claims := &Models.TypeClaims{
				Username: parsedPatient["email"].(string),
				Type:     "email-confirmation-" + newID.Hex(),
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: expirationTime.Unix(),
				},
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

			tokenString, _ := token.SignedString(jwtKey)

			dateInfo := strings.Split(appointment.AppointmentDate, " ")

			hour := dateInfo[1]

			//fmt.Println("dateInfo", dateInfo[0], hour[:len(hour)-3])

			go sendAppointmentConfirmationEmail(tokenString, parsedPatient["email"].(string), newID.Hex(), userParsed["name"].(string)+" "+userParsed["lastName"].(string), dateInfo[0], hour[:len(hour)-3])
		}
	} else {
		err, appointment = appointmentsValidator(r)

		if len(err["validationError"].(url.Values)) > 0 {
			//fmt.Println("appointment", appointment.State)
			Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
			return
		}

	}

	//fmt.Print("fappointment", appointment)

	appointment.ID = newID
	appointment.Date = time.Now().String()
	appointment.UpdateDate = time.Now().String()
	appointment.CreatedBy = userParsed["_id"].(bson.ObjectId).Hex()
	appointment.UpdatedBy = userParsed["_id"].(bson.ObjectId).Hex()

	if userType.(int) == 2 {
		appointment.Doctor = userParsed["_id"].(bson.ObjectId).Hex()
	}

	if err := dao.Insert("appointments", appointment, nil); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusCreated, appointment)

}

func findAppointmentsEndPoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	pet, err := dao.FindByID("appointments", params["id"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Appointment ID")
		return
	}
	Helpers.RespondWithJSON(w, http.StatusOK, pet)

}

func removeAppointmentsEndPoint(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	err := dao.DeleteByID("appointments", params["id"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Appointment ID")
		return
	}
	Helpers.RespondWithJSON(w, http.StatusOK, nil)

}

func updateAppointmentsEndPoint(w http.ResponseWriter, r *http.Request) {

	user := context.Get(r, "user")

	userParsed := user.(bson.M)

	// temporary buffer
	b := bytes.NewBuffer(make([]byte, 0))

	// TeeReader returns a Reader that writes to b what it reads from r.Body.
	reader := io.TeeReader(r.Body, b)

	var appointment Models.Appointments

	var err map[string]interface{}
	// Get the JSON body and decode into credentials
	err0 := json.NewDecoder(reader).Decode(&appointment)

	if err0 != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// we are done with body
	defer r.Body.Close()

	r.Body = ioutil.NopCloser(b)

	params := mux.Vars(r)

	w.Header().Set("Content-type", "application/json")

	if appointment.State == "PENDING" {
		err, appointment = appointmentsScheduleValidator(r)

		fmt.Println("err", err, appointment)

		if len(err["validationError"].(url.Values)) > 0 {
			Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
			return
		}
	} else {
		err, appointment = appointmentsValidator(r)

		if len(err["validationError"].(url.Values)) > 0 {
			//fmt.Println(len(e))
			Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
			return
		}
	}

	prevData, err2 := dao.FindByID("appointments", params["id"])
	if err2 != nil {
		Helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Appointment ID")
		return
	}

	parsedData := prevData.(bson.M)

	appointment.ID = parsedData["_id"].(bson.ObjectId)

	appointment.Date = parsedData["date"].(string)

	appointment.UpdateDate = time.Now().String()

	appointment.CreatedBy = parsedData["createdBy"].(string)

	appointment.UpdatedBy = userParsed["_id"].(bson.ObjectId).Hex()

	if err := dao.Update("appointments", appointment.ID, appointment); err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})

}

func confirmPatientAppointment(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	userType := context.Get(r, "userType")

	user := context.Get(r, "user")

	userParsed := user.(bson.M)

	if userType.(int) == 2 {
		sendEmailConfirmationToPatient(params["email"], userParsed["phone"].(string))
	}

	go dao.PartialUpdate("appointments", params["appointments"], bson.M{"state": "CONFIRMED"})

	Helpers.RespondWithJSON(w, http.StatusOK, nil)

}

func cancelPatientAppointment(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	userType := context.Get(r, "userType")

	user := context.Get(r, "user")

	userParsed := user.(bson.M)

	if userType.(int) == 2 {
		sendEmailCancelationToPatient(params["email"], userParsed["phone"].(string))
	}

	go dao.PartialUpdate("appointments", params["appointments"], bson.M{"state": "CANCELLED"})

	Helpers.RespondWithJSON(w, http.StatusOK, nil)

}
