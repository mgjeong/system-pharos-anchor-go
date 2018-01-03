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
	errors "commons/errors"
	imagedbmocks "db/mongo/image/mocks"
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
	collectionName  = "REGISTRY"
	registryUrl     = "http://localhost:5000"
	imageId         = "000000000000000000000000"
	registryId      = "000000000000000000000001"
	invalidObjectId = ""
)

var (
	image = map[string]interface{}{
		"id":         imageId,
		"repository": "http://localhost:5000",
		"tag":        "latest",
		"size":       "3453",
		"action":     "push",
		"timestamp":  "2017-11-27T12:42:56.914783506Z",
	}
	images             = []map[string]interface{}{image}
	dummySession       = mgomocks.MockSession{}
	connectionError    = errors.DBConnectionError{}
	invalidObjectError = errors.InvalidObjectId{invalidObjectId}
)

func TestCalledConnectWithEmptyURL_ExpectErrorReturn(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	invalidUrl := ""

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

func TestCalledAddDockerRegistry_ExpectSuccess(t *testing.T) {
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
	Executor := Executor{}

	_, err := Executor.AddDockerRegistry(registryUrl)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}
}

func TestCalledAddDockerRegistryWhenDBReturnsError_ExpectErrorReturn(t *testing.T) {
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
	Executor := Executor{}

	_, err := Executor.AddDockerRegistry(registryUrl)

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", err.Error())
	case errors.NotFound:
	}
}

func TestCalledGetDockerRegistries_ExpectSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	args := []Registry{{ID: bson.ObjectIdHex(registryId), Url: registryUrl, Images: []string{}}}
	expectedRes := []map[string]interface{}{{
		"id":  registryId,
		"url": registryUrl,
	}}

	connectionMockObj := mgomocks.NewMockConnection(mockCtrl)
	sessionMockObj := mgomocks.NewMockSession(mockCtrl)
	dbMockObj := mgomocks.NewMockDatabase(mockCtrl)
	collectionMockObj := mgomocks.NewMockCollection(mockCtrl)
	queryMockObj := mgomocks.NewMockQuery(mockCtrl)

	gomock.InOrder(
		connectionMockObj.EXPECT().Dial(validUrl).Return(sessionMockObj, nil),
		sessionMockObj.EXPECT().DB(gomock.Any()).Return(dbMockObj),
		dbMockObj.EXPECT().C(gomock.Any()).Return(collectionMockObj),
		collectionMockObj.EXPECT().Find(nil).Return(queryMockObj),
		queryMockObj.EXPECT().All(gomock.Any()).SetArg(0, args).Return(nil),
		sessionMockObj.EXPECT().Close(),
	)

	mgoDial = connectionMockObj

	Executor := Executor{}
	res, err := Executor.GetDockerRegistries()

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if !reflect.DeepEqual(expectedRes, res) {
		t.Errorf("Expected res: %s, actual res: %s", expectedRes, res)
	}
}

func TestCalledGetDockerRegistriesWhenDBReturnsError_ExpectErrorReturn(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	connectionMockObj := mgomocks.NewMockConnection(mockCtrl)
	sessionMockObj := mgomocks.NewMockSession(mockCtrl)
	dbMockObj := mgomocks.NewMockDatabase(mockCtrl)
	collectionMockObj := mgomocks.NewMockCollection(mockCtrl)
	queryMockObj := mgomocks.NewMockQuery(mockCtrl)

	gomock.InOrder(
		connectionMockObj.EXPECT().Dial(validUrl).Return(sessionMockObj, nil),
		sessionMockObj.EXPECT().DB(gomock.Any()).Return(dbMockObj),
		dbMockObj.EXPECT().C(gomock.Any()).Return(collectionMockObj),
		collectionMockObj.EXPECT().Find(nil).Return(queryMockObj),
		queryMockObj.EXPECT().All(gomock.Any()).Return(mgo.ErrNotFound),
		sessionMockObj.EXPECT().Close(),
	)

	mgoDial = connectionMockObj

	Executor := Executor{}
	_, err := Executor.GetDockerRegistries()

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", err.Error())
	case errors.NotFound:
	}
}

