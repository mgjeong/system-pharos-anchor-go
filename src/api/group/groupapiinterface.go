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
package group

import "net/http"

var SdamGroupHandle SDAMGroupAPIHandlerInterface

var SdamGroup SDAMGroupAPIInterface

type SDAMGroupAPIHandlerInterface interface {
	Handle(w http.ResponseWriter, req *http.Request)
}

type SDAMGroupAPIInterface interface {
	createGroup(w http.ResponseWriter, req *http.Request)
	group(w http.ResponseWriter, req *http.Request, groupID string)
	groups(w http.ResponseWriter, req *http.Request)
	groupJoin(w http.ResponseWriter, req *http.Request, groupID string)
	groupLeave(w http.ResponseWriter, req *http.Request, groupID string)
	groupDeployApp(w http.ResponseWriter, req *http.Request, groupID string)
	groupInfoApps(w http.ResponseWriter, req *http.Request, groupID string)
	groupInfoApp(w http.ResponseWriter, req *http.Request, groupID string, appID string)
	groupUpdateAppInfo(w http.ResponseWriter, req *http.Request, groupID string, appID string)
	groupDeleteApp(w http.ResponseWriter, req *http.Request, groupID string, appID string)
	groupStartApp(w http.ResponseWriter, req *http.Request, groupID string, appID string)
	groupStopApp(w http.ResponseWriter, req *http.Request, groupID string, appID string)
	groupUpdateApp(w http.ResponseWriter, req *http.Request, groupID string, appID string)
}
