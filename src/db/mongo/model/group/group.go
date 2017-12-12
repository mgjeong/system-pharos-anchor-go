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
package group

import (
	"commons/errors"
	"commons/logger"
	. "db/modelinterface"
	agentDB "db/mongo/model/agent"
	. "db/mongo/wrapper"

	"gopkg.in/mgo.v2/bson"
)

const (
	DB_NAME          = "DeploymentManagerDB"
	GROUP_COLLECTION = "GROUP"
	DB_URL           = "127.0.0.1:27017"
)

type Group struct {
	ID      bson.ObjectId `bson:"_id,omitempty"`
	Members []string
}

type DBManager struct {
	GroupInterface
}

var mgoDial Connection
var agentDBManager AgentInterface

func init() {
	mgoDial = MongoDial{}
	agentDBManager = agentDB.DBManager{}
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

// convertToMap converts Group object into a map.
func (group Group) convertToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":      group.ID.Hex(),
		"members": group.Members,
	}
}

// CreateGroup inserts new Group to 'group' collection.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (DBManager) CreateGroup() (map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	session, err := connect(DB_URL)
	if err != nil {
		return nil, err
	}
	defer close(session)

	group := Group{
		ID: bson.NewObjectId(),
	}

	err = getCollection(session, DB_NAME, GROUP_COLLECTION).Insert(group)
	if err != nil {
		return nil, ConvertMongoError(err)
	}

	result := group.convertToMap()
	return result, err
}

// GetGroup returns single document specified by group_id parameter.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (DBManager) GetGroup(group_id string) (map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	session, err := connect(DB_URL)
	if err != nil {
		return nil, err
	}
	defer close(session)

	// Verify id is ObjectId, otherwise fail
	if !bson.IsObjectIdHex(group_id) {
		err = errors.InvalidObjectId{group_id}
		return nil, err
	}

	group := Group{}
	query := bson.M{"_id": bson.ObjectIdHex(group_id)}
	err = getCollection(session, DB_NAME, GROUP_COLLECTION).Find(query).One(&group)
	if err != nil {
		return nil, ConvertMongoError(err, group_id)
	}

	result := group.convertToMap()
	return result, err
}

// GetAllGroups returns all documents from 'group' collection.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (DBManager) GetAllGroups() ([]map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	session, err := connect(DB_URL)
	if err != nil {
		return nil, err
	}
	defer close(session)

	groups := []Group{}
	err = getCollection(session, DB_NAME, GROUP_COLLECTION).Find(nil).All(&groups)
	if err != nil {
		return nil, ConvertMongoError(err)
	}

	result := make([]map[string]interface{}, len(groups))
	for i, group := range groups {
		result[i] = group.convertToMap()
	}
	return result, err
}

// JoinGroup adds the specific agent to a list of group members.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (DBManager) JoinGroup(group_id string, agent_id string) error {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	session, err := connect(DB_URL)
	if err != nil {
		return err
	}
	defer close(session)

	// Verify id is ObjectId, otherwise fail
	if !bson.IsObjectIdHex(group_id) {
		err := errors.InvalidObjectId{group_id}
		return err
	}
	if !bson.IsObjectIdHex(agent_id) {
		err := errors.InvalidObjectId{agent_id}
		return err
	}

	query := bson.M{"_id": bson.ObjectIdHex(group_id)}
	update := bson.M{"$addToSet": bson.M{"members": agent_id}}
	err = getCollection(session, DB_NAME, GROUP_COLLECTION).Update(query, update)
	if err != nil {
		return ConvertMongoError(err, group_id)
	}
	return err
}

// LeaveGroup deletes the specific agent from a list of group members.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (DBManager) LeaveGroup(group_id string, agent_id string) error {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	session, err := connect(DB_URL)
	if err != nil {
		return err
	}
	defer close(session)

	// Verify id is ObjectId, otherwise fail
	if !bson.IsObjectIdHex(group_id) {
		err = errors.InvalidObjectId{group_id}
		return err
	}
	if !bson.IsObjectIdHex(agent_id) {
		err = errors.InvalidObjectId{agent_id}
		return err
	}

	query := bson.M{"_id": bson.ObjectIdHex(group_id)}
	update := bson.M{"$pull": bson.M{"members": agent_id}}
	err = getCollection(session, DB_NAME, GROUP_COLLECTION).Update(query, update)
	if err != nil {
		return ConvertMongoError(err, group_id)
	}
	return err
}

// GetGroupMembers returns all agents who belong to the target group.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (client *DBManager) GetGroupMembers(group_id string) ([]map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Verify id is ObjectId, otherwise fail
	if !bson.IsObjectIdHex(group_id) {
		err := errors.InvalidObjectId{group_id}
		return nil, err
	}

	group, err := client.GetGroup(group_id)
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, len(group["members"].([]string)))
	for i, agent_id := range group["members"].([]string) {
		var agent map[string]interface{}
		agent, err := agentDBManager.GetAgent(agent_id)
		if err != nil {
			return nil, err
		}
		result[i] = agent
	}
	return result, err
}

// GetGroupMembersByAppID returns all agents including the app identified
// by the given appid on the target group.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (client *DBManager) GetGroupMembersByAppID(group_id string, app_id string) ([]map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Verify id is ObjectId, otherwise fail
	if !bson.IsObjectIdHex(group_id) {
		err := errors.InvalidObjectId{group_id}
		return nil, err
	}

	group, err := client.GetGroup(group_id)
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, len(group["members"].([]string)))
	for i, agent_id := range group["members"].([]string) {
		var agent map[string]interface{}
		agent, err := agentDBManager.GetAgentByAppID(agent_id, app_id)
		if err != nil {
			return nil, err
		}
		result[i] = agent
	}

	return result, err
}

// DeleteGroup deletes single document specified by group_id parameter.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (DBManager) DeleteGroup(group_id string) error {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	session, err := connect(DB_URL)
	if err != nil {
		return err
	}
	defer close(session)

	// Verify id is ObjectId, otherwise fail
	if !bson.IsObjectIdHex(group_id) {
		err = errors.InvalidObjectId{group_id}
		return err
	}

	query := bson.M{"_id": bson.ObjectIdHex(group_id)}
	err = getCollection(session, DB_NAME, GROUP_COLLECTION).Remove(query)
	if err != nil {
		return ConvertMongoError(err, group_id)
	}
	return err
}