func TestCalledGetDockerRegistry_ExpectSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	query := bson.M{"url": registryUrl}
	arg := Registry{ID: bson.ObjectIdHex(registryId), Url: registryUrl}
	expectedRes := map[string]interface{}{
		"id":  registryId,
		"url": registryUrl,
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

	Executor := Executor{}
	res, err := Executor.GetDockerRegistry(registryUrl)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if !reflect.DeepEqual(expectedRes, res) {
		t.Errorf("Expected res: %s, actual res: %s", expectedRes, res)
	}
}

func TestCalledGetDockerRegistryWhenDBHasNotMatchedAgent_ExpectErrorReturn(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	query := bson.M{"url": registryUrl}

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

	Executor := Executor{}
	_, err := Executor.GetDockerRegistry(registryUrl)

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", err.Error())
	case errors.NotFound:
	}
}

func TestCalledDeleteDockerRegistry_ExpectSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	query := bson.M{"_id": bson.ObjectIdHex(registryId)}
	arg := Registry{ID: bson.ObjectIdHex(registryId), Url: registryUrl, Images: []string{}}

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
		sessionMockObj.EXPECT().DB(gomock.Any()).Return(dbMockObj),
		dbMockObj.EXPECT().C(gomock.Any()).Return(collectionMockObj),
		collectionMockObj.EXPECT().Remove(query).Return(nil),
		sessionMockObj.EXPECT().Close(),
	)

	mgoDial = connectionMockObj

	Executor := Executor{}
	err := Executor.DeleteDockerRegistry(registryId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}
}

func TestCalledDeleteDockerRegistryWithInvalidObjectId_ExpectErrorReturn(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	connectionMockObj := mgomocks.NewMockConnection(mockCtrl)
	sessionMockObj := mgomocks.NewMockSession(mockCtrl)

	gomock.InOrder(
		connectionMockObj.EXPECT().Dial(validUrl).Return(sessionMockObj, nil),
		sessionMockObj.EXPECT().Close(),
	)

	mgoDial = connectionMockObj

	Executor := Executor{}
	err := Executor.DeleteDockerRegistry(invalidObjectId)

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", invalidObjectError.Error(), "nil")
	}

	if err.Error() != invalidObjectError.Error() {
		t.Errorf("Expected err: %s, actual err: %s", invalidObjectError.Error(), err.Error())
	}
}

func TestCalledDeleteDockerRegistryWhenDBReturnsError_ExpectErrorReturn(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	query := bson.M{"_id": bson.ObjectIdHex(registryId)}

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

	Executor := Executor{}
	err := Executor.DeleteDockerRegistry(registryId)

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", err.Error())
	case errors.NotFound:
	}
}

func TestCalledAddDockerImages_ExpectSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	query := bson.M{"_id": bson.ObjectIdHex(registryId)}
	update := bson.M{"$addToSet": bson.M{"images": imageId}}
	arg := Registry{ID: bson.ObjectIdHex(registryId), Url: registryUrl, Images: []string{}}

	connectionMockObj := mgomocks.NewMockConnection(mockCtrl)
	sessionMockObj := mgomocks.NewMockSession(mockCtrl)
	dbMockObj := mgomocks.NewMockDatabase(mockCtrl)
	collectionMockObj := mgomocks.NewMockCollection(mockCtrl)
	queryMockObj := mgomocks.NewMockQuery(mockCtrl)
	imageDBMockObj := imagedbmocks.NewMockCommand(mockCtrl)

	gomock.InOrder(
		connectionMockObj.EXPECT().Dial(validUrl).Return(sessionMockObj, nil),
		sessionMockObj.EXPECT().DB(gomock.Any()).Return(dbMockObj),
		dbMockObj.EXPECT().C(gomock.Any()).Return(collectionMockObj),
		collectionMockObj.EXPECT().Find(query).Return(queryMockObj),
		queryMockObj.EXPECT().One(gomock.Any()).SetArg(0, arg).Return(nil),
		imageDBMockObj.EXPECT().AddDockerImage(image).Return(image, nil),
		sessionMockObj.EXPECT().DB(gomock.Any()).Return(dbMockObj),
		dbMockObj.EXPECT().C(gomock.Any()).Return(collectionMockObj),
		collectionMockObj.EXPECT().Update(query, update).Return(nil),
		sessionMockObj.EXPECT().Close(),
	)

	mgoDial = connectionMockObj
	imageExecutor = imageDBMockObj

	Executor := Executor{}
	err := Executor.AddDockerImages(registryId, images)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}
}

