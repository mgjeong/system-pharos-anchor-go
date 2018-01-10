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
package image

import (
	"commons/errors"
	"commons/logger"
	. "db/mongo/wrapper"
	"gopkg.in/mgo.v2/bson"
)

type Command interface {
	// AddDockerImage insert a new docker image.
	AddDockerImage(image map[string]interface{}) (map[string]interface{}, error)

	// DeleteDockerImage delete a specific image from db related to image.
	DeleteDockerImage(imageId string) error

	// UpdateDockerImage update status of docker image from db related to image.
	UpdateDockerImage(imageId string, image map[string]interface{}) error

	// GetDockerImage returns a single document from db related to image.
	GetDockerImage(imageId string) (map[string]interface{}, error)
}

const (
	DB_NAME          = "DeploymentManagerDB"
	IMAGE_COLLECTION = "IMAGE"
	DB_URL           = "127.0.0.1:27017"
)

type Image struct {
	ID         bson.ObjectId `bson:"_id,omitempty"`
	Repository string
	Tag        string
	Size       string
	Action     string
	Timestamp  string
}

type Executor struct {}

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

// convertToMap converts Image object into a map.
func (image Image) convertToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":         image.ID.Hex(),
		"repository": image.Repository,
		"tag":        image.Tag,
		"size":       image.Size,
		"action":     image.Action,
		"timestamp":  image.Timestamp,
	}
}

// AddDockerImage insert a new docker image to 'image' collection.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (Executor) AddDockerImage(image map[string]interface{}) (map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	session, err := connect(DB_URL)
	if err != nil {
		return nil, err
	}
	defer close(session)

	newImage := Image{
		ID:         bson.NewObjectId(),
		Repository: image["repository"].(string),
		Tag:        image["tag"].(string),
		Size:       image["size"].(string),
		Action:     image["action"].(string),
		Timestamp:  image["timestamp"].(string),
	}

	err = getCollection(session, DB_NAME, IMAGE_COLLECTION).Insert(newImage)

	if err != nil {
		return nil, ConvertMongoError(err)
	}

	result := newImage.convertToMap()
	return result, err
}

// UpdateDockerImage update status of docker image specified by imageId parameter.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (Executor) UpdateDockerImage(imageId string, image map[string]interface{}) error {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	session, err := connect(DB_URL)
	if err != nil {
		return err
	}
	defer close(session)

	// Verify id is ObjectId, otherwise fail
	if !bson.IsObjectIdHex(imageId) {
		err := errors.InvalidObjectId{imageId}
		return err
	}

	query := bson.M{"_id": bson.ObjectIdHex(imageId), "repository": image["repository"].(string)}
	update := bson.M{
		"$set": bson.M{
			"tag":        image["tag"].(string),
			"size":       image["size"].(string),
			"action":     image["action"].(string),
			"timestamp":  image["timestamp"].(string)}}

	err = getCollection(session, DB_NAME, IMAGE_COLLECTION).Update(query, update)
	if err != nil {
		return ConvertMongoError(err, "Failed to update address")
	}
	return err
}

// GetDockerImage returns a single document specified by imageId parameter.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (Executor) GetDockerImage(imageId string) (map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	session, err := connect(DB_URL)
	if err != nil {
		return nil, err
	}
	defer close(session)

	// Verify id is ObjectId, otherwise fail
	if !bson.IsObjectIdHex(imageId) {
		err := errors.InvalidObjectId{imageId}
		return nil, err
	}

	image := Image{}
	query := bson.M{"_id": bson.ObjectIdHex(imageId)}
	err = getCollection(session, DB_NAME, IMAGE_COLLECTION).Find(query).One(&image)
	if err != nil {
		return nil, ConvertMongoError(err, imageId)
	}

	result := image.convertToMap()
	return result, err
}

// DeleteDockerImage delete a single document from 'image' collection.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (Executor) DeleteDockerImage(imageId string) error {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	session, err := connect(DB_URL)
	if err != nil {
		return err
	}
	defer close(session)

	// Verify id is ObjectId, otherwise fail
	if !bson.IsObjectIdHex(imageId) {
		err = errors.InvalidObjectId{imageId}
		return err
	}

	query := bson.M{"_id": bson.ObjectIdHex(imageId)}
	err = getCollection(session, DB_NAME, IMAGE_COLLECTION).Remove(query)
	if err != nil {
		return ConvertMongoError(err, imageId)
	}
	return err
}
