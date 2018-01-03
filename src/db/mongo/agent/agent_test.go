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
	collectionName  = "AGENT"
	status          = "connected"
	appId           = "000000000000000000000000"
	agentId         = "000000000000000000000001"
	invalidObjectId = ""
)

var (
	dummySession       = mgomocks.MockSession{}
	connectionError    = errors.DBConnectionError{}
	invalidObjectError = errors.InvalidObjectId{invalidObjectId}

	configuration = map[string]interface{}{
		"devicename":   "Edge Device #1",
		"deviceid":     "54919CA5-4101-4AE4-595B-353C51AA983C",
		"manufacturer": "Manufacturer Name",
		"modelnumber":  "Model number as designated by the manufacturer",
		"serialnumber": "Serial number",
		"platform":     "Platform name and version",
		"os":           "Operationg system name and version",
		"location":     "Human readable location",
		"pinginterval": "10",
	}
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

func TestCalledAddAgent_ExpectSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	ip := "192.168.0.1"

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

	_, err := Executor.AddAgent(ip, status, configuration)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}
}

func TestCalledAddAgentWhenDBReturnsError_ExpectErrorReturn(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	ip := "192.168.0.1"

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
	_, err := Executor.AddAgent(ip, status, configuration)

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", err.Error())
	case errors.NotFound:
	}
}

func TestCalledUpdateAgentAddress_ExpectSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	query := bson.M{"_id": bson.ObjectIdHex(agentId)}
	update := bson.M{"$set": bson.M{"host": "192.168.0.1", "port": "48098"}}

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
	Executor := Executor{}
	err := Executor.UpdateAgentAddress(agentId, "192.168.0.1", "48098")

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}
}

func TestCalledUpdateAgentAddressWithInvalidObjectId_ExpectErrorReturn(t *testing.T) {
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

	err := Executor.UpdateAgentAddress(invalidObjectId, "192.168.0.1", "48098")

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", invalidObjectError.Error(), "nil")
	}

	if err.Error() != invalidObjectError.Error() {
		t.Errorf("Expected err: %s, actual err: %s", invalidObjectError.Error(), err.Error())
	}
}

func TestCalledUpdateAgentAddressWhenDBReturnsError_ExpectErrorReturn(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	query := bson.M{"_id": bson.ObjectIdHex(agentId)}
	update := bson.M{"$set": bson.M{"host": "192.168.0.1", "port": "48098"}}

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
	Executor := Executor{}
	err := Executor.UpdateAgentAddress(agentId, "192.168.0.1", "48098")

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", err.Error())
	case errors.NotFound:
	}
}

func TestCalledUpdateAgentStatus_ExpectSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	query := bson.M{"_id": bson.ObjectIdHex(agentId)}
	update := bson.M{"$set": bson.M{"status": "connected"}}

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
	Executor := Executor{}
	err := Executor.UpdateAgentStatus(agentId, "connected")

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}
}

func TestCalledUpdateAgentStatusWithInvalidObjectId_ExpectErrorReturn(t *testing.T) {
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
	err := Executor.UpdateAgentStatus(invalidObjectId, "connected")

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", invalidObjectError.Error(), "nil")
	}

	if err.Error() != invalidObjectError.Error() {
		t.Errorf("Expected err: %s, actual err: %s", invalidObjectError.Error(), err.Error())
	}
}

func TestCalledUpdateAgentStatusWhenDBReturnsError_ExpectErrorReturn(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	query := bson.M{"_id": bson.ObjectIdHex(agentId)}
	update := bson.M{"$set": bson.M{"status": "connected"}}

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
	Executor := Executor{}
	err := Executor.UpdateAgentStatus(agentId, "connected")

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", err.Error())
	case errors.NotFound:
	}
}

func TestCalledGetAgent_ExpectSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	query := bson.M{"_id": bson.ObjectIdHex(agentId)}
	arg := Agent{ID: bson.ObjectIdHex(agentId), IP: "192.168.0.1", Apps: []string{}, Status: status, Config: configuration}
	expectedRes := map[string]interface{}{
		"id":     agentId,
		"ip":     "192.168.0.1",
		"apps":   []string{},
		"status": status,
		"config": configuration,
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
	res, err := Executor.GetAgent(agentId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if !reflect.DeepEqual(expectedRes, res) {
		t.Errorf("Expected res: %s, actual res: %s", expectedRes, res)
	}
}

func TestCalledGetAgentWithInvalidObjectId_ExpectErrorReturn(t *testing.T) {
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
	_, err := Executor.GetAgent(invalidObjectId)

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", invalidObjectError.Error(), "nil")
	}

	if err.Error() != invalidObjectError.Error() {
		t.Errorf("Expected err: %s, actual err: %s", invalidObjectError.Error(), err.Error())
	}
}

func TestCalledGetAgentWhenDBHasNotMatchedAgent_ExpectErrorReturn(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	query := bson.M{"_id": bson.ObjectIdHex(agentId)}

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
	_, err := Executor.GetAgent(agentId)

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", err.Error())
	case errors.NotFound:
	}
}

