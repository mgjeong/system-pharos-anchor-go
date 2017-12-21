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
	errors "commons/errors"
	mgomocks "db/mongo/wrapper/mocks"
	"github.com/golang/mock/gomock"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"reflect"
	"testing"
)

const (
	validUrl        = "127.0.0.1:27017"
	dbName          = "DeploymentManagerDB"
	collectionName  = "IMAGE"
	imageId         = "000000000000000000000000"
	invalidObjectId = ""
	invalidUrl      = ""
)

var (
	image = map[string]interface{}{
		"repository": "http://localhost:5000",
		"tag":        "latest",
		"size":       "3453",
		"action":     "push",
		"timestamp":  "2017-11-27T12:42:56.914783506Z",
	}
	dummySession       = mgomocks.MockSession{}
	connectionError    = errors.DBConnectionError{}
	invalidObjectError = errors.InvalidObjectId{invalidObjectId}
)

func TestCalledConnectWithEmptyURL_ExpectErrorReturn(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	connectionMockObj := mgomocks.NewMockConnection(mockCtrl)

	gomock.InOrder(
		connectionMockObj.EXPECT().Dial(invalidUrl).Return(&dummySession, connectionError),
	)
	mgoDial = connectionMockObj

	_, err := connect(invalidUrl)

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "UnknownError", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "UnknownError", err.Error())
	case errors.DBOperationError:
	}
}

func TestCalledConnectWithValidURL_ExpectSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	connectionMockObj := mgomocks.NewMockConnection(mockCtrl)

	gomock.InOrder(
		connectionMockObj.EXPECT().Dial(validUrl).Return(&dummySession, nil),
	)
	mgoDial = connectionMockObj

	_, err := connect(validUrl)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}
}

func TestCalledClose_ExpectSessionClosed(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	sessionMockObj := mgomocks.NewMockSession(mockCtrl)
	gomock.InOrder(
		sessionMockObj.EXPECT().Close(),
	)

	close(sessionMockObj)
}

func TestCalled_GetCollcetion_ExpectToCCalled(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	sessionMockObj := mgomocks.NewMockSession(mockCtrl)
	dbMockObj := mgomocks.NewMockDatabase(mockCtrl)
	collectionMockObj := mgomocks.NewMockCollection(mockCtrl)

	gomock.InOrder(
		sessionMockObj.EXPECT().DB(dbName).Return(dbMockObj),
		dbMockObj.EXPECT().C(collectionName).Return(collectionMockObj),
	)

	collection := getCollection(sessionMockObj, dbName, collectionName)

	if collection == nil {
		t.Errorf("Unexpected err: getCollection returns nil")
	}
}

func TestCalledAddDockerImage_ExpectSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	connectionMockObj := mgomocks.NewMockConnection(mockCtrl)
	sessionMockObj := mgomocks.NewMockSession(mockCtrl)
	dbMockObj := mgomocks.NewMockDatabase(mockCtrl)
	collectionMockObj := mgomocks.NewMockCollection(mockCtrl)

	gomock.InOrder(
		connectionMockObj.EXPECT().Dial(validUrl).Return(sessionMockObj, nil),
		sessionMockObj.EXPECT().DB(dbName).Return(dbMockObj),
		dbMockObj.EXPECT().C(gomock.Any()).Return(collectionMockObj),
		collectionMockObj.EXPECT().Insert(gomock.Any()).Return(nil),
		sessionMockObj.EXPECT().Close(),
	)

	mgoDial = connectionMockObj
	dbManager := DBManager{}

	_, err := dbManager.AddDockerImage(image)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}
}

func TestCalledAddDockerImageWhenDBReturnsError_ExpectErrorReturn(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	connectionMockObj := mgomocks.NewMockConnection(mockCtrl)
	sessionMockObj := mgomocks.NewMockSession(mockCtrl)
	dbMockObj := mgomocks.NewMockDatabase(mockCtrl)
	collectionMockObj := mgomocks.NewMockCollection(mockCtrl)

	gomock.InOrder(
		connectionMockObj.EXPECT().Dial(validUrl).Return(sessionMockObj, nil),
		sessionMockObj.EXPECT().DB(gomock.Any()).Return(dbMockObj),
		dbMockObj.EXPECT().C(gomock.Any()).Return(collectionMockObj),
		collectionMockObj.EXPECT().Insert(gomock.Any()).Return(mgo.ErrNotFound),
		sessionMockObj.EXPECT().Close(),
	)

	mgoDial = connectionMockObj
	dbManager := DBManager{}

	_, err := dbManager.AddDockerImage(image)

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", err.Error())
	case errors.NotFound:
	}
}

