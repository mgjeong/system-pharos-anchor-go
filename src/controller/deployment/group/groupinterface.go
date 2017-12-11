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

type DeploymentInterface interface {
	// DeployApp request an deployment of edge services to a group specified by groupId parameter.
	DeployApp(groupId string, body string) (int, map[string]interface{}, error)

	// GetApps request a list of applications that is deployed to a group specified by groupId parameter.
	GetApps(groupId string) (int, map[string]interface{}, error)

	// GetApp gets the application's information of the group specified by groupId parameter.
	GetApp(groupId string, appId string) (int, map[string]interface{}, error)

	// UpdateApp request to update an application specified by appId parameter to all members of the group.
	UpdateAppInfo(groupId string, appId string, body string) (int, map[string]interface{}, error)

	// DeleteApp request to delete an application specified by appId parameter to all members of the group.
	DeleteApp(groupId string, appId string) (int, map[string]interface{}, error)

	// UpdateAppInfo request to update all of images which is included an application specified by
	// appId parameter to all members of the group.
	UpdateApp(groupId string, appId string) (int, map[string]interface{}, error)

	// StartApp request to start an application specified by appId parameter to all members of the group.
	StartApp(groupId string, appId string) (int, map[string]interface{}, error)

	// StopApp request to stop an application specified by appId parameter to all members of the group.
	StopApp(groupId string, appId string) (int, map[string]interface{}, error)
}
