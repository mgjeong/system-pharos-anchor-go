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
	. "db/mongo/wrapper"
	"gopkg.in/mgo.v2/bson"
)

type Command interface {
	// AddDockerRegistry insert a new docker registry information.
	AddDockerRegistry(url string) (map[string]interface{}, error)

	// GetDockerRegistries returns all documents from db related to docker registry.
	GetDockerRegistries() ([]map[string]interface{}, error)

	// DeleteDockerRegistry delete a specific docker registry information from db related to registry.
	DeleteDockerRegistry(registryId string) error
}

const (
	DB_NAME             = "DeploymentManagerDB"
	REGISTRY_COLLECTION = "REGISTRY"
	DB_URL              = "127.0.0.1:27017"
)

type Registry struct {
	ID  bson.ObjectId `bson:"_id,omitempty"`
	Url string
}

type Executor struct{}

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

// convertToMap converts Registry object into a map.
func (registry Registry) convertToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":  registry.ID.Hex(),
		"url": registry.Url,
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
	return nil
}