func TestCalledUpdateDockerImage_ExpectSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	query := bson.M{"_id": bson.ObjectIdHex(imageId), "repository": image["repository"].(string)}
	update := bson.M{
		"$set": bson.M{
			"tag":       image["tag"].(string),
			"size":      image["size"].(string),
			"action":    image["action"].(string),
			"timestamp": image["timestamp"].(string)}}

	connectionMockObj := mgomocks.NewMockConnection(mockCtrl)
	sessionMockObj := mgomocks.NewMockSession(mockCtrl)
	dbMockObj := mgomocks.NewMockDatabase(mockCtrl)
	collectionMockObj := mgomocks.NewMockCollection(mockCtrl)

	gomock.InOrder(
		connectionMockObj.EXPECT().Dial(validUrl).Return(sessionMockObj, nil),
		sessionMockObj.EXPECT().DB(gomock.Any()).Return(dbMockObj),
		dbMockObj.EXPECT().C(gomock.Any()).Return(collectionMockObj),
		collectionMockObj.EXPECT().Update(query, update).Return(nil),
		sessionMockObj.EXPECT().Close(),
	)

	mgoDial = connectionMockObj
	dbManager := DBManager{}
	err := dbManager.UpdateDockerImage(imageId, image)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}
}

func TestCalledUpdateDockerImageWithInvalidObjectId_ExpectErrorReturn(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	connectionMockObj := mgomocks.NewMockConnection(mockCtrl)
	sessionMockObj := mgomocks.NewMockSession(mockCtrl)

	gomock.InOrder(
		connectionMockObj.EXPECT().Dial(validUrl).Return(sessionMockObj, nil),
		sessionMockObj.EXPECT().Close(),
	)

	mgoDial = connectionMockObj
	dbManager := DBManager{}

	err := dbManager.UpdateDockerImage(invalidObjectId, image)

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", invalidObjectError.Error(), "nil")
	}

	if err.Error() != invalidObjectError.Error() {
		t.Errorf("Expected err: %s, actual err: %s", invalidObjectError.Error(), err.Error())
	}
}

func TestCalledUpdateDockerImageWhenDBReturnsError_ExpectErrorReturn(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	query := bson.M{"_id": bson.ObjectIdHex(imageId), "repository": image["repository"].(string)}
	update := bson.M{
		"$set": bson.M{
			"tag":       image["tag"].(string),
			"size":      image["size"].(string),
			"action":    image["action"].(string),
			"timestamp": image["timestamp"].(string)}}

	connectionMockObj := mgomocks.NewMockConnection(mockCtrl)
	sessionMockObj := mgomocks.NewMockSession(mockCtrl)
	dbMockObj := mgomocks.NewMockDatabase(mockCtrl)
	collectionMockObj := mgomocks.NewMockCollection(mockCtrl)

	gomock.InOrder(
		connectionMockObj.EXPECT().Dial(validUrl).Return(sessionMockObj, nil),
		sessionMockObj.EXPECT().DB(gomock.Any()).Return(dbMockObj),
		dbMockObj.EXPECT().C(gomock.Any()).Return(collectionMockObj),
		collectionMockObj.EXPECT().Update(query, update).Return(mgo.ErrNotFound),
		sessionMockObj.EXPECT().Close(),
	)

	mgoDial = connectionMockObj

	dbManager := DBManager{}
	err := dbManager.UpdateDockerImage(imageId, image)

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", err.Error())
	case errors.NotFound:
	}
}

func TestCalledGetDockerImage_ExpectSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	query := bson.M{"_id": bson.ObjectIdHex(imageId)}
	arg := Image{
		ID:         bson.ObjectIdHex(imageId),
		Repository: image["repository"].(string),
		Tag:        image["tag"].(string),
		Size:       image["size"].(string),
		Action:     image["action"].(string),
		Timestamp:  image["timestamp"].(string)}

	expectedRes := map[string]interface{}{
		"id":         imageId,
		"repository": image["repository"].(string),
		"tag":        image["tag"].(string),
		"size":       image["size"].(string),
		"action":     image["action"].(string),
		"timestamp":  image["timestamp"].(string),
	}

	connectionMockObj := mgomocks.NewMockConnection(mockCtrl)
	sessionMockObj := mgomocks.NewMockSession(mockCtrl)
	dbMockObj := mgomocks.NewMockDatabase(mockCtrl)
	collectionMockObj := mgomocks.NewMockCollection(mockCtrl)
	queryMockObj := mgomocks.NewMockQuery(mockCtrl)

	gomock.InOrder(
		connectionMockObj.EXPECT().Dial(validUrl).Return(sessionMockObj, nil),
		sessionMockObj.EXPECT().DB(gomock.Any()).Return(dbMockObj),
		dbMockObj.EXPECT().C(gomock.Any()).Return(collectionMockObj),
		collectionMockObj.EXPECT().Find(query).Return(queryMockObj),
		queryMockObj.EXPECT().One(gomock.Any()).SetArg(0, arg).Return(nil),
		sessionMockObj.EXPECT().Close(),
	)

	mgoDial = connectionMockObj

	dbManager := DBManager{}
	res, err := dbManager.GetDockerImage(imageId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if !reflect.DeepEqual(expectedRes, res) {
		t.Errorf("Expected res: %s, actual res: %s", expectedRes, res)
	}
}

