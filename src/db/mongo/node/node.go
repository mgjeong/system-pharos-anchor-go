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
package node

import (
	"commons/errors"
	"commons/logger"
	. "db/mongo/wrapper"
	"gopkg.in/mgo.v2/bson"
)

type Command interface {
	// AddNode insert new Node.
	AddNode(ip string, status string, config map[string]interface{}) (map[string]interface{}, error)

	// UpdateNodeAddress updates ip,port of node from db related to node.
	UpdateNodeAddress(nodeId string, host string, port string) error

	// UpdateNodeStatus updates status of node from db related to node.
	UpdateNodeStatus(nodeId string, status string) error

	// UpdateNodeConfiguration updates configuration information of node from db related to node.
	UpdateNodeConfiguration(nodeId string, config map[string]interface{}) error

	// GetNode returns single document from db related to node.
	GetNode(nodeId string) (map[string]interface{}, error)

	// GetNodes returns all matches for the query-string which is passed in call to function.
	GetNodes(queryOptional ...map[string]interface{}) ([]map[string]interface{}, error)

	// GetNodeByAppID returns single document including specific app.
	GetNodeByAppID(nodeId string, appId string) (map[string]interface{}, error)

	// GetNodeByIP returns single document from db related to node.
	GetNodeByIP(ip string) (map[string]interface{}, error)

	// AddAppToNode add specific app to the target node.
	AddAppToNode(nodeId string, appId string) error

	// DeleteAppFromNode delete specific app from the target node.
	DeleteAppFromNode(nodeId string, appId string) error

	// DeleteNode delete single document from db related to node.
	DeleteNode(nodeId string) error
}

const (
	DB_NAME         = "DeploymentManagerDB"
	NODE_COLLECTION = "NODE"
	DB_URL          = "127.0.0.1:27017"
)

type Node struct {
	ID     bson.ObjectId `bson:"_id,omitempty"`
	IP     string
	Apps   []string
	Status string
	Config map[string]interface{}
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

// convertToMap converts Node object into a map.
func (node Node) convertToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":     node.ID.Hex(),
		"ip":     node.IP,
		"apps":   node.Apps,
		"status": node.Status,
		"config": node.Config,
	}
}

// AddNode inserts new node to 'node' collection.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (Executor) AddNode(ip string, status string, config map[string]interface{}) (map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	session, err := connect(DB_URL)
	if err != nil {
		return nil, err
	}
	defer close(session)

	node := Node{
		ID:     bson.NewObjectId(),
		IP:     ip,
		Status: status,
		Config: config,
	}

	err = getCollection(session, DB_NAME, NODE_COLLECTION).Insert(node)

	if err != nil {
		return nil, ConvertMongoError(err)
	}

	result := node.convertToMap()
	return result, err
}

// UpdateNodeAddress updates ip,port of node specified by nodeId parameter.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (Executor) UpdateNodeAddress(nodeId string, host string, port string) error {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	session, err := connect(DB_URL)
	if err != nil {
		return err
	}
	defer close(session)

	// Verify id is ObjectId, otherwise fail
	if !bson.IsObjectIdHex(nodeId) {
		err := errors.InvalidObjectId{nodeId}
		return err
	}

	query := bson.M{"_id": bson.ObjectIdHex(nodeId)}
	update := bson.M{"$set": bson.M{"host": host, "port": port}}
	err = getCollection(session, DB_NAME, NODE_COLLECTION).Update(query, update)
	if err != nil {
		return ConvertMongoError(err, "Failed to update address")
	}
	return err
}

// UpdateNodeStatus updates status of node specified by nodeId parameter.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (Executor) UpdateNodeStatus(nodeId string, status string) error {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	session, err := connect(DB_URL)
	if err != nil {
		return err
	}
	defer close(session)

	// Verify id is ObjectId, otherwise fail
	if !bson.IsObjectIdHex(nodeId) {
		err = errors.InvalidObjectId{nodeId}
		return err
	}

	query := bson.M{"_id": bson.ObjectIdHex(nodeId)}
	update := bson.M{"$set": bson.M{"status": status}}
	err = getCollection(session, DB_NAME, NODE_COLLECTION).Update(query, update)
	if err != nil {
		return ConvertMongoError(err, "Failed to update status")
	}
	return err
}

func (Executor) UpdateNodeConfiguration(nodeId string, config map[string]interface{}) error {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	session, err := connect(DB_URL)
	if err != nil {
		return err
	}
	defer close(session)

	// Verify id is ObjectId, otherwise fail
	if !bson.IsObjectIdHex(nodeId) {
		err = errors.InvalidObjectId{nodeId}
		return err
	}

	query := bson.M{"_id": bson.ObjectIdHex(nodeId)}
	update := bson.M{"$set": bson.M{"config": config}}
	err = getCollection(session, DB_NAME, NODE_COLLECTION).Update(query, update)
	if err != nil {
		return ConvertMongoError(err, "Failed to update status")
	}
	return err
}

