package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Kagami/go-face"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"

	Helpers "github.com/sumaikun/clickal-rest-api/helpers"
)

//-------------------------------------- file Upload -----------------------------------------

func fileUpload(w http.ResponseWriter, r *http.Request) {

	fmt.Println("File Upload EndPoint Hit")

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		Helpers.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	defer file.Close()

	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	var extension = filepath.Ext(handler.Filename)

	fmt.Printf("Extension: %+v\n", extension)

	tempFile, err := ioutil.TempFile("files", "upload-*"+extension)

	if err != nil {
		fmt.Println(err)
		Helpers.RespondWithJSON(w, http.StatusInternalServerError, err)
	}

	var tempPath = tempFile.Name()

	fmt.Println("temp file before trim" + tempPath)

	var tempName = strings.Replace(tempPath, "files/", "", -1)

	fmt.Println("tempName " + tempName)

	defer tempFile.Close()

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		Helpers.RespondWithJSON(w, http.StatusInternalServerError, err)
	}
	// write this byte array to our temporary file
	tempFile.Write(fileBytes)

	Helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"filename": tempName})

}

func deleteImage(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var fileName = params["file"]

	var err = os.Remove("./files/" + fileName)
	if err != nil {
		//log.Fatal(err) // perhaps handle this nicer
		Helpers.RespondWithJSON(w, http.StatusInternalServerError, err)
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"status": "fileDelete"})
	return

}

func serveImage(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	var fileName = params["image"]

	if !strings.Contains(fileName, "png") && !strings.Contains(fileName, "jpg") && !strings.Contains(fileName, "jpeg") && !strings.Contains(fileName, "gif") {
		Helpers.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"result": "invalid file extension"})
		return
	}

	img, err := os.Open("./files/" + params["image"])
	if err != nil {
		//log.Fatal(err) // perhaps handle this nicer
		Helpers.RespondWithJSON(w, http.StatusInternalServerError, err)
		return
	}
	defer img.Close()
	w.Header().Set("Content-Type", "image/jpeg") // <-- set the content-type header
	io.Copy(w, img)

}

func downloadFile(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	var fileName = params["file"]

	/*fmt.Println("fileName " + fileName)

	download, err := os.Open("./files/upload-815043770.pdf")

	if err != nil {

		Helpers.RespondWithJSON(w, http.StatusInternalServerError, err)
		return
	}

	defer download.Close()

	contentType, err := getFileContentType(download)

	if err != nil {
		Helpers.RespondWithJSON(w, http.StatusInternalServerError, err)
		return
	}

	fmt.Println("detected contentType", contentType)

	w.Header().Set("Content-Type", "application/pdf")

	w.Header().Set("Content-Disposition: attachment", "filename=test.pdf")

	_, err = io.Copy(w, download)*/

	http.ServeFile(w, r, "./files/"+fileName)
}

func getFileContentType(out *os.File) (string, error) {

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}

// Enums --------------------------------------------------------------------

func userRoles(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-type", "application/json")

	x := [3]string{"ADMIN", "OPERATOR", "AUDITOR"}

	Helpers.RespondWithJSON(w, http.StatusOK, x)
}

func contactStratus(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-type", "application/json")

	x := [6]string{"estrato 1", "estrato 2", "estrato 3", "estrato 4", "estrato 5", "estrato 6"}

	Helpers.RespondWithJSON(w, http.StatusOK, x)
}

func contactDocumentType(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-type", "application/json")

	x := [4]string{"CC", "CE", "Pasaporte", "TI"}

	Helpers.RespondWithJSON(w, http.StatusOK, x)
}

func parametersType(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-type", "application/json")

	x := []string{"Tipo de especialización", "Ciudades de atención"}

	Helpers.RespondWithJSON(w, http.StatusOK, x)
}

func administrationWayType(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-type", "application/json")

	x := [7]string{"Oral", "Intravenosa", "Intramuscular", "Subcutanea", "tópica", "rectal", "inhalatoria"}

	Helpers.RespondWithJSON(w, http.StatusOK, x)
}