func TestCalledAddDockerImagesWithInvalidObjectId_ExpectErrorReturn(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	connectionMockObj := mgomocks.NewMockConnection(mockCtrl)
	sessionMockObj := mgomocks.NewMockSession(mockCtrl)

	gomock.InOrder(
		connectionMockObj.EXPECT().Dial(validUrl).Return(sessionMockObj, nil),
		sessionMockObj.EXPECT().Close(),
	)

	mgoDial = connectionMockObj

	Executor := Executor{}
	err := Executor.AddDockerImages(invalidObjectId, images)

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", invalidObjectError.Error(), "nil")
	}

	if err.Error() != invalidObjectError.Error() {
		t.Errorf("Expected err: %s, actual err: %s", invalidObjectError.Error(), err.Error())
	}
}

func TestCalledAddDockerImagesWhenDBHasNotMatchedRegistry_ExpectErrorReturn(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	query := bson.M{"_id": bson.ObjectIdHex(registryId)}

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

	Executor := Executor{}
	err := Executor.AddDockerImages(registryId, images)

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", err.Error())
	case errors.NotFound:
	}
}

func TestCalledGetDockerImages_ExpectSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	query := bson.M{"_id": bson.ObjectIdHex(registryId)}
	arg := Registry{ID: bson.ObjectIdHex(registryId), Url: registryUrl, Images: []string{imageId}}
	expectedRes := []map[string]interface{}{{
		"repository": image["repository"].(string),
		"tag":        image["tag"].(string),
		"size":       image["size"].(string),
		"action":     image["action"].(string),
		"timestamp":  image["timestamp"].(string),
	}}

	connectionMockObj := mgomocks.NewMockConnection(mockCtrl)
	sessionMockObj := mgomocks.NewMockSession(mockCtrl)
	dbMockObj := mgomocks.NewMockDatabase(mockCtrl)
	collectionMockObj := mgomocks.NewMockCollection(mockCtrl)
	queryMockObj := mgomocks.NewMockQuery(mockCtrl)
	imageDBMockObj := imagedbmocks.NewMockCommand(mockCtrl)

	gomock.InOrder(
		connectionMockObj.EXPECT().Dial(validUrl).Return(sessionMockObj, nil),
		sessionMockObj.EXPECT().DB(gomock.Any()).Return(dbMockObj),
		dbMockObj.EXPECT().C(gomock.Any()).Return(collectionMockObj),
		collectionMockObj.EXPECT().Find(query).Return(queryMockObj),
		queryMockObj.EXPECT().One(gomock.Any()).SetArg(0, arg).Return(nil),
		imageDBMockObj.EXPECT().GetDockerImage(imageId).Return(image, nil),
		sessionMockObj.EXPECT().Close(),
	)

	mgoDial = connectionMockObj
	imageExecutor = imageDBMockObj

	Executor := Executor{}
	res, err := Executor.GetDockerImages(registryId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if !reflect.DeepEqual(expectedRes, res) {
		t.Errorf("Expected res: %s, actual res: %s", expectedRes, res)
	}
}

func TestCalledGetDockerImagesWithInvalidObjectId_ExpectErrorReturn(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	connectionMockObj := mgomocks.NewMockConnection(mockCtrl)
	sessionMockObj := mgomocks.NewMockSession(mockCtrl)

	gomock.InOrder(
		connectionMockObj.EXPECT().Dial(validUrl).Return(sessionMockObj, nil),
		sessionMockObj.EXPECT().Close(),
	)

	mgoDial = connectionMockObj

	Executor := Executor{}
	_, err := Executor.GetDockerImages(invalidObjectId)

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", invalidObjectError.Error(), "nil")
	}

	if err.Error() != invalidObjectError.Error() {
		t.Errorf("Expected err: %s, actual err: %s", invalidObjectError.Error(), err.Error())
	}
}

