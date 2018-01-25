/*******************************************************************************
 * Copyright 2017 Samsung Electronics All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 *******************************************************************************/

package app

import (
	"commons/errors"
	"commons/logger"
	. "db/mongo/wrapper"
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/yaml.v2"
	"strings"
)

// Interface of App model's operations.
type Command interface {
	// AddApp insert a deployed application information.
	AddApp(appId string, description []byte) error

	// GetApps returns all matches for the query-string which is passed in call to function.
	GetApps(queryOptional ...map[string]interface{}) ([]map[string]interface{}, error)

	// DeleteApp delete a deployed application information.
	DeleteApp(appId string) error
}

const (
	DB_NAME        = "DeploymentManagerDB"
	APP_COLLECTION = "APP"
	SERVICES_FIELD = "services"
	IMAGE_FIELD    = "image"
	DB_URL         = "localhost:27017"
)

type App struct {
	ID       string `bson:"_id,omitempty"`
	Images   []string
	Services []string
	RefCnt   int
}

type Executor struct {
}

var mgoDial Connection

func init() {
	mgoDial = MongoDial{}
}

// Try to connect with mongo db server.
// if succeed to connect with mongo db server, return error as nil,
// otherwise, return error.
func connect(url string) (Session, error) {
	// Create a MongoDB Session
	session, err := mgoDial.Dial(url)

	if err != nil {
		return nil, ConvertMongoError(err, "")
	}

	return session, err
}

// close of mongodb session.
func close(mgoSession Session) {
	mgoSession.Close()
}

// Getting collection by name.
// return mongodb Collection
func getCollection(mgoSession Session, dbname string, collectionName string) Collection {
	return mgoSession.DB(dbname).C(collectionName)
}

// Convert to map by object of struct App.
// will return App information as map.
func (app App) convertToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":       app.ID,
		"images":   app.Images,
		"services": app.Services,
	}
}

// AddApp insert a deployed application information.
// if succeed to add, return app information as map.
// otherwise, return error.
func (Executor) AddApp(appId string, description []byte) error {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	if len(appId) == 0 {
		err := errors.InvalidParam{"Invalid param error : app_id is empty."}
		return err
	}

	session, err := connect(DB_URL)
	if err != nil {
		return err
	}
	defer close(session)

	// Get application information specified by appId parameter.
	app := App{}
	query := bson.M{"_id": appId}
	err = getCollection(session, DB_NAME, APP_COLLECTION).Find(query).One(&app)
	if err != nil {
		err = ConvertMongoError(err)
		switch err.(type) {
		default:
			return err
		case errors.NotFound:
			images, services, err := getImageAndServiceNames(description)
			if err != nil {
				return err
			}

			// Add a newly deployed application information.
			app := App{
				ID:       appId,
				Images:   images,
				Services: services,
				RefCnt:   1,
			}

			err = getCollection(session, DB_NAME, APP_COLLECTION).Insert(app)
			if err != nil {
				return ConvertMongoError(err, "")
			}
			return nil
		}
	}

	// Increase the reference count.
	query = bson.M{"_id": appId}
	update := bson.M{"$set": bson.M{"refCnt": app.RefCnt + 1}}
	err = getCollection(session, DB_NAME, APP_COLLECTION).Update(query, update)
	if err != nil {
		return ConvertMongoError(err, "Failed to increase reference count")
	}
	return err
}

// GetApps returns all matches for the query-string which is passed in call to function.
// if succeed to get, return list of all app information as slice.
// otherwise, return error.
func (Executor) GetApps(queryOptional ...map[string]interface{}) ([]map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	session, err := connect(DB_URL)
	if err != nil {
		return nil, err
	}
	defer close(session)

	var query interface{}
	switch len(queryOptional) {
	case 1:
		for key, val := range queryOptional[0] {
			query = bson.M{key: bson.M{"$in": []string{val.(string)}}}
		}
	}

	apps := []App{}
	err = getCollection(session, DB_NAME, APP_COLLECTION).Find(query).All(&apps)
	if err != nil {
		err = ConvertMongoError(err, "Failed to get all apps")
		return nil, err
	}

	result := make([]map[string]interface{}, len(apps))
	for i, app := range apps {
		result[i] = app.convertToMap()
	}

	return result, err
}

