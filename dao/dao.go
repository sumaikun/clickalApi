package dao

import (
	"log"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var db *mgo.Database

//MongoConnector struct for connections access
type MongoConnector struct {
	Server   string
	Database string
}

//Connect golang to mongo sb
func (mongo *MongoConnector) Connect() {
	session, err := mgo.Dial(mongo.Server)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(mongo.Database)
}

//FindAll from repository
func (mongo *MongoConnector) FindAll(collection string) ([]interface{}, error) {
	var data []interface{}
	err := db.C(collection).Find(bson.M{}).All(&data)
	return data, err
}

//FindAllWithUsers from repository
func (mongo *MongoConnector) FindAllWithUsers(collection string) ([]interface{}, error) {

	var data []interface{}

	/*query := []bson.M{{
	"$lookup": bson.M{
		"let":  bson.M{"userObjId": bson.M{"$toObjectId": "$createdBy"}},
		"from": "users",
		"pipeline": []bson.M{{
			"$match": bson.M{"$expr": bson.M{"$eq":[]string{"$_id","$$userObjId"}}},
		}},
		"as": "userDetails",
	}}}*/

	query := []bson.M{{
		"$lookup": bson.M{
			"let":  bson.M{"userObjId": "$createdBy"},
			"from": "users",
			"pipeline": []bson.M{{
				"$match": bson.M{"$expr": bson.M{"$eq": []interface{}{bson.M{"$toString": "$_id"}, "$$userObjId"}}},
			}},
			"as": "userDetails",
		}}, {
		"$project": bson.M{
			"userDetails._id":         0,
			"userDetails.role":        0,
			"userDetails.password":    0,
			"userDetails.date":        0,
			"userDetails.update_date": 0,
		},
	}}

	err := db.C(collection).Pipe(query).All(&data)

	return data, err
}

//FindAllWithCities from repository
func (mongo *MongoConnector) FindAllWithCities(collection string) ([]interface{}, error) {

	var data []interface{}

	/*query := []bson.M{{
	"$lookup": bson.M{
		"let":  bson.M{"userObjId": bson.M{"$toObjectId": "$createdBy"}},
		"from": "users",
		"pipeline": []bson.M{{
			"$match": bson.M{"$expr": bson.M{"$eq":[]string{"$_id","$$userObjId"}}},
		}},
		"as": "userDetails",
	}}}*/

	query := []bson.M{{
		"$lookup": bson.M{
			"let":  bson.M{"userObjId": "$createdBy"},
			"from": "users",
			"pipeline": []bson.M{{
				"$match": bson.M{"$expr": bson.M{"$eq": []interface{}{bson.M{"$toString": "$_id"}, "$$userObjId"}}},
			}},
			"as": "userDetails",
		}}, {
		"$lookup": bson.M{
			"let":  bson.M{"cityObjId": "$city"},
			"from": "cityTypes",
			"pipeline": []bson.M{{
				"$match": bson.M{"$expr": bson.M{"$eq": []interface{}{bson.M{"$toString": "$_id"}, "$$cityObjId"}}},
			}},
			"as": "cityDetails",
		}}, {
		"$project": bson.M{
			"userDetails._id":         0,
			"userDetails.role":        0,
			"userDetails.password":    0,
			"userDetails.date":        0,
			"userDetails.update_date": 0,
		},
	}}

	err := db.C(collection).Pipe(query).All(&data)

	return data, err
}

//FindAllWithPatients from repository
func (mongo *MongoConnector) FindAllWithPatients(collection string) ([]interface{}, error) {

	var data []interface{}

	/*query := []bson.M{{
	"$lookup": bson.M{
		"let":  bson.M{"userObjId": bson.M{"$toObjectId": "$createdBy"}},
		"from": "users",
		"pipeline": []bson.M{{
			"$match": bson.M{"$expr": bson.M{"$eq":[]string{"$_id","$$userObjId"}}},
		}},
		"as": "userDetails",
	}}}*/

	query := []bson.M{{
		"$lookup": bson.M{
			"let":  bson.M{"doctorObjId": "$doctor"},
			"from": "doctors",
			"pipeline": []bson.M{{
				"$match": bson.M{"$expr": bson.M{"$eq": []interface{}{bson.M{"$toString": "$_id"}, "$$doctorObjId"}}},
			}},
			"as": "doctorDetails",
		}}, {
		"$lookup": bson.M{
			"let":  bson.M{"patientObjId": "$patient"},
			"from": "patients",
			"pipeline": []bson.M{{
				"$match": bson.M{"$expr": bson.M{"$eq": []interface{}{bson.M{"$toString": "$_id"}, "$$patientObjId"}}},
			}},
			"as": "patientDetails",
		}}, {
		"$project": bson.M{
			"userDetails._id":         0,
			"userDetails.role":        0,
			"userDetails.password":    0,
			"userDetails.date":        0,
			"userDetails.update_date": 0,
		},
	}}

	err := db.C(collection).Pipe(query).All(&data)

	return data, err
}

//FindManyByKey from repository
func (mongo *MongoConnector) FindManyByKey(collection string, key string, value string) ([]interface{}, error) {
	var data []interface{}

	query := []bson.M{
		{
			"$match": bson.M{key: value},
		}, {
			"$lookup": bson.M{
				"let":  bson.M{"userObjId": "$createdBy"},
				"from": "users",
				"pipeline": []bson.M{{
					"$match": bson.M{"$expr": bson.M{"$eq": []interface{}{bson.M{"$toString": "$_id"}, "$$userObjId"}}},
				}},
				"as": "userDetails",
			}}, {
			"$lookup": bson.M{
				"let":  bson.M{"patientObjId": "$patient"},
				"from": "patients",
				"pipeline": []bson.M{{
					"$match": bson.M{"$expr": bson.M{"$eq": []interface{}{bson.M{"$toString": "$_id"}, "$$patientObjId"}}},
				}},
				"as": "userDetails",
			}}, {
			"$project": bson.M{
				"userDetails._id":         0,
				"userDetails.role":        0,
				"userDetails.password":    0,
				"userDetails.date":        0,
				"userDetails.update_date": 0,
			},
		}}

	err := db.C(collection).Pipe(query).All(&data)
	return data, err
}

//FindManyByKeyWithPatiens from repository
func (mongo *MongoConnector) FindManyByKeyWithPatiens(collection string, key string, value string) ([]interface{}, error) {
	var data []interface{}

	query := []bson.M{
		{
			"$match": bson.M{key: value},
		}, {
			"$lookup": bson.M{
				"let":  bson.M{"userObjId": "$createdBy"},
				"from": "users",
				"pipeline": []bson.M{{
					"$match": bson.M{"$expr": bson.M{"$eq": []interface{}{bson.M{"$toString": "$_id"}, "$$userObjId"}}},
				}},
				"as": "userDetails",
			}}, {
			"$lookup": bson.M{
				"let":  bson.M{"doctorObjId": "$doctor"},
				"from": "doctors",
				"pipeline": []bson.M{{
					"$match": bson.M{"$expr": bson.M{"$eq": []interface{}{bson.M{"$toString": "$_id"}, "$$doctorObjId"}}},
				}},
				"as": "doctorDetails",
			}}, {
			"$lookup": bson.M{
				"let":  bson.M{"patientObjId": "$patient"},
				"from": "patients",
				"pipeline": []bson.M{{
					"$match": bson.M{"$expr": bson.M{"$eq": []interface{}{bson.M{"$toString": "$_id"}, "$$patientObjId"}}},
				}},
				"as": "patientDetails",
			}}, {
			"$project": bson.M{
				"userDetails._id":         0,
				"userDetails.role":        0,
				"userDetails.password":    0,
				"userDetails.date":        0,
				"userDetails.update_date": 0,
			},
		}}

	err := db.C(collection).Pipe(query).All(&data)
	return data, err
}

//Insert into repository
func (mongo *MongoConnector) Insert(collection string, data interface{}, uniqueKeys []string) error {

	for _, key := range uniqueKeys {
		index := mgo.Index{
			Key:    []string{key},
			Unique: true,
		}
		if err := db.C(collection).EnsureIndex(index); err != nil {
			return err
		}
	}

	err := db.C(collection).Insert(&data)
	return err
}

//FindByID in repository
func (mongo *MongoConnector) FindByID(collection string, id string) (interface{}, error) {

	//fmt.Println(collection, id)

	var data interface{}
	err := db.C(collection).FindId(bson.ObjectIdHex(id)).One(&data)
	return data, err
}

// DeleteByID by id on repository
func (mongo *MongoConnector) DeleteByID(collection string, id string) error {

	err := db.C(collection).RemoveId(bson.ObjectIdHex(id))
	//err := db.C(COLLECTION).Remove(&movie)
	return err
}

//FindInArrayKey in repository
func (mongo *MongoConnector) FindInArrayKey(collection string, key string, id string) (interface{}, error) {

	var data []interface{}
	err := db.C(collection).Find(bson.M{
		key: bson.M{
			"$elemMatch": bson.M{"$eq": id},
		},
	}).All(&data)
	return data, err

}

// Update an existing collection
func (mongo *MongoConnector) Update(collection string, id interface{}, data interface{}) error {

	err := db.C(collection).UpdateId(id, &data)
	return err
}

// PartialUpdate an existing collection
func (mongo *MongoConnector) PartialUpdate(collection string, id string, data interface{}) error {

	err := db.C(collection).Update(bson.M{"_id": bson.ObjectIdHex(id)}, bson.M{"$set": data})
	return err
}

//FindOneByKEY with key and value specified in repository
func (mongo *MongoConnector) FindOneByKEY(collection string, key string, value string) (interface{}, error) {
	var data interface{}
	err := db.C(collection).Find(bson.M{key: value}).One(&data)
	return data, err
}

//FindByDateRange based on range by unique key
func (mongo *MongoConnector) FindByDateRange(collection string, key string, fromDate time.Time, toDate time.Time) (interface{}, error) {
	//fromDate := time.Date(2014, time.November, 4, 0, 0, 0, 0, time.UTC)
	//toDate := time.Date(2014, time.November, 5, 0, 0, 0, 0, time.UTC)

	var data []interface{}
	err := db.C(collection).Find(bson.M{
		key: bson.M{
			"$gt": fromDate,
			"$lt": toDate,
		},
	}).All(&data)
	return data, err
}

//FindByKeyLike find based on regex that is like an includes
func (mongo *MongoConnector) FindByKeyLike(collection string, key string, value string) (interface{}, error) {

	var data []interface{}
	err := db.C(collection).Find(bson.M{key: bson.RegEx{value + ".*", ""}}).All(&data)
	return data, err
}

/** Specific services  */

//FindAppointmentByDateAndPatient specific query for get day appointment
func (mongo *MongoConnector) FindAppointmentByDateAndPatient(patient string, date string) (interface{}, error) {

	var data []interface{}
	err := db.C("appointments").Find(bson.M{
		"$and": []bson.M{
			bson.M{"appointmentDate": bson.RegEx{date + ".*", ""}},
			bson.M{"patient": patient},
		},
	}).All(&data)
	return data, err
}

//FindAppointmentsByDateAndDoctor specific query for get day appointment
func (mongo *MongoConnector) FindAppointmentsByDateAndDoctor(doctor string, date string) (interface{}, error) {

	var data []interface{}
	err := db.C("appointments").Find(bson.M{
		"$and": []bson.M{
			bson.M{"appointmentDate": bson.RegEx{date + ".*", ""}},
			bson.M{"doctor": doctor},
		},
	}).All(&data)
	return data, err
}

//FindAnnotationsByDateAndDoctor specific query for get day annotations
func (mongo *MongoConnector) FindAnnotationsByDateAndDoctor(doctor string, date string) (interface{}, error) {

	var data []interface{}
	err := db.C("agendaAnnotations").Find(bson.M{
		"$and": []bson.M{
			bson.M{"AnnotationDate": bson.RegEx{date + ".*", ""}},
			bson.M{"doctor": doctor},
		},
	}).All(&data)
	return data, err
}

//FindDoctorsWithCitiesAndSpecializations from repository
func (mongo *MongoConnector) FindDoctorsWithCitiesAndSpecializations() ([]interface{}, error) {

	var data []interface{}

	query := []bson.M{{
		"$lookup": bson.M{
			"let":  bson.M{"specialisgObjId": "$specialistType"},
			"from": "specialistTypes",
			"pipeline": []bson.M{{
				"$match": bson.M{"$expr": bson.M{"$in": []interface{}{bson.M{"$toString": "$_id"}, "$$specialisgObjId"}}},
			}},
			"as": "specialistDetails",
		}}, {
		"$lookup": bson.M{
			"let":  bson.M{"cityObjId": "$city"},
			"from": "cityTypes",
			"pipeline": []bson.M{{
				"$match": bson.M{"$expr": bson.M{"$eq": []interface{}{bson.M{"$toString": "$_id"}, "$$cityObjId"}}},
			}},
			"as": "cityDetails",
		}}, {
		"$project": bson.M{
			"userDetails._id":         0,
			"userDetails.role":        0,
			"userDetails.password":    0,
			"userDetails.date":        0,
			"userDetails.update_date": 0,
		},
	}}

	err := db.C("doctors").Pipe(query).All(&data)

	return data, err
}
