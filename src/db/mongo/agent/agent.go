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
package agent

import (
	"commons/errors"
	"commons/logger"
	. "db/mongo/wrapper"
	"gopkg.in/mgo.v2/bson"
)

type Command interface {
	// AddAgent insert new Agent.
	AddAgent(ip string, status string, config map[string]interface{}) (map[string]interface{}, error)

	// UpdateAgentAddress updates ip,port of agent from db related to agent.
	UpdateAgentAddress(agent_id string, host string, port string) error

	// UpdateAgentStatus updates status of agent from db related to agent.
	UpdateAgentStatus(agent_id string, status string) error

	// GetAgent returns single document from db related to agent.
	GetAgent(agent_id string) (map[string]interface{}, error)

	// GetAllAgents returns all documents from db related to agent.
	GetAllAgents() ([]map[string]interface{}, error)

	// GetAgentByAppID returns single document including specific app.
	GetAgentByAppID(agent_id string, app_id string) (map[string]interface{}, error)

	// AddAppToAgent add specific app to the target agent.
	AddAppToAgent(agent_id string, app_id string) error

	// DeleteAppFromAgent delete specific app from the target agent.
	DeleteAppFromAgent(agent_id string, app_id string) error

	// DeleteAgent delete single document from db related to agent.
	DeleteAgent(agent_id string) error
}

const (
	DB_NAME          = "DeploymentManagerDB"
	AGENT_COLLECTION = "AGENT"
	DB_URL           = "127.0.0.1:27017"
)

type Agent struct {
	ID     bson.ObjectId `bson:"_id,omitempty"`
	IP     string
	Apps   []string
	Status string
	Config map[string]interface{}
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

// convertToMap converts Agent object into a map.
func (agent Agent) convertToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":     agent.ID.Hex(),
		"ip":     agent.IP,
		"apps":   agent.Apps,
		"status": agent.Status,
		"config": agent.Config,
	}
}

// AddAgent inserts new agent to 'agent' collection.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (Executor) AddAgent(ip string, status string, config map[string]interface{}) (map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	session, err := connect(DB_URL)
	if err != nil {
		return nil, err
	}
	defer close(session)

	agent := Agent{
		ID:     bson.NewObjectId(),
		IP:     ip,
		Status: status,
		Config: config,
	}

	err = getCollection(session, DB_NAME, AGENT_COLLECTION).Insert(agent)

	if err != nil {
		return nil, ConvertMongoError(err)
	}

	result := agent.convertToMap()
	return result, err
}

// UpdateAgentAddress updates ip,port of agent specified by agent_id parameter.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (Executor) UpdateAgentAddress(agent_id string, host string, port string) error {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	session, err := connect(DB_URL)
	if err != nil {
		return err
	}
	defer close(session)

	// Verify id is ObjectId, otherwise fail
	if !bson.IsObjectIdHex(agent_id) {
		err := errors.InvalidObjectId{agent_id}
		return err
	}

	query := bson.M{"_id": bson.ObjectIdHex(agent_id)}
	update := bson.M{"$set": bson.M{"host": host, "port": port}}
	err = getCollection(session, DB_NAME, AGENT_COLLECTION).Update(query, update)
	if err != nil {
		return ConvertMongoError(err, "Failed to update address")
	}
	return err
}

// UpdateAgentStatus updates status of agent specified by agent_id parameter.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (Executor) UpdateAgentStatus(agent_id string, status string) error {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	session, err := connect(DB_URL)
	if err != nil {
		return err
	}
	defer close(session)

	// Verify id is ObjectId, otherwise fail
	if !bson.IsObjectIdHex(agent_id) {
		err = errors.InvalidObjectId{agent_id}
		return err
	}

	query := bson.M{"_id": bson.ObjectIdHex(agent_id)}
	update := bson.M{"$set": bson.M{"status": status}}
	err = getCollection(session, DB_NAME, AGENT_COLLECTION).Update(query, update)
	if err != nil {
		return ConvertMongoError(err, "Failed to update status")
	}
	return err
}

