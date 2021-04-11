package main

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	Models "github.com/sumaikun/clickal-rest-api/models"
	"github.com/thedevsaddam/govalidator"
	"gopkg.in/mgo.v2/bson"
)

func confirmAccountValidator(r *http.Request) (map[string]interface{}, Models.ConfirmAccount) {

	var confirm Models.ConfirmAccount

	rules := govalidator.MapData{
		"token": []string{"required"},
	}

	opts := govalidator.Options{
		Request:         r,
		Data:            &confirm,
		Rules:           rules,
		RequiredDefault: true,
	}

	v := govalidator.New(opts)
	e := v.ValidateJSON()
	//fmt.Println(user)

	err := map[string]interface{}{"validationError": e}

	return err, confirm
}

func resetPasswordValidator(r *http.Request) (map[string]interface{}, Models.ResetPassword) {

	var reset Models.ResetPassword

	rules := govalidator.MapData{
		"password": []string{"required"},
		"token":    []string{"required"},
	}

	opts := govalidator.Options{
		Request:         r,
		Data:            &reset,
		Rules:           rules,
		RequiredDefault: true,
	}

	v := govalidator.New(opts)
	e := v.ValidateJSON()
	//fmt.Println(user)

	err := map[string]interface{}{"validationError": e}

	return err, reset
}

func forgotPasswordValidator(r *http.Request) (map[string]interface{}, Models.ForgotPassword) {

	var forgot Models.ForgotPassword

	rules := govalidator.MapData{
		"email": []string{"required"},
	}

	opts := govalidator.Options{
		Request:         r,
		Data:            &forgot,
		Rules:           rules,
		RequiredDefault: true,
	}

	v := govalidator.New(opts)
	e := v.ValidateJSON()
	//fmt.Println(user)

	err := map[string]interface{}{"validationError": e}

	return err, forgot
}

func userValidator(r *http.Request) (map[string]interface{}, Models.User) {

	var user Models.User

	rules := govalidator.MapData{
		"name":           []string{"required"},
		"lastName":       []string{"required"},
		"email":          []string{"required", "email"},
		"phone":          []string{"required", "min:7", "max:10"},
		"address":        []string{"required"},
		"role":           []string{"required", "roleEnum"},
		"typeId":         []string{"required", "documentTypeEnum"},
		"identification": []string{"required"},
		"birthDate":      []string{"required"},
		"city":           []string{"required", "cityParam"},
		"state":          []string{"required", "stateEnum"},
	}

	opts := govalidator.Options{
		Request:         r,
		Data:            &user,
		Rules:           rules,
		RequiredDefault: true,
	}

	v := govalidator.New(opts)
	e := v.ValidateJSON()
	//fmt.Println(user)

	err := map[string]interface{}{"validationError": e}

	return err, user
}

func doctorValidator(r *http.Request) (map[string]interface{}, Models.Doctor) {

	var user Models.Doctor

	rules := govalidator.MapData{
		"name":           []string{"required"},
		"lastName":       []string{"required"},
		"email":          []string{"required", "email"},
		"phone":          []string{"required", "min:7", "max:10"},
		"address":        []string{"required"},
		"birthDate":      []string{"required"},
		"city":           []string{"required", "cityParam"},
		"specialistType": []string{"required", "specialistTypeParam"},
		"typeId":         []string{"required", "documentTypeEnum"},
		"identification": []string{"required"},
		"state":          []string{"stateEnum"},
		"medicalCenter":  []string{"medicalCenterParam"},
	}

	opts := govalidator.Options{
		Request:         r,
		Data:            &user,
		Rules:           rules,
		RequiredDefault: true,
	}

	v := govalidator.New(opts)
	e := v.ValidateJSON()
	//fmt.Println(user)

	err := map[string]interface{}{"validationError": e}

	return err, user
}

