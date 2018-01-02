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

import "net/http"

var SdamAgentHandle SDAMAgentAPIHandlerInterface

var SdamAgent SDAMAgentAPIInterface

type SDAMAgentAPIHandlerInterface interface {
	Handle(w http.ResponseWriter, req *http.Request)
}

type SDAMAgentAPIInterface interface {
	agentRegister(w http.ResponseWriter, req *http.Request)
	agentPing(w http.ResponseWriter, req *http.Request, agentID string)
	agentUnregister(w http.ResponseWriter, req *http.Request, agentID string)
	agent(w http.ResponseWriter, req *http.Request, agentID string)
	agents(w http.ResponseWriter, req *http.Request)
	agentDeployApp(w http.ResponseWriter, req *http.Request, agentID string)
	agentInfoApps(w http.ResponseWriter, req *http.Request, agentID string)
	agentInfoApp(w http.ResponseWriter, req *http.Request, agentID string, appID string)
	agentUpdateAppInfo(w http.ResponseWriter, req *http.Request, agentID string, appID string)
	agentDeleteApp(w http.ResponseWriter, req *http.Request, agentID string, appID string)
	agentStartApp(w http.ResponseWriter, req *http.Request, agentID string, appID string)
	agentStopApp(w http.ResponseWriter, req *http.Request, agentID string, appID string)
	agentUpdateApp(w http.ResponseWriter, req *http.Request, agentID string, appID string)
	agentGetResourceInfo(w http.ResponseWriter, req *http.Request, agentId string)
	agentGetPerformanceInfo(w http.ResponseWriter, req *http.Request, agentId string)
}