// GetAgent returns single document specified by agent_id parameter.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (Executor) GetAgent(agent_id string) (map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	session, err := connect(DB_URL)
	if err != nil {
		return nil, err
	}
	defer close(session)

	// Verify id is ObjectId, otherwise fail
	if !bson.IsObjectIdHex(agent_id) {
		err := errors.InvalidObjectId{agent_id}
		return nil, err
	}

	agent := Agent{}
	query := bson.M{"_id": bson.ObjectIdHex(agent_id)}
	err = getCollection(session, DB_NAME, AGENT_COLLECTION).Find(query).One(&agent)
	if err != nil {
		return nil, ConvertMongoError(err, agent_id)
	}

	result := agent.convertToMap()
	return result, err
}

// GetAllAgents returns all documents from 'agent' collection.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (Executor) GetAllAgents() ([]map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	session, err := connect(DB_URL)
	if err != nil {
		return nil, err
	}
	defer close(session)

	agents := []Agent{}
	err = getCollection(session, DB_NAME, AGENT_COLLECTION).Find(nil).All(&agents)
	if err != nil {
		return nil, ConvertMongoError(err)
	}

	result := make([]map[string]interface{}, len(agents))
	for i, agent := range agents {
		result[i] = agent.convertToMap()
	}
	return result, err
}

// GetAgentByAppID returns single document specified by agent_id parameter.
// If successful, this function returns an error as nil.
// But if the target agent does not include the given app_id,
// an appropriate error will be returned.
func (Executor) GetAgentByAppID(agent_id string, app_id string) (map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	session, err := connect(DB_URL)
	if err != nil {
		return nil, err
	}
	defer close(session)

	// Verify id is ObjectId, otherwise fail
	if !bson.IsObjectIdHex(agent_id) {
		err = errors.InvalidObjectId{agent_id}
		return nil, err
	}

	agent := Agent{}
	query := bson.M{"_id": bson.ObjectIdHex(agent_id), "apps": bson.M{"$in": []string{app_id}}}
	err = getCollection(session, DB_NAME, AGENT_COLLECTION).Find(query).One(&agent)
	if err != nil {
		return nil, ConvertMongoError(err, agent_id)
	}

	result := agent.convertToMap()
	return result, err
}

// AddAppToAgent adds the specific app to the target agent.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (Executor) AddAppToAgent(agent_id string, app_id string) error {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	session, err := connect(DB_URL)
	if err != nil {
		return err
	}
	defer close(session)

	// Verify id is ObjectId, otherwise fail
	if !bson.IsObjectIdHex(agent_id) {
		err := errors.InvalidObjectId{agent_id}
		return err
	}

	query := bson.M{"_id": bson.ObjectIdHex(agent_id)}
	update := bson.M{"$addToSet": bson.M{"apps": app_id}}
	err = getCollection(session, DB_NAME, AGENT_COLLECTION).Update(query, update)
	if err != nil {
		return ConvertMongoError(err, agent_id)
	}
	return err
}

// DeleteAppFromAgent deletes the specific app from the target agent.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (Executor) DeleteAppFromAgent(agent_id string, app_id string) error {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	session, err := connect(DB_URL)
	if err != nil {
		return err
	}
	defer close(session)

	// Verify id is ObjectId, otherwise fail
	if !bson.IsObjectIdHex(agent_id) {
		err = errors.InvalidObjectId{agent_id}
		return err
	}

	query := bson.M{"_id": bson.ObjectIdHex(agent_id)}
	update := bson.M{"$pull": bson.M{"apps": app_id}}
	err = getCollection(session, DB_NAME, AGENT_COLLECTION).Update(query, update)
	if err != nil {
		return ConvertMongoError(err, agent_id)
	}
	return err
}

// DeleteAgent deletes single document from 'agent' collection.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (Executor) DeleteAgent(agent_id string) error {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	session, err := connect(DB_URL)
	if err != nil {
		return err
	}
	defer close(session)

	// Verify id is ObjectId, otherwise fail
	if !bson.IsObjectIdHex(agent_id) {
		err = errors.InvalidObjectId{agent_id}
		return err
	}

	query := bson.M{"_id": bson.ObjectIdHex(agent_id)}
	err = getCollection(session, DB_NAME, AGENT_COLLECTION).Remove(query)
	if err != nil {
		return ConvertMongoError(err, agent_id)
	}
	return err
}
