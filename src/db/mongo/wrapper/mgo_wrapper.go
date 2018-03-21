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
package wrapper

import (
	"commons/errors"
	"gopkg.in/mgo.v2"
)

type (
	Session interface {
		DB(name string) Database
		Close()
	}

	MongoSession struct {
		Session *mgo.Session
	}

	Connection interface {
		Dial(url string) (Session, error)
	}

	MongoDial struct{}

	Database interface {
		C(name string) Collection
	}

	MongoDatabase struct {
		Database *mgo.Database
	}

	Collection interface {
		Find(query interface{}) Query
		Insert(docs ...interface{}) error
		Remove(selector interface{}) error
		Update(selector interface{}, update interface{}) error
	}

	MongoCollection struct {
		Collection *mgo.Collection
	}

	Query interface {
		All(result interface{}) error
		One(result interface{}) error
	}

	MongoQuery struct {
		Query *mgo.Query
	}
)

func (s MongoSession) DB(name string) Database {
	return &MongoDatabase{Database: s.Session.DB(name)}
}

func (s MongoSession) Close() {
	s.Session.Close()
}

// Dial is a wrapper function used to abstract mgo Dial function.
func (MongoDial) Dial(url string) (Session, error) {
	session, err := mgo.Dial(url)
	return MongoSession{Session: session}, err
}

// C is a wrapper function used to abstract mgo C function.
func (d MongoDatabase) C(name string) Collection {
	return &MongoCollection{Collection: d.Database.C(name)}
}

// Find is a wrapper function used to abstract mgo Find function.
func (c MongoCollection) Find(query interface{}) Query {
	return MongoQuery{Query: c.Collection.Find(query)}
}

// Insert is a wrapper function used to abstract mgo Insert function.
func (c MongoCollection) Insert(docs ...interface{}) error {
	return c.Collection.Insert(docs...)
}

// Remove is a wrapper function used to abstract mgo Remove function.
func (c MongoCollection) Remove(selector interface{}) error {
	return c.Collection.Remove(selector)
}

// Update is a wrapper function used to abstract mgo Update function.
func (c MongoCollection) Update(selector interface{}, update interface{}) error {
	return c.Collection.Update(selector, update)
}

// All is a wrapper function used to abstract mgo All function.
func (q MongoQuery) All(result interface{}) error {
	return q.Query.All(result)
}

// One is a wrapper function used to abstract mgo One function.
func (q MongoQuery) One(result interface{}) error {
	return q.Query.One(result)
}

// ConvertMongoError converts a mongo error into an error defined in errors package.
func ConvertMongoError(mgoError error, message ...string) (err error) {
	switch mgoError {
	case mgo.ErrNotFound:
		targetId := ""
		if len(message) != 0 {
			targetId = message[0]
		}
		return errors.NotFound{targetId}
	default:
		return errors.DBOperationError{mgoError.Error()}
	}
}
