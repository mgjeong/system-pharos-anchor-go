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

type DeploymentInterface interface {
	// DeployApp request an deployment of edge services to an agent specified by
	// agentId parameter.
	DeployApp(agentId string, body string) (int, map[string]interface{}, error)

	// GetApps request a list of applications that is deployed to an agent specified
	// by agentId parameter.
	GetApps(agentId string) (int, map[string]interface{}, error)

	// GetApp gets the application's information of the agent specified by agentId parameter.
	GetApp(agentId string, appId string) (int, map[string]interface{}, error)

	// UpdateApp request to update an application specified by appId parameter.
	UpdateAppInfo(agentId string, appId string, body string) (int, map[string]interface{}, error)

	// DeleteApp request to delete an application specified by appId parameter.
	DeleteApp(agentId string, appId string) (int, map[string]interface{}, error)

	// UpdateAppInfo request to update all of images which is included an application
	// specified by appId parameter.
	UpdateApp(agentId string, appId string) (int, map[string]interface{}, error)

	// StartApp request to start an application specified by appId parameter.
	StartApp(agentId string, appId string) (int, map[string]interface{}, error)

	// StopApp request to stop an application specified by appId parameter.
	StopApp(agentId string, appId string) (int, map[string]interface{}, error)
}
