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

// Package app provides an interfaces to add, delete, get
// an target edge device.
package app

import (
	"commons/logger"
	"commons/results"
	appDB "db/mongo/app"
)

// Command is an interface of app operations.
type Command interface {
	GetAppsWithImageName(imageName string) (int, map[string]interface{}, error)
}

const (
	APPS       = "apps"      // used to indicate a list of apps.
	IMAGE_NAME = "imagename" // used to indicate name of image.
)

// Executor implements the Command interface.
type Executor struct{}

var appDbExecutor appDB.Command

func init() {
	appDbExecutor = appDB.Executor{}
}

func (Executor) GetAppsWithImageName(imageName string) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	query := make(map[string]interface{})
	query[IMAGE_NAME] = imageName

	// Get matched apps with query stored in the database.
	apps, err := appDbExecutor.GetApps(query)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return results.ERROR, nil, err
	}

	res := make(map[string]interface{})
	res[APPS] = apps

	return results.OK, res, err
}