func patientValidator(r *http.Request) (map[string]interface{}, Models.Patient) {

	var contact Models.Patient

	rules := govalidator.MapData{
		"name":           []string{"required"},
		"lastName":       []string{"required"},
		"email":          []string{"required", "email"},
		"phone":          []string{"required", "min:7", "max:10"},
		"address":        []string{"required"},
		"stratus":        []string{"required"},
		"city":           []string{"required", "cityParam"},
		"ocupation":      []string{"required"},
		"typeId":         []string{"required", "documentTypeEnum"},
		"identification": []string{"required"},
		"birthDate":      []string{"required"},
		"sex":            []string{"required", "sexTypeEnum"},
	}

	opts := govalidator.Options{
		Request:         r,
		Data:            &contact,
		Rules:           rules,
		RequiredDefault: true,
	}

	v := govalidator.New(opts)
	e := v.ValidateJSON()
	//fmt.Println(user)

	err := map[string]interface{}{"validationError": e}

	return err, contact
}

func userRegisterValidator(r *http.Request) (map[string]interface{}, Models.UserRegister) {

	var userRegister Models.UserRegister

	rules := govalidator.MapData{
		"name":      []string{"required"},
		"lastName":  []string{"required"},
		"email":     []string{"required", "email"},
		"phone":     []string{"required", "min:7", "max:10"},
		"password":  []string{"required"},
		"confirmed": []string{"bool"},
		"city":      []string{"required", "cityParam"},
	}

	opts := govalidator.Options{
		Request:         r,
		Data:            &userRegister,
		Rules:           rules,
		RequiredDefault: true,
	}

	v := govalidator.New(opts)
	e := v.ValidateJSON()
	//fmt.Println(user)

	err := map[string]interface{}{"validationError": e}

	return err, userRegister
}

func medicalCenterValidator(r *http.Request) (map[string]interface{}, Models.MedicalCenters) {

	var product Models.MedicalCenters

	rules := govalidator.MapData{
		"name":           []string{"required"},
		"phone":          []string{"required", "min:7", "max:10"},
		"address":        []string{"required"},
		"identification": []string{"required"},
	}

	opts := govalidator.Options{
		Request:         r,
		Data:            &product,
		Rules:           rules,
		RequiredDefault: true,
	}

	v := govalidator.New(opts)
	e := v.ValidateJSON()
	//fmt.Println(user)

	err := map[string]interface{}{"validationError": e}

	return err, product
}

func specialistTypesValidator(r *http.Request) (map[string]interface{}, Models.SpecialistTypes) {

	var parameters Models.SpecialistTypes

	rules := govalidator.MapData{
		"name": []string{"required"},
	}

	opts := govalidator.Options{
		Request:         r,
		Data:            &parameters,
		Rules:           rules,
		RequiredDefault: true,
	}

	v := govalidator.New(opts)
	e := v.ValidateJSON()
	//fmt.Println(user)

	err := map[string]interface{}{"validationError": e}

	return err, parameters
}

func cityTypesValidator(r *http.Request) (map[string]interface{}, Models.CitiesTypes) {

	var parameters Models.CitiesTypes

	rules := govalidator.MapData{
		"name": []string{"required"},
	}

	opts := govalidator.Options{
		Request:         r,
		Data:            &parameters,
		Rules:           rules,
		RequiredDefault: true,
	}

	v := govalidator.New(opts)
	e := v.ValidateJSON()
	//fmt.Println(user)

	err := map[string]interface{}{"validationError": e}

	return err, parameters
}

func validatorSelector(r *http.Request, entity string) (map[string]interface{}, interface{}, []string) {

	var err map[string]interface{} = nil

	switch entity {

	case "specialistTypes":
		err, data := specialistTypesValidator(r)
		if len(err["validationError"].(url.Values)) == 0 {
			data.ID = bson.NewObjectId()
			data.Date = time.Now().String()
			data.UpdateDate = time.Now().String()
		}
		fmt.Println(data)
		return err, data, []string{"name"}

	case "cityTypes":
		err, data := cityTypesValidator(r)
		if len(err["validationError"].(url.Values)) == 0 {
			data.ID = bson.NewObjectId()
			data.Date = time.Now().String()
			data.UpdateDate = time.Now().String()
		}
		fmt.Println(data)
		return err, data, []string{"name"}

	}

	return err, nil, nil

}

