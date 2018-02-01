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

package app

import (
	"commons/logger"
	"commons/results"
)

// Command is an interface of apps operations.
type Command interface {
	Search(query map[string]interface{}) (int, map[string]interface{}, error)
}

// Executor implements the Command interface.
type Executor struct{}

//var appDbExecutor appDB.Command
//
//func init() {
//	appDbExecutor = appDB.Executor{}
//}

func (Executor) Search(query map[string]interface{}) (int, map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	return results.OK, nil, nil
//	return results.OK, res, err
}