func TestCalledGetDockerImagesWhenDBHasNotMatchedRegistry_ExpectErrorReturn(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	query := bson.M{"_id": bson.ObjectIdHex(registryId)}

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

	Executor := Executor{}
	_, err := Executor.GetDockerImages(registryId)

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

	query := bson.M{"_id": bson.ObjectIdHex(registryId)}
	arg := Registry{ID: bson.ObjectIdHex(registryId), Url: registryUrl, Images: []string{imageId}}

	connectionMockObj := mgomocks.NewMockConnection(mockCtrl)
	sessionMockObj := mgomocks.NewMockSession(mockCtrl)
	dbMockObj := mgomocks.NewMockDatabase(mockCtrl)
	collectionMockObj := mgomocks.NewMockCollection(mockCtrl)
	queryMockObj := mgomocks.NewMockQuery(mockCtrl)
	imageDBMockObj := imagedbmocks.NewMockCommand(mockCtrl)

	gomock.InOrder(
		connectionMockObj.EXPECT().Dial(validUrl).Return(sessionMockObj, nil),
		sessionMockObj.EXPECT().DB(gomock.Any()).Return(dbMockObj),
		dbMockObj.EXPECT().C(gomock.Any()).Return(collectionMockObj),
		collectionMockObj.EXPECT().Find(query).Return(queryMockObj),
		queryMockObj.EXPECT().One(gomock.Any()).SetArg(0, arg).Return(nil),
		imageDBMockObj.EXPECT().GetDockerImage(imageId).Return(image, nil),
		imageDBMockObj.EXPECT().UpdateDockerImage(imageId, image).Return(nil),
		sessionMockObj.EXPECT().Close(),
	)

	mgoDial = connectionMockObj
	imageExecutor = imageDBMockObj

	Executor := Executor{}
	err := Executor.UpdateDockerImage(registryId, image)

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

	Executor := Executor{}
	err := Executor.UpdateDockerImage(invalidObjectId, image)

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", invalidObjectError.Error(), "nil")
	}

	if err.Error() != invalidObjectError.Error() {
		t.Errorf("Expected err: %s, actual err: %s", invalidObjectError.Error(), err.Error())
	}
}

func TestCalledUpdateDockerImageWhenDBHasNotMatchedRegistry_ExpectErrorReturn(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	query := bson.M{"_id": bson.ObjectIdHex(registryId)}

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

	Executor := Executor{}
	err := Executor.UpdateDockerImage(registryId, image)

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", err.Error())
	case errors.NotFound:
	}
}