func validatorSelectorUpdate(r *http.Request, entity string, prevData bson.M) (map[string]interface{}, interface{}, interface{}) {

	var err map[string]interface{} = nil

	switch entity {

	case "specialistTypes":
		err, data := specialistTypesValidator(r)
		if len(err["validationError"].(url.Values)) == 0 {
			data.ID = prevData["_id"].(bson.ObjectId)
			data.Date = prevData["date"].(string)
			data.UpdateDate = time.Now().String()
		}

		return err, data, data.ID

	case "cityTypes":
		err, data := cityTypesValidator(r)
		if len(err["validationError"].(url.Values)) == 0 {
			data.ID = prevData["_id"].(bson.ObjectId)
			data.Date = prevData["date"].(string)
			data.UpdateDate = time.Now().String()
		}

		return err, data, data.ID

	}

	return err, nil, nil

}

///////////////////////////////////////////////////////////////////////

func patientReviewValidator(r *http.Request) (map[string]interface{}, Models.PatientReview) {

	var patientReview Models.PatientReview

	rules := govalidator.MapData{
		"patient": []string{"required"},
	}

	opts := govalidator.Options{
		Request:         r,
		Data:            &patientReview,
		Rules:           rules,
		RequiredDefault: true,
	}

	v := govalidator.New(opts)
	e := v.ValidateJSON()
	//fmt.Println(user)

	err := map[string]interface{}{"validationError": e}

	return err, patientReview
}

///////////////////////////////////////////////////////////////////////

func physiologicalConstantsValidator(r *http.Request) (map[string]interface{}, Models.PhysiologicalConstants) {

	var physiologicalConstants Models.PhysiologicalConstants

	rules := govalidator.MapData{
		"patient":         []string{"required"},
		"bloodPressure":   []string{"required"},
		"heartRate":       []string{"required"},
		"respiratoryRate": []string{"required"},
		"heartBeat":       []string{"required"},
		"temperature":     []string{"required"},
		"weight":          []string{"required"},
		"height":          []string{"required"},
		"hidrationStatus": []string{"required", "hidrationStatusEnum"},
	}

	opts := govalidator.Options{
		Request:         r,
		Data:            &physiologicalConstants,
		Rules:           rules,
		RequiredDefault: true,
	}

	v := govalidator.New(opts)
	e := v.ValidateJSON()
	//fmt.Println(user)

	err := map[string]interface{}{"validationError": e}

	return err, physiologicalConstants
}

//////////////////////////////////////////////////////////////////////

func appointmentsValidator(r *http.Request) (map[string]interface{}, Models.Appointments) {

	var appointments Models.Appointments

	rules := govalidator.MapData{
		"patient":                      []string{"required"},
		"reasonForConsultation":        []string{"required"},
		"resultsForConsultation":       []string{"required"},
		"diagnosticCode":               []string{"required"},
		"appointmentDate":              []string{"required"},
		"state":                        []string{"required"},
		"haveMedicalTest":              []string{"bool"},
		"medicalReasonForConsultation": []string{"required"},
	}

	opts := govalidator.Options{
		Request:         r,
		Data:            &appointments,
		Rules:           rules,
		RequiredDefault: true,
	}

	v := govalidator.New(opts)
	e := v.ValidateJSON()
	//fmt.Println(user)

	err := map[string]interface{}{"validationError": e}

	return err, appointments
}

//////////////////////////////////////////////////////////////////////

func appointmentsScheduleValidator(r *http.Request) (map[string]interface{}, Models.Appointments) {

	var appointments Models.Appointments

	rules := govalidator.MapData{
		"patient":          []string{"required"},
		"appointmentDate":  []string{"required"},
		"state":            []string{"required"},
		"agendaAnnotation": []string{"required"},
	}

	opts := govalidator.Options{
		Request:         r,
		Data:            &appointments,
		Rules:           rules,
		RequiredDefault: true,
	}

	v := govalidator.New(opts)
	e := v.ValidateJSON()
	//fmt.Println(user)

	err := map[string]interface{}{"validationError": e}

	return err, appointments
}