// DeleteApp delete a deployed application information.
// if succeed to delete, return error as nil.
// otherwise, return error.
func (Executor) DeleteApp(appId string) error {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	if len(appId) == 0 {
		err := errors.InvalidParam{"Invalid param error : appId is empty."}
		return err
	}

	session, err := connect(DB_URL)
	if err != nil {
		return err
	}
	defer close(session)

	// Get application information specified by appId parameter.
	app := App{}
	query := bson.M{"_id": appId}
	err = getCollection(session, DB_NAME, APP_COLLECTION).Find(query).One(&app)
	if err != nil {
		return ConvertMongoError(err, "")
	}

	refCnt := app.RefCnt - 1
	if refCnt == 0 {
		err = getCollection(session, DB_NAME, APP_COLLECTION).Remove(bson.M{"_id": appId})
		if err != nil {
			errMsg := "Failed to remove a app by " + appId
			return ConvertMongoError(err, errMsg)
		}
		return nil
	}

	// Decrease the reference count.
	query = bson.M{"_id": appId}
	update := bson.M{"$set": bson.M{"refCnt": refCnt}}
	err = getCollection(session, DB_NAME, APP_COLLECTION).Update(query, update)
	if err != nil {
		return ConvertMongoError(err, "Failed to decrease reference count")
	}
	return err
}

func getImageAndServiceNames(source []byte) ([]string, []string, error) {
	var yamlData interface{}
	err := yaml.Unmarshal([]byte(source), &yamlData)
	if err != nil {
		return nil, nil, errors.InvalidYaml{"Invalid YAML error : description has not service information."}
	}

	jsonData, err := json.Marshal(convert(yamlData))
	if err != nil {
		return nil, nil, errors.InvalidYaml{"Invalid YAML error : description has not service information."}
	}

	description := make(map[string]interface{})
	err = json.Unmarshal(jsonData, &description)
	if err != nil {
		return nil, nil, convertJsonError(err)
	}

	if len(description[SERVICES_FIELD].(map[string]interface{})) == 0 || description[SERVICES_FIELD] == nil {
		return nil, nil, errors.InvalidYaml{"Invalid YAML error : description has not service information."}
	}

	var images []string
	var services []string
	for service_name, service_info := range description[SERVICES_FIELD].(map[string]interface{}) {
		services = append(services, service_name)

		if service_info.(map[string]interface{})[IMAGE_FIELD] == nil {
			return nil, nil, errors.InvalidYaml{"Invalid YAML error : description has not image information."}
		}

		fullImageName := service_info.(map[string]interface{})[IMAGE_FIELD].(string)
		var registryUrl, imageName string
		words := strings.Split(fullImageName, "/")
		if len(words) == 2 {
			registryUrl += words[0] + "/"
			imageName += words[1]
		} else {
			imageName += words[0]
		}

		words = strings.Split(imageName, ":")
		imageNameWithoutTag := words[0]
		images = append(images, registryUrl+imageNameWithoutTag)
	}
	return images, services, nil
}

// Converting to commons/errors by Json error
func convertJsonError(jsonError error) (err error) {
	switch jsonError.(type) {
	case *json.SyntaxError,
		*json.InvalidUTF8Error,
		*json.InvalidUnmarshalError,
		*json.UnmarshalFieldError,
		*json.UnmarshalTypeError:
		return errors.InvalidYaml{}
	default:
		return errors.Unknown{}
	}
}

// convert function changes the type of key from interface{} to string.
// yaml package unmarshal key-value pairs with map[interface{}]interface{}.
// but map[interface{}]interface{} type is not supported in json package.
// this function is available to resolve the problem.
func convert(in interface{}) interface{} {
	switch x := in.(type) {
	case map[interface{}]interface{}:
		out := map[string]interface{}{}
		for key, value := range x {
			out[key.(string)] = convert(value)
		}
		return out
	case []interface{}:
		for key, value := range x {
			x[key] = convert(value)
		}
	}
	return in
}
