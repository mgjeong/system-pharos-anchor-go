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

package subscriber

import (
	"commons/errors"
	"commons/logger"
	. "db/mongo/wrapper"
	"gopkg.in/mgo.v2/bson"
)

type Command interface {
	// AddSubscriber insert new Subscriber.
	AddSubscriber(id string, eventType string, URL string, Status []string, eventId []string) (map[string]interface{}, error)
	GetSubscriber(id string) (map[string]interface{}, error)
	DeleteSubscriber(id string) error
}

const (
	DB_NAME               = "DeploymentManagerDB"
	SUBSCRIBER_COLLECTION = "SUBSCRIBER"
	DB_URL                = "127.0.0.1:27017"
)

type Subscriber struct {
	ID      string `bson:"_id,omitempty"`
	Type    string
	URL     string
	Status  []string
	EventId []string
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

// convertToMap converts Subscriber object into a map.
func (subscriber Subscriber) convertToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":      subscriber.ID,
		"type":    subscriber.Type,
		"url":     subscriber.URL,
		"status":  subscriber.Status,
		"eventId": subscriber.EventId,
	}
}

func (Executor) AddSubscriber(id string, eventType string, url string, status []string,
	eventId []string) (map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	session, err := connect(DB_URL)
	if err != nil {
		return nil, err
	}
	defer close(session)

	subscriber := Subscriber{
		ID:      id,
		Type:    eventType,
		URL:     url,
		Status:  status,
		EventId: eventId,
	}

	err = getCollection(session, DB_NAME, SUBSCRIBER_COLLECTION).Insert(subscriber)

	if err != nil {
		return nil, ConvertMongoError(err)
	}

	result := subscriber.convertToMap()
	return result, err
}

func (Executor) GetSubscriber(id string) (map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	session, err := connect(DB_URL)
	if err != nil {
		return nil, err
	}
	defer close(session)

	subscriber := Subscriber{}
	query := bson.M{"_id": id}
	err = getCollection(session, DB_NAME, SUBSCRIBER_COLLECTION).Find(query).One(&subscriber)
	if err != nil {
		return nil, ConvertMongoError(err, id)
	}

	result := subscriber.convertToMap()
	return result, err
}

func (Executor) DeleteSubscriber(id string) error {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	if len(id) == 0 {
		err := errors.InvalidParam{"Invalid param error : subscriberId is empty."}
		return err
	}

	session, err := connect(DB_URL)
	if err != nil {
		return err
	}
	defer close(session)
	
	err = getCollection(session, DB_NAME, SUBSCRIBER_COLLECTION).Remove(bson.M{"_id": id})
	if err != nil {
		errMsg := "Failed to remove a subscriber by " + id
		return ConvertMongoError(err, errMsg)
	}

	return err
}
