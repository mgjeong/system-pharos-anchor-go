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
	nodeDB "db/mongo/node"
	. "db/mongo/wrapper"

	"gopkg.in/mgo.v2/bson"
)

type Command interface {
	// CreateGroup insert new Group.
	CreateGroup(name string) (map[string]interface{}, error)

	// GetGroup returns single document from db related to group.
	GetGroup(groupId string) (map[string]interface{}, error)

	// GetGroups returns all documents from db related to group.
	GetGroups() ([]map[string]interface{}, error)

	// GetGroupMembers returns all nodes who belong to the target group.
	GetGroupMembers(groupId string) ([]map[string]interface{}, error)

	// GetGroupMembersByAppID returns all nodes including specific app on the target group.
	GetGroupMembersByAppID(groupId string, appId string) ([]map[string]interface{}, error)

	// JoinGroup add specific node to the target group.
	JoinGroup(groupId string, nodeId string) error

	// LeaveGroup delete specific node from the target group.
	LeaveGroup(groupId string, nodeId string) error

	// DeleteGroup delete single document from db related to group.
	DeleteGroup(groupId string) error
}

const (
	DB_NAME          = "DeploymentManagerDB"
	GROUP_COLLECTION = "GROUP"
	DB_URL           = "127.0.0.1:27017"
)

type Group struct {
	ID      bson.ObjectId `bson:"_id,omitempty"`
	Name    string
	Members []string
}

type Executor struct{}

var mgoDial Connection
var nodeExecutor nodeDB.Command

func init() {
	mgoDial = MongoDial{}
	nodeExecutor = nodeDB.Executor{}
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
		"name":    group.Name,
		"members": group.Members,
	}
}

// CreateGroup inserts new Group to 'group' collection.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (Executor) CreateGroup(name string) (map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	session, err := connect(DB_URL)
	if err != nil {
		return nil, err
	}
	defer close(session)

	group := Group{
		ID:      bson.NewObjectId(),
		Name:    name,
		Members: []string{},
	}

	err = getCollection(session, DB_NAME, GROUP_COLLECTION).Insert(group)
	if err != nil {
		return nil, ConvertMongoError(err)
	}

	result := group.convertToMap()
	return result, err
}

// GetGroup returns single document specified by groupId parameter.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (Executor) GetGroup(groupId string) (map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	session, err := connect(DB_URL)
	if err != nil {
		return nil, err
	}
	defer close(session)

	// Verify id is ObjectId, otherwise fail
	if !bson.IsObjectIdHex(groupId) {
		err = errors.InvalidObjectId{groupId}
		return nil, err
	}

	group := Group{}
	query := bson.M{"_id": bson.ObjectIdHex(groupId)}
	err = getCollection(session, DB_NAME, GROUP_COLLECTION).Find(query).One(&group)
	if err != nil {
		return nil, ConvertMongoError(err, groupId)
	}

	result := group.convertToMap()
	return result, err
}

// GetGroups returns all documents from 'group' collection.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (Executor) GetGroups() ([]map[string]interface{}, error) {
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

// JoinGroup adds the specific node to a list of group members.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (Executor) JoinGroup(groupId string, nodeId string) error {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	session, err := connect(DB_URL)
	if err != nil {
		return err
	}
	defer close(session)

	// Verify id is ObjectId, otherwise fail
	if !bson.IsObjectIdHex(groupId) {
		err := errors.InvalidObjectId{groupId}
		return err
	}
	if !bson.IsObjectIdHex(nodeId) {
		err := errors.InvalidObjectId{nodeId}
		return err
	}

	query := bson.M{"_id": bson.ObjectIdHex(groupId)}
	update := bson.M{"$addToSet": bson.M{"members": nodeId}}
	err = getCollection(session, DB_NAME, GROUP_COLLECTION).Update(query, update)
	if err != nil {
		return ConvertMongoError(err, groupId)
	}
	return err
}

// LeaveGroup deletes the specific node from a list of group members.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (Executor) LeaveGroup(groupId string, nodeId string) error {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	session, err := connect(DB_URL)
	if err != nil {
		return err
	}
	defer close(session)

	// Verify id is ObjectId, otherwise fail
	if !bson.IsObjectIdHex(groupId) {
		err = errors.InvalidObjectId{groupId}
		return err
	}
	if !bson.IsObjectIdHex(nodeId) {
		err = errors.InvalidObjectId{nodeId}
		return err
	}

	query := bson.M{"_id": bson.ObjectIdHex(groupId)}
	update := bson.M{"$pull": bson.M{"members": nodeId}}
	err = getCollection(session, DB_NAME, GROUP_COLLECTION).Update(query, update)
	if err != nil {
		return ConvertMongoError(err, groupId)
	}
	return err
}

// GetGroupMembers returns all nodes who belong to the target group.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (client Executor) GetGroupMembers(groupId string) ([]map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Verify id is ObjectId, otherwise fail
	if !bson.IsObjectIdHex(groupId) {
		err := errors.InvalidObjectId{groupId}
		return nil, err
	}

	group, err := client.GetGroup(groupId)
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, len(group["members"].([]string)))
	for i, nodeId := range group["members"].([]string) {
		var node map[string]interface{}
		node, err := nodeExecutor.GetNode(nodeId)
		if err != nil {
			return nil, err
		}
		result[i] = node
	}
	return result, err
}

// GetGroupMembersByAppID returns all nodes including the app identified
// by the given appid on the target group.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (client Executor) GetGroupMembersByAppID(groupId string, appId string) ([]map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Verify id is ObjectId, otherwise fail
	if !bson.IsObjectIdHex(groupId) {
		err := errors.InvalidObjectId{groupId}
		return nil, err
	}

	group, err := client.GetGroup(groupId)
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, len(group["members"].([]string)))
	for i, nodeId := range group["members"].([]string) {
		var node map[string]interface{}
		node, err := nodeExecutor.GetNodeByAppID(nodeId, appId)
		if err != nil {
			return nil, err
		}
		result[i] = node
	}

	return result, err
}

// DeleteGroup deletes single document specified by groupId parameter.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (Executor) DeleteGroup(groupId string) error {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	session, err := connect(DB_URL)
	if err != nil {
		return err
	}
	defer close(session)

	// Verify id is ObjectId, otherwise fail
	if !bson.IsObjectIdHex(groupId) {
		err = errors.InvalidObjectId{groupId}
		return err
	}

	query := bson.M{"_id": bson.ObjectIdHex(groupId)}
	err = getCollection(session, DB_NAME, GROUP_COLLECTION).Remove(query)
	if err != nil {
		return ConvertMongoError(err, groupId)
	}
	return err
}