func TestCalledGetDockerImageWithInvalidObjectId_ExpectErrorReturn(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	connectionMockObj := mgomocks.NewMockConnection(mockCtrl)
	sessionMockObj := mgomocks.NewMockSession(mockCtrl)

	gomock.InOrder(
		connectionMockObj.EXPECT().Dial(validUrl).Return(sessionMockObj, nil),
		sessionMockObj.EXPECT().Close(),
	)

	mgoDial = connectionMockObj

	dbManager := DBManager{}
	_, err := dbManager.GetDockerImage(invalidObjectId)

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", invalidObjectError.Error(), "nil")
	}

	if err.Error() != invalidObjectError.Error() {
		t.Errorf("Expected err: %s, actual err: %s", invalidObjectError.Error(), err.Error())
	}
}

func TestCalledGetDockerImageWhenDBHasNotMatchedAgent_ExpectErrorReturn(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	query := bson.M{"_id": bson.ObjectIdHex(imageId)}

	connectionMockObj := mgomocks.NewMockConnection(mockCtrl)
	sessionMockObj := mgomocks.NewMockSession(mockCtrl)
	dbMockObj := mgomocks.NewMockDatabase(mockCtrl)
	collectionMockObj := mgomocks.NewMockCollection(mockCtrl)
	queryMockObj := mgomocks.NewMockQuery(mockCtrl)

	gomock.InOrder(
		connectionMockObj.EXPECT().Dial(validUrl).Return(sessionMockObj, nil),
		sessionMockObj.EXPECT().DB(gomock.Any()).Return(dbMockObj),
		dbMockObj.EXPECT().C(gomock.Any()).Return(collectionMockObj),
		collectionMockObj.EXPECT().Find(query).Return(queryMockObj),
		queryMockObj.EXPECT().One(gomock.Any()).Return(mgo.ErrNotFound),
		sessionMockObj.EXPECT().Close(),
	)

	mgoDial = connectionMockObj

	dbManager := DBManager{}
	_, err := dbManager.GetDockerImage(imageId)

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", err.Error())
	case errors.NotFound:
	}
}

func TestCalledDeleteDockerImage_ExpectSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	query := bson.M{"_id": bson.ObjectIdHex(imageId)}

	connectionMockObj := mgomocks.NewMockConnection(mockCtrl)
	sessionMockObj := mgomocks.NewMockSession(mockCtrl)
	dbMockObj := mgomocks.NewMockDatabase(mockCtrl)
	collectionMockObj := mgomocks.NewMockCollection(mockCtrl)

	gomock.InOrder(
		connectionMockObj.EXPECT().Dial(validUrl).Return(sessionMockObj, nil),
		sessionMockObj.EXPECT().DB(gomock.Any()).Return(dbMockObj),
		dbMockObj.EXPECT().C(gomock.Any()).Return(collectionMockObj),
		collectionMockObj.EXPECT().Remove(query).Return(nil),
		sessionMockObj.EXPECT().Close(),
	)

	mgoDial = connectionMockObj

	dbManager := DBManager{}
	err := dbManager.DeleteDockerImage(imageId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}
}

func TestCalledDeleteDockerImageWithInvalidObjectId_ExpectErrorReturn(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	connectionMockObj := mgomocks.NewMockConnection(mockCtrl)
	sessionMockObj := mgomocks.NewMockSession(mockCtrl)

	gomock.InOrder(
		connectionMockObj.EXPECT().Dial(validUrl).Return(sessionMockObj, nil),
		sessionMockObj.EXPECT().Close(),
	)

	mgoDial = connectionMockObj

	dbManager := DBManager{}
	err := dbManager.DeleteDockerImage(invalidObjectId)

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", invalidObjectError.Error(), "nil")
	}

	if err.Error() != invalidObjectError.Error() {
		t.Errorf("Expected err: %s, actual err: %s", invalidObjectError.Error(), err.Error())
	}
}

func TestCalledDeleteDockerImageWhenDBReturnsError_ExpectErrorReturn(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	query := bson.M{"_id": bson.ObjectIdHex(imageId)}

	connectionMockObj := mgomocks.NewMockConnection(mockCtrl)
	sessionMockObj := mgomocks.NewMockSession(mockCtrl)
	dbMockObj := mgomocks.NewMockDatabase(mockCtrl)
	collectionMockObj := mgomocks.NewMockCollection(mockCtrl)

	gomock.InOrder(
		connectionMockObj.EXPECT().Dial(validUrl).Return(sessionMockObj, nil),
		sessionMockObj.EXPECT().DB(gomock.Any()).Return(dbMockObj),
		dbMockObj.EXPECT().C(gomock.Any()).Return(collectionMockObj),
		collectionMockObj.EXPECT().Remove(query).Return(mgo.ErrNotFound),
		sessionMockObj.EXPECT().Close(),
	)

	mgoDial = connectionMockObj

	dbManager := DBManager{}
	err := dbManager.DeleteDockerImage(imageId)

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", err.Error())
	case errors.NotFound:
	}
}