//////////////////////////////////////////////////////////////////////

func medicinesValidator(r *http.Request) (map[string]interface{}, Models.Medicines) {

	var medicines Models.Medicines

	rules := govalidator.MapData{
		"patient":           []string{"required"},
		"administrationWay": []string{"required"},
		"duration":          []string{"required"},
		"posology":          []string{"required"},
		"presentation":      []string{"required"},
		"product":           []string{"required"},
		"appointment":       []string{"required"},
	}

	opts := govalidator.Options{
		Request:         r,
		Data:            &medicines,
		Rules:           rules,
		RequiredDefault: true,
	}

	v := govalidator.New(opts)
	e := v.ValidateJSON()
	//fmt.Println(user)

	err := map[string]interface{}{"validationError": e}

	return err, medicines
}

//////////////////////////////////////////////////////////////////////

func patientsFilesValidator(r *http.Request) (map[string]interface{}, Models.PatientFiles) {

	var patientFile Models.PatientFiles

	rules := govalidator.MapData{
		"patient":     []string{"required"},
		"filePath":    []string{"required"},
		"description": []string{"required"},
	}

	opts := govalidator.Options{
		Request:         r,
		Data:            &patientFile,
		Rules:           rules,
		RequiredDefault: true,
	}

	v := govalidator.New(opts)
	e := v.ValidateJSON()
	//fmt.Println(user)

	err := map[string]interface{}{"validationError": e}

	return err, patientFile
}

//////////////////////////////////////////////////////////////////////

func agendaAnnotationValidator(r *http.Request) (map[string]interface{}, Models.AgendaAnnotation) {

	var agendaAnnotation Models.AgendaAnnotation

	rules := govalidator.MapData{
		"annotationDate":   []string{"required"},
		"annotationToDate": []string{"required"},
		"description":      []string{"required"},
	}

	opts := govalidator.Options{
		Request:         r,
		Data:            &agendaAnnotation,
		Rules:           rules,
		RequiredDefault: true,
	}

	v := govalidator.New(opts)
	e := v.ValidateJSON()
	//fmt.Println(user)

	err := map[string]interface{}{"validationError": e}

	return err, agendaAnnotation
}

//////////////////////////////////////////////////////////////////////

func doctorSettingsValidator(r *http.Request) (map[string]interface{}, Models.DoctorSettings) {

	var doctorSettings Models.DoctorSettings

	rules := govalidator.MapData{
		"hoursRange":   []string{"required", "hoursRangeType"},
		"daysRange":    []string{"required", "daysRangeType"},
		"isScheduling": []string{"bool"},
		"doctor":       []string{"required", "doctorParam"},
	}

	opts := govalidator.Options{
		Request:         r,
		Data:            &doctorSettings,
		Rules:           rules,
		RequiredDefault: true,
	}

	v := govalidator.New(opts)
	e := v.ValidateJSON()
	//fmt.Println(user)

	err := map[string]interface{}{"validationError": e}

	return err, doctorSettings
}

//////////////////////////////////////////////////////////////////////

func customerRegisterByLandingAppointment(r *http.Request) (map[string]interface{}, Models.PatientAppointment) {

	var contact Models.PatientAppointment

	rules := govalidator.MapData{
		"name":            []string{"required"},
		"lastName":        []string{"required"},
		"email":           []string{"required", "email"},
		"phone":           []string{"required", "min:7", "max:10"},
		"city":            []string{"required", "cityParam"},
		"ocupation":       []string{"required"},
		"typeId":          []string{"required", "documentTypeEnum"},
		"identification":  []string{"required"},
		"appointmentDate": []string{"required"},
		"doctor":          []string{"required", "doctorParam"},
	}

	opts := govalidator.Options{
		Request:         r,
		Data:            &contact,
		Rules:           rules,
		RequiredDefault: true,
	}

	v := govalidator.New(opts)
	e := v.ValidateJSON()
	//fmt.Println(user)

	err := map[string]interface{}{"validationError": e}

	return err, contact
}