func TestCalledGetAllAgents_ExpectSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	args := []Agent{{ID: bson.ObjectIdHex(agentId), IP: "192.168.0.1", Apps: []string{}, Status: status, Config: configuration}}
	expectedRes := []map[string]interface{}{{
		"id":     agentId,
		"ip":     "192.168.0.1",
		"apps":   []string{},
		"status": status,
		"config": configuration,
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
	res, err := Executor.GetAllAgents()

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if !reflect.DeepEqual(expectedRes, res) {
		t.Errorf("Expected res: %s, actual res: %s", expectedRes, res)
	}
}

func TestCalledGetAllAgentsWhenDBReturnsError_ExpectErrorReturn(t *testing.T) {
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
	_, err := Executor.GetAllAgents()

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", err.Error())
	case errors.NotFound:
	}
}

func TestCalledGetAgentByAppID_ExpectSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	query := bson.M{"_id": bson.ObjectIdHex(agentId), "apps": bson.M{"$in": []string{appId}}}
	arg := Agent{ID: bson.ObjectIdHex(agentId), IP: "192.168.0.1", Apps: []string{}, Status: status, Config: configuration}
	expectedRes := map[string]interface{}{
		"id":     agentId,
		"ip":     "192.168.0.1",
		"apps":   []string{},
		"status": status,
		"config": configuration,
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
	res, err := Executor.GetAgentByAppID(agentId, appId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if !reflect.DeepEqual(expectedRes, res) {
		t.Errorf("Expected res: %s, actual res: %s", expectedRes, res)
	}
}

func TestCalledGetAgentByAppIDWhenDBHasNotMatchedAgent_ExpectErrorReturn(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	query := bson.M{"_id": bson.ObjectIdHex(agentId), "apps": bson.M{"$in": []string{appId}}}

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
	_, err := Executor.GetAgentByAppID(agentId, appId)

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", err.Error())
	case errors.NotFound:
	}
}

func TestCalledGetAgentByAppIDWithInvalidObjectId_ExpectErrorReturn(t *testing.T) {
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
	_, err := Executor.GetAgentByAppID(invalidObjectId, appId)

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", invalidObjectError.Error(), "nil")
	}

	if err.Error() != invalidObjectError.Error() {
		t.Errorf("Expected err: %s, actual err: %s", invalidObjectError.Error(), err.Error())
	}
}

func TestCalledAddAppToAgent_ExpectSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	query := bson.M{"_id": bson.ObjectIdHex(agentId)}
	update := bson.M{"$addToSet": bson.M{"apps": appId}}

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
	Executor := Executor{}
	err := Executor.AddAppToAgent(agentId, appId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}
}

func TestCalledAddAppToAgentWithInvalidObjectId_ExpectErrorReturn(t *testing.T) {
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
	err := Executor.AddAppToAgent(invalidObjectId, appId)

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", invalidObjectError.Error(), "nil")
	}

	if err.Error() != invalidObjectError.Error() {
		t.Errorf("Expected err: %s, actual err: %s", invalidObjectError.Error(), err.Error())
	}
}

func TestCalledAddAppToAgentWhenDBReturnsError_ExpectErrorReturn(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	query := bson.M{"_id": bson.ObjectIdHex(agentId)}
	update := bson.M{"$addToSet": bson.M{"apps": appId}}

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
	Executor := Executor{}
	err := Executor.AddAppToAgent(agentId, appId)

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", err.Error())
	case errors.NotFound:
	}
}

func TestCalledDeleteAppFromAgent_ExpectSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	query := bson.M{"_id": bson.ObjectIdHex(agentId)}
	update := bson.M{"$pull": bson.M{"apps": appId}}

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
	Executor := Executor{}
	err := Executor.DeleteAppFromAgent(agentId, appId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}
}

func TestCalledDeleteAppFromAgentWithInvalidObjectId_ExpectErrorReturn(t *testing.T) {
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
	err := Executor.DeleteAppFromAgent(invalidObjectId, appId)

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", invalidObjectError.Error(), "nil")
	}

	if err.Error() != invalidObjectError.Error() {
		t.Errorf("Expected err: %s, actual err: %s", invalidObjectError.Error(), err.Error())
	}
}

func TestCalledDeleteAppFromAgentWhenDBReturnsError_ExpectErrorReturn(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	query := bson.M{"_id": bson.ObjectIdHex(agentId)}
	update := bson.M{"$pull": bson.M{"apps": appId}}

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
	Executor := Executor{}
	err := Executor.DeleteAppFromAgent(agentId, appId)

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", err.Error())
	case errors.NotFound:
	}
}

func TestCalledDeleteAgent_ExpectSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	query := bson.M{"_id": bson.ObjectIdHex(agentId)}

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
	Executor := Executor{}
	err := Executor.DeleteAgent(agentId)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}
}

func TestCalledDeleteAgentWithInvalidObjectId_ExpectErrorReturn(t *testing.T) {
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
	err := Executor.DeleteAgent(invalidObjectId)

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", invalidObjectError.Error(), "nil")
	}

	if err.Error() != invalidObjectError.Error() {
		t.Errorf("Expected err: %s, actual err: %s", invalidObjectError.Error(), err.Error())
	}
}

func TestCalledDeleteAgentWhenDBReturnsError_ExpectErrorReturn(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	query := bson.M{"_id": bson.ObjectIdHex(agentId)}

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
	Executor := Executor{}
	err := Executor.DeleteAgent(agentId)

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", err.Error())
	case errors.NotFound:
	}
}