// GetNode returns single document specified by nodeId parameter.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (Executor) GetNode(nodeId string) (map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	session, err := connect(DB_URL)
	if err != nil {
		return nil, err
	}
	defer close(session)

	// Verify id is ObjectId, otherwise fail
	if !bson.IsObjectIdHex(nodeId) {
		err := errors.InvalidObjectId{nodeId}
		return nil, err
	}

	node := Node{}
	query := bson.M{"_id": bson.ObjectIdHex(nodeId)}
	err = getCollection(session, DB_NAME, NODE_COLLECTION).Find(query).One(&node)
	if err != nil {
		return nil, ConvertMongoError(err, nodeId)
	}

	result := node.convertToMap()
	return result, err
}

// GetNodes returns all documents from 'node' collection.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (Executor) GetNodes(queryOptional ...map[string]interface{}) ([]map[string]interface{}, error) {
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

	nodes := []Node{}
	err = getCollection(session, DB_NAME, NODE_COLLECTION).Find(query).All(&nodes)
	if err != nil {
		return nil, ConvertMongoError(err)
	}

	result := make([]map[string]interface{}, len(nodes))
	for i, node := range nodes {
		result[i] = node.convertToMap()
	}
	return result, err
}

// GetNodeByAppID returns single document specified by nodeId parameter.
// If successful, this function returns an error as nil.
// But if the target node does not include the given appId,
// an appropriate error will be returned.
func (Executor) GetNodeByAppID(nodeId string, appId string) (map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	session, err := connect(DB_URL)
	if err != nil {
		return nil, err
	}
	defer close(session)

	// Verify id is ObjectId, otherwise fail
	if !bson.IsObjectIdHex(nodeId) {
		err = errors.InvalidObjectId{nodeId}
		return nil, err
	}

	node := Node{}
	query := bson.M{"_id": bson.ObjectIdHex(nodeId), "apps": bson.M{"$in": []string{appId}}}
	err = getCollection(session, DB_NAME, NODE_COLLECTION).Find(query).One(&node)
	if err != nil {
		return nil, ConvertMongoError(err, nodeId)
	}

	result := node.convertToMap()
	return result, err
}

// GetNodeByIP returns single document specified by ip parameter.
// If successful, this function returns an error as nil.
// But if the target node does not include the given appId,
// an appropriate error will be returned.
func (Executor) GetNodeByIP(ip string) (map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	session, err := connect(DB_URL)
	if err != nil {
		return nil, err
	}
	defer close(session)

	node := Node{}
	query := bson.M{"ip": ip}
	err = getCollection(session, DB_NAME, NODE_COLLECTION).Find(query).One(&node)
	if err != nil {
		return nil, ConvertMongoError(err, ip)
	}

	result := node.convertToMap()
	return result, err
}

// AddAppToNode adds the specific app to the target node.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (Executor) AddAppToNode(nodeId string, appId string) error {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	session, err := connect(DB_URL)
	if err != nil {
		return err
	}
	defer close(session)

	// Verify id is ObjectId, otherwise fail
	if !bson.IsObjectIdHex(nodeId) {
		err := errors.InvalidObjectId{nodeId}
		return err
	}

	query := bson.M{"_id": bson.ObjectIdHex(nodeId)}
	update := bson.M{"$addToSet": bson.M{"apps": appId}}
	err = getCollection(session, DB_NAME, NODE_COLLECTION).Update(query, update)
	if err != nil {
		return ConvertMongoError(err, nodeId)
	}
	return err
}

// DeleteAppFromNode deletes the specific app from the target node.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (Executor) DeleteAppFromNode(nodeId string, appId string) error {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	session, err := connect(DB_URL)
	if err != nil {
		return err
	}
	defer close(session)

	// Verify id is ObjectId, otherwise fail
	if !bson.IsObjectIdHex(nodeId) {
		err = errors.InvalidObjectId{nodeId}
		return err
	}

	query := bson.M{"_id": bson.ObjectIdHex(nodeId)}
	update := bson.M{"$pull": bson.M{"apps": appId}}
	err = getCollection(session, DB_NAME, NODE_COLLECTION).Update(query, update)
	if err != nil {
		return ConvertMongoError(err, nodeId)
	}
	return err
}

// DeleteNode deletes single document from 'node' collection.
// If successful, this function returns an error as nil.
// otherwise, an appropriate error will be returned.
func (Executor) DeleteNode(nodeId string) error {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	session, err := connect(DB_URL)
	if err != nil {
		return err
	}
	defer close(session)

	// Verify id is ObjectId, otherwise fail
	if !bson.IsObjectIdHex(nodeId) {
		err = errors.InvalidObjectId{nodeId}
		return err
	}

	query := bson.M{"_id": bson.ObjectIdHex(nodeId)}
	err = getCollection(session, DB_NAME, NODE_COLLECTION).Remove(query)
	if err != nil {
		return ConvertMongoError(err, nodeId)
	}
	return err
}