func TestCalledUpdateDockerImageWhenRegistryHasNotMatchedImage_ExpectErrorReturn(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	query := bson.M{"_id": bson.ObjectIdHex(registryId)}
	arg := Registry{ID: bson.ObjectIdHex(registryId), Url: registryUrl, Images: []string{}}

	connectionMockObj := mgomocks.NewMockConnection(mockCtrl)
	sessionMockObj := mgomocks.NewMockSession(mockCtrl)
	dbMockObj := mgomocks.NewMockDatabase(mockCtrl)
	collectionMockObj := mgomocks.NewMockCollection(mockCtrl)
	queryMockObj := mgomocks.NewMockQuery(mockCtrl)
	imageDBMockObj := imagedbmocks.NewMockCommand(mockCtrl)

	gomock.InOrder(
		connectionMockObj.EXPECT().Dial(validUrl).Return(sessionMockObj, nil),
		sessionMockObj.EXPECT().DB(gomock.Any()).Return(dbMockObj),
		dbMockObj.EXPECT().C(gomock.Any()).Return(collectionMockObj),
		collectionMockObj.EXPECT().Find(query).Return(queryMockObj),
		queryMockObj.EXPECT().One(gomock.Any()).SetArg(0, arg).Return(nil),
		sessionMockObj.EXPECT().Close(),
	)

	mgoDial = connectionMockObj
	imageExecutor = imageDBMockObj

	Executor := Executor{}
	err := Executor.UpdateDockerImage(registryId, image)

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

	query := bson.M{"_id": bson.ObjectIdHex(registryId)}
	update := bson.M{"$pull": bson.M{"images": imageId}}
	arg := Registry{ID: bson.ObjectIdHex(registryId), Url: registryUrl, Images: []string{imageId}}

	connectionMockObj := mgomocks.NewMockConnection(mockCtrl)
	sessionMockObj := mgomocks.NewMockSession(mockCtrl)
	dbMockObj := mgomocks.NewMockDatabase(mockCtrl)
	collectionMockObj := mgomocks.NewMockCollection(mockCtrl)
	queryMockObj := mgomocks.NewMockQuery(mockCtrl)
	imageDBMockObj := imagedbmocks.NewMockCommand(mockCtrl)

	gomock.InOrder(
		connectionMockObj.EXPECT().Dial(validUrl).Return(sessionMockObj, nil),
		sessionMockObj.EXPECT().DB(gomock.Any()).Return(dbMockObj),
		dbMockObj.EXPECT().C(gomock.Any()).Return(collectionMockObj),
		collectionMockObj.EXPECT().Find(query).Return(queryMockObj),
		queryMockObj.EXPECT().One(gomock.Any()).SetArg(0, arg).Return(nil),
		imageDBMockObj.EXPECT().GetDockerImage(imageId).Return(image, nil),
		sessionMockObj.EXPECT().DB(gomock.Any()).Return(dbMockObj),
		dbMockObj.EXPECT().C(gomock.Any()).Return(collectionMockObj),
		collectionMockObj.EXPECT().Update(query, update).Return(nil),
		imageDBMockObj.EXPECT().DeleteDockerImage(imageId).Return(nil),
		sessionMockObj.EXPECT().Close(),
	)

	mgoDial = connectionMockObj
	imageExecutor = imageDBMockObj

	Executor := Executor{}
	err := Executor.DeleteDockerImage(registryId, image)

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

	Executor := Executor{}
	err := Executor.DeleteDockerImage(invalidObjectId, image)

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", invalidObjectError.Error(), "nil")
	}

	if err.Error() != invalidObjectError.Error() {
		t.Errorf("Expected err: %s, actual err: %s", invalidObjectError.Error(), err.Error())
	}
}

func TestCalledDeleteDockerImageWhenDBHasNotMatchedRegistry_ExpectErrorReturn(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	query := bson.M{"_id": bson.ObjectIdHex(registryId)}

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

	Executor := Executor{}
	err := Executor.DeleteDockerImage(registryId, image)

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", err.Error())
	case errors.NotFound:
	}
}

func TestCalledDeleteDockerImageWhenRegistryHasNotMatchedImage_ExpectErrorReturn(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	query := bson.M{"_id": bson.ObjectIdHex(registryId)}
	arg := Registry{ID: bson.ObjectIdHex(registryId), Url: registryUrl, Images: []string{}}

	connectionMockObj := mgomocks.NewMockConnection(mockCtrl)
	sessionMockObj := mgomocks.NewMockSession(mockCtrl)
	dbMockObj := mgomocks.NewMockDatabase(mockCtrl)
	collectionMockObj := mgomocks.NewMockCollection(mockCtrl)
	queryMockObj := mgomocks.NewMockQuery(mockCtrl)
	imageDBMockObj := imagedbmocks.NewMockCommand(mockCtrl)

	gomock.InOrder(
		connectionMockObj.EXPECT().Dial(validUrl).Return(sessionMockObj, nil),
		sessionMockObj.EXPECT().DB(gomock.Any()).Return(dbMockObj),
		dbMockObj.EXPECT().C(gomock.Any()).Return(collectionMockObj),
		collectionMockObj.EXPECT().Find(query).Return(queryMockObj),
		queryMockObj.EXPECT().One(gomock.Any()).SetArg(0, arg).Return(nil),
		sessionMockObj.EXPECT().Close(),
	)

	mgoDial = connectionMockObj
	imageExecutor = imageDBMockObj

	Executor := Executor{}
	err := Executor.DeleteDockerImage(registryId, image)

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", err.Error())
	case errors.NotFound:
	}
}
