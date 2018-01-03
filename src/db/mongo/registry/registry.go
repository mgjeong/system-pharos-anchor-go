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
package registry

import (
	"commons/errors"
	"commons/logger"
	imageDB "db/mongo/image"
	. "db/mongo/wrapper"
	"gopkg.in/mgo.v2/bson"
)

type Command interface {
	// AddDockerRegistry insert a new docker registry information.
	AddDockerRegistry(url string) (map[string]interface{}, error)

	// GetDockerRegistries returns all documents from db related to docker registry.
	GetDockerRegistries() ([]map[string]interface{}, error)

	// GetDockerRegistry returns a single document from db related to docker registry.
	GetDockerRegistry(url string) (map[string]interface{}, error)

	// DeleteDockerRegistry delete a specific docker registry information from db related to registry.
	DeleteDockerRegistry(registryId string) error

	// AddDockerImages add a specific docker image to the target registry.
	AddDockerImages(registryId string, images []map[string]interface{}) error

	// GetDockerImages returns all docker images which belong to the target registry.
	GetDockerImages(registryId string) ([]map[string]interface{}, error)

	// UpdateDockerImage update status of docker image which belong to the target registry.
	UpdateDockerImage(registryId string, image map[string]interface{}) error

	// DeleteDockerImage delete a specific docker image from the target registry.
	DeleteDockerImage(registryId string, image map[string]interface{})
}

const (
	DB_NAME             = "DeploymentManagerDB"
	REGISTRY_COLLECTION = "REGISTRY"
	DB_URL              = "127.0.0.1:27017"
)

type Registry struct {
	ID     bson.ObjectId `bson:"_id,omitempty"`
	Url    string
	Images []string
}

type Executor struct {}

var mgoDial Connection
var imageExecutor imageDB.Command

func init() {
	mgoDial = MongoDial{}
	imageExecutor = imageDB.Executor{}
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

// convertToMap converts Registry object into a map.
func (registry Registry) convertToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":     registry.ID.Hex(),
		"url":    registry.Url,
		"images": registry.Images,
	}
}

// AddDockerRegistry insert a new docker registry information to 'registry' collection.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (Executor) AddDockerRegistry(url string) (map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	session, err := connect(DB_URL)
	if err != nil {
		return nil, err
	}
	defer close(session)

	newRegistry := Registry{
		ID:  bson.NewObjectId(),
		Url: url,
	}

	err = getCollection(session, DB_NAME, REGISTRY_COLLECTION).Insert(newRegistry)

	if err != nil {
		return nil, ConvertMongoError(err)
	}

	result := newRegistry.convertToMap()
	return result, err
}

// GetDockerRegistries returns all documents.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (Executor) GetDockerRegistries() ([]map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	session, err := connect(DB_URL)
	if err != nil {
		return nil, err
	}
	defer close(session)

	registries := []Registry{}
	err = getCollection(session, DB_NAME, REGISTRY_COLLECTION).Find(nil).All(&registries)
	if err != nil {
		return nil, ConvertMongoError(err)
	}

	result := make([]map[string]interface{}, len(registries))
	for i, registry := range registries {
		result[i] = registry.convertToMap()
		// Remove unused 'images' field from map.
		delete(result[i], "images")
	}
	return result, err
}

// GetDockerRegistry returns a single document specified by url parameter.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (Executor) GetDockerRegistry(url string) (map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	session, err := connect(DB_URL)
	if err != nil {
		return nil, err
	}
	defer close(session)

	registry := Registry{}
	query := bson.M{"url": url}
	err = getCollection(session, DB_NAME, REGISTRY_COLLECTION).Find(query).One(&registry)
	if err != nil {
		return nil, ConvertMongoError(err, url)
	}

	result := registry.convertToMap()
	// Remove unused 'images' field from map.
	delete(result, "images")
	return result, err
}

// DeleteDockerRegistry delete a single document from 'registry' collection.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (Executor) DeleteDockerRegistry(registryId string) error {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	session, err := connect(DB_URL)
	if err != nil {
		return err
	}
	defer close(session)

	// Verify id is ObjectId, otherwise fail
	if !bson.IsObjectIdHex(registryId) {
		err = errors.InvalidObjectId{registryId}
		return err
	}

	// Get a docker registry information specified by registryId parameter.
	registry := Registry{}
	query := bson.M{"_id": bson.ObjectIdHex(registryId)}
	err = getCollection(session, DB_NAME, REGISTRY_COLLECTION).Find(query).One(&registry)
	if err != nil {
		return ConvertMongoError(err, registryId)
	}

	// Delete a docker registry specified by registryId parameter.
	err = getCollection(session, DB_NAME, REGISTRY_COLLECTION).Remove(query)
	if err != nil {
		return ConvertMongoError(err, registryId)
	}

	// Delete all docker images included in the target registry from database.
	for _, imageId := range registry.Images {
		imageExecutor.DeleteDockerImage(imageId)
	}
	return nil
}