func presentationType(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-type", "application/json")

	x := [7]string{"Jarabes", "Gotas", "Capsulas", "Polvo", "Granulado", "Emulsión", "Bebible"}

	Helpers.RespondWithJSON(w, http.StatusOK, x)
}

func stateType(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-type", "application/json")

	x := [3]string{"ACTIVE", "INACTIVE", "PENDING"}

	Helpers.RespondWithJSON(w, http.StatusOK, x)
}

func recognizeFace(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	rec, err := face.NewRecognizer(modelsDir)
	if err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Free the resources when you're finished.
	defer rec.Close()

	// Test image with 10 faces.
	testImage := filepath.Join(imagesDir, "jesus.jpeg")
	// Recognize faces on that image.
	faces, err2 := rec.RecognizeFile(testImage)
	if err2 != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err2.Error())
		return
	}

	var samples []face.Descriptor
	var usersFaces []int32
	for i, f := range faces {
		samples = append(samples, f.Descriptor)
		// Each face is unique on that image so goes to its own category.
		usersFaces = append(usersFaces, int32(i))
	}

	fmt.Println("faces", len(faces))

	//fmt.Println("faces", faces[0].Descriptor)

	testImage2 := filepath.Join(imagesDir, "jenifer.jpeg")
	userFace, err3 := rec.RecognizeSingleFile(testImage2)
	if err3 != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err3.Error())
		return
	}
	if userFace == nil {
		Helpers.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"result": "Not a single face on the image"})
		return
	}

	//fmt.Println("face", userFace.Descriptor)

	rec.SetSamples(samples, usersFaces)

	userFaceID := rec.ClassifyThreshold(userFace.Descriptor, 0.1)
	if userFaceID < 0 {
		Helpers.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"result": "could not classify"})
		return
	}

	fmt.Println("userFaceID", userFaceID)

	Helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

//------------------------------------- LANDING PAGE FUNCTIONS ---------------------------------------

func doctorsLandingPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	var doctorsLanding []map[string]interface{}

	doctors, err := dao.FindDoctorsWithCitiesAndSpecializations()
	if err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	//fmt.Println(reflect.TypeOf(doctors))

	for _, doctor := range doctors {
		parsedDoctor := doctor.(bson.M)
		//fmt.Println(parsedDoctor["name"])

		doctorSettings, _ := dao.FindOneByKEY("doctorSettings", "doctor", parsedDoctor["_id"].(bson.ObjectId).Hex())

		//fmt.Println("doctorSettings", doctorSettings)

		doctorL := make(map[string]interface{})
		doctorL["name"] = parsedDoctor["name"]
		doctorL["lastName"] = parsedDoctor["lastName"]
		doctorL["CityDetails"] = parsedDoctor["cityDetails"]
		doctorL["aboutDoctor"] = parsedDoctor["aboutDoctor"]
		doctorL["id"] = parsedDoctor["_id"]
		doctorL["picture"] = parsedDoctor["picture"]
		doctorL["qualification"] = parsedDoctor["qualification"]
		doctorL["specialistDetails"] = parsedDoctor["specialistDetails"]
		doctorL["settings"] = doctorSettings
		doctorsLanding = append(doctorsLanding, doctorL)
	}

	Helpers.RespondWithJSON(w, http.StatusOK, doctorsLanding)
}

func doctorDaySchedule(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	w.Header().Set("Content-type", "application/json")

	//params["patient"]

	appointments, err := dao.FindAppointmentsByDateAndDoctor(params["doctor"], params["date"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	annotations, err := dao.FindAnnotationsByDateAndDoctor(params["doctor"], params["date"])
	if err != nil {
		Helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Helpers.RespondWithJSON(w, http.StatusOK, map[string]interface{}{"appointments": appointments, "annotations": annotations})
}
