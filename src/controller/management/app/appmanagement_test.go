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
package app

import (
	"commons/errors"
	"commons/results"
	dbmocks "db/mongo/app/mocks"
	"github.com/golang/mock/gomock"
	"testing"
)

const (
	IMAGENAME         = "testImageName"
	INVALID_IMAGENAME = "invalidImageName"
	APPID             = "000000000000000000000000"
	REFCNT            = 1
)

var (
	app = map[string]interface{}{
		"id":       APPID,
		"images":   []string{},
		"services": []string{},
		"refcnt":   REFCNT,
	}
	notFoundError = errors.NotFound{}
)

var manager Command

func init() {
	manager = Executor{}
}

func TestCalledGetAppsWithImageName_ExpectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apps := []map[string]interface{}{app}

	query := make(map[string]interface{})
	query[IMAGES] = IMAGENAME

	dbExecutorMockObj := dbmocks.NewMockCommand(ctrl)

	gomock.InOrder(
		dbExecutorMockObj.EXPECT().GetApps(query).Return(apps, nil),
	)
	// pass mockObj to a real object.
	appDbExecutor = dbExecutorMockObj

	code, _, err := manager.GetAppsWithImageName(IMAGENAME)

	if err != nil {
		t.Errorf("Unexpected err: %s", err.Error())
	}

	if code != results.OK {
		t.Errorf("Expected code: %d, actual code: %d", results.OK, code)
	}
}

func TestCalledGetAppsWithInvalidImageName_ExpectErrorReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbExecutorMockObj := dbmocks.NewMockCommand(ctrl)

	query := make(map[string]interface{})
	query[IMAGES] = INVALID_IMAGENAME

	gomock.InOrder(
		dbExecutorMockObj.EXPECT().GetApps(query).Return(nil, notFoundError),
	)

	// pass mockObj to a real object.
	appDbExecutor = dbExecutorMockObj

	code, _, err := manager.GetAppsWithImageName(INVALID_IMAGENAME)

	if code != results.ERROR {
		t.Errorf("Expected code: %d, actual code: %d", results.ERROR, code)
	}

	if err == nil {
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", "nil")
	}

	switch err.(type) {
	default:
		t.Errorf("Expected err: %s, actual err: %s", "NotFound", err.Error())
	case errors.NotFound:
	}
}