// AddDockerImages insert new images to 'image' collection.
// and add it to the list of image ids in the target 'registry' document.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (Executor) AddDockerImages(registryId string, images []map[string]interface{}) error {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	session, err := connect(DB_URL)
	if err != nil {
		return err
	}
	defer close(session)

	// Verify id is ObjectId, otherwise fail
	if !bson.IsObjectIdHex(registryId) {
		err = errors.InvalidObjectId{registryId}
		return err
	}

	// Get registry information specified by registryId parameter.
	registry := Registry{}
	query := bson.M{"_id": bson.ObjectIdHex(registryId)}
	err = getCollection(session, DB_NAME, REGISTRY_COLLECTION).Find(query).One(&registry)
	if err != nil {
		return ConvertMongoError(err, registryId)
	}

	// Store newly added docker image information.
	newImages := make([]map[string]interface{}, 0)

	// Cleanup function in case of any error.
	cleanup := func() {
		for _, image := range newImages {
			query := bson.M{"_id": bson.ObjectIdHex(registryId)}
			update := bson.M{"$pull": bson.M{"images": image["id"].(string)}}
			getCollection(session, DB_NAME, REGISTRY_COLLECTION).Update(query, update)
			imageExecutor.DeleteDockerImage(image["id"].(string))
		}
	}

	for _, image := range images {
		// Insert a new docker image information to 'image' collection.
		newImage, err := imageExecutor.AddDockerImage(image)
		if err != nil {
			// Delete docker images already added to database.
			cleanup()
			return err
		}
		newImages = append(newImages, newImage)

		// A newly added image id is added to the list of image ids in the target 'registry' document.
		query := bson.M{"_id": bson.ObjectIdHex(registryId)}
		update := bson.M{"$addToSet": bson.M{"images": newImage["id"].(string)}}
		err = getCollection(session, DB_NAME, REGISTRY_COLLECTION).Update(query, update)
		if err != nil {
			// Delete docker images already added to database.
			cleanup()
			return ConvertMongoError(err, registryId)
		}
	}
	return nil
}

// GetDockerImages returns all images document included in the target registry.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (Executor) GetDockerImages(registryId string) ([]map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	session, err := connect(DB_URL)
	if err != nil {
		return nil, err
	}
	defer close(session)

	// Verify id is ObjectId, otherwise fail
	if !bson.IsObjectIdHex(registryId) {
		err = errors.InvalidObjectId{registryId}
		return nil, err
	}

	// Get registry information specified by registryId parameter.
	registry := Registry{}
	query := bson.M{"_id": bson.ObjectIdHex(registryId)}
	err = getCollection(session, DB_NAME, REGISTRY_COLLECTION).Find(query).One(&registry)
	if err != nil {
		return nil, ConvertMongoError(err, registryId)
	}

	result := make([]map[string]interface{}, len(registry.Images))
	for i, imageId := range registry.Images {
		image, err := imageExecutor.GetDockerImage(imageId)
		if err != nil {
			return nil, err
		}
		// Remove unused 'id' field from map.
		delete(image, "id")
		result[i] = image
	}
	return result, nil
}

// UpdateDockerImage update status of docker image included in the target registry.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (client Executor) UpdateDockerImage(registryId string, image map[string]interface{}) error {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	session, err := connect(DB_URL)
	if err != nil {
		return err
	}
	defer close(session)

	// Verify id is ObjectId, otherwise fail
	if !bson.IsObjectIdHex(registryId) {
		err = errors.InvalidObjectId{registryId}
		return err
	}

	// Get registry information specified by registryId parameter.
	registry := Registry{}
	query := bson.M{"_id": bson.ObjectIdHex(registryId)}
	err = getCollection(session, DB_NAME, REGISTRY_COLLECTION).Find(query).One(&registry)
	if err != nil {
		return ConvertMongoError(err, registryId)
	}

	for _, imageId := range registry.Images {
		targetImage, err := imageExecutor.GetDockerImage(imageId)
		if err != nil {
			return err
		}

		// If the repository name specified by image parameter is the same as targetImage,
		// update status of docker image from 'image' collection.
		if targetImage["repository"].(string) == image["repository"].(string) {
			err = imageExecutor.UpdateDockerImage(imageId, image)
			if err != nil {
				return err
			}
			return nil
		}
	}
	return errors.NotFound{Message: "there is no " + image["repository"].(string)}
}

// DeleteDockerImage delete a single docker image document from 'image' collection.
// and remove it from the list of image ids in the target 'registry' document.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (Executor) DeleteDockerImage(registryId string, image map[string]interface{}) error {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	session, err := connect(DB_URL)
	if err != nil {
		return err
	}
	defer close(session)

	// Verify id is ObjectId, otherwise fail
	if !bson.IsObjectIdHex(registryId) {
		err = errors.InvalidObjectId{registryId}
		return err
	}

	// Get registry information specified by registryId parameter.
	registry := Registry{}
	query := bson.M{"_id": bson.ObjectIdHex(registryId)}
	err = getCollection(session, DB_NAME, REGISTRY_COLLECTION).Find(query).One(&registry)
	if err != nil {
		return ConvertMongoError(err, registryId)
	}

	for _, imageId := range registry.Images {
		targetImage, err := imageExecutor.GetDockerImage(imageId)
		if err != nil {
			return err
		}

		// If the repository name specified by image parameter is the same as targetImage,
		// delete it from 'image' collection and the list of image ids in the target 'registry' document.
		if targetImage["repository"].(string) == image["repository"].(string) {
			query := bson.M{"_id": bson.ObjectIdHex(registryId)}
			update := bson.M{"$pull": bson.M{"images": imageId}}
			err = getCollection(session, DB_NAME, REGISTRY_COLLECTION).Update(query, update)
			if err != nil {
				return ConvertMongoError(err)
			}

			err := imageExecutor.DeleteDockerImage(imageId)
			if err != nil {
				// Restore already deleted docker image id.
				query := bson.M{"_id": bson.ObjectIdHex(registryId)}
				update := bson.M{"$addToSet": bson.M{"images": imageId}}
				getCollection(session, DB_NAME, REGISTRY_COLLECTION).Update(query, update)
				return ConvertMongoError(err)
			}
			return nil
		}
	}
	return errors.NotFound{Message: "there is no " + image["repository"].(string)}
}
