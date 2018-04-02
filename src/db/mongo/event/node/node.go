/*******************************************************************************
 * Copyright 2018 Samsung Electronics All Rights Reserved.
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
	"strings"
)

type Command interface {
	AddEvent(id string, subscriberId string) error
	GetEvent(id string) (map[string]interface{}, error)
	DeleteEvent(id string) error
	UnRegisterEvent(id string, subscriberId string) error
}

const (
	DB_NAME               = "DeploymentManagerDB"
	NODE_EVENT_COLLECTION = "NODE_EVENT"
	DB_URL                = "127.0.0.1:27017"
)

type NodeEvent struct {
	ID         string `bson:"_id,omitempty"`
	Subscriber []string
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

// Convert to map by object of struct NodeEvent.
// will return App information as map.
func (event NodeEvent) convertToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":         event.ID,
		"subscriber": event.Subscriber,
	}
}

func (Executor) AddEvent(eventId string, subscriberId string) error {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	if len(eventId) == 0 {
		err := errors.InvalidParam{"Invalid param error : eventId is empty."}
		return err
	}

	session, err := connect(DB_URL)
	if err != nil {
		return err
	}
	defer close(session)

	nodeEvent := NodeEvent{}
	query := bson.M{"_id": eventId}

	// Check whether the corresponding EventId (that is, the Event having the same Option)
	// exists in NodeEventDB, add it to the Subscriber list if it exists,
	// and create an Event if it does not exist.
	err = getCollection(session, DB_NAME, NODE_EVENT_COLLECTION).Find(query).One(&nodeEvent)
	if err != nil {
		err = ConvertMongoError(err)
		switch err.(type) {
		default:
			return err
		case errors.NotFound:
			subscriber := make([]string, 0)
			subscriber = append(subscriber, subscriberId)
			nodeEvent = NodeEvent{
				ID:         eventId,
				Subscriber: subscriber,
			}

			err = getCollection(session, DB_NAME, NODE_EVENT_COLLECTION).Insert(nodeEvent)
			if err != nil {
				return ConvertMongoError(err, "")
			}
			return nil
		}
	} else {
		for _, subsId := range nodeEvent.Subscriber {
			if strings.Compare(subsId, subscriberId) == 0 {
				return nil
			}
		}
		nodeEvent.Subscriber = append(nodeEvent.Subscriber, subscriberId)
		err = getCollection(session, DB_NAME, NODE_EVENT_COLLECTION).Update(query, nodeEvent)
		if err != nil {
			return ConvertMongoError(err, "")
		}
		return nil
	}
	return err
}

func (Executor) GetEvent(id string) (map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	session, err := connect(DB_URL)
	if err != nil {
		return nil, err
	}
	defer close(session)

	nodeEvent := NodeEvent{}
	query := bson.M{"_id": id}
	err = getCollection(session, DB_NAME, NODE_EVENT_COLLECTION).Find(query).One(&nodeEvent)
	if err != nil {
		return nil, ConvertMongoError(err, id)
	}

	result := nodeEvent.convertToMap()
	return result, err
}

func (Executor) DeleteEvent(id string) error {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	if len(id) == 0 {
		err := errors.InvalidParam{"Invalid param error : nodeEventId is empty."}
		return err
	}

	session, err := connect(DB_URL)
	if err != nil {
		return err
	}
	defer close(session)

	err = getCollection(session, DB_NAME, NODE_EVENT_COLLECTION).Remove(bson.M{"_id": id})
	if err != nil {
		errMsg := "Failed to remove a nodeEvent by " + id
		return ConvertMongoError(err, errMsg)
	}

	return err
}

func (Executor) UnRegisterEvent(id string, subscriberId string) error {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	session, err := connect(DB_URL)
	if err != nil {
		return err
	}
	defer close(session)

	nodeEvent := NodeEvent{}
	query := bson.M{"_id": id}
	err = getCollection(session, DB_NAME, NODE_EVENT_COLLECTION).Find(query).One(&nodeEvent)
	if err != nil {
		return ConvertMongoError(err, id)
	}

	for i, subs := range nodeEvent.Subscriber {
		if subs == subscriberId {
			nodeEvent.Subscriber = append(nodeEvent.Subscriber[:i], nodeEvent.Subscriber[i+1:]...)
			break
		}
	}
	err = getCollection(session, DB_NAME, NODE_EVENT_COLLECTION).Update(query, nodeEvent)
	if err != nil {
		return ConvertMongoError(err, id)
	}

	return err
}
