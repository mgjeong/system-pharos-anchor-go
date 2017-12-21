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
package modelinterface

type RegistryInterface interface {
	// AddDockerRegistry insert a new docker registry information.
	AddDockerRegistry(url string) (map[string]interface{}, error)

	// GetDockerRegistries returns all documents from db related to docker registry.
	GetDockerRegistries() ([]map[string]interface{}, error)

	// GetDockerRegistry returns a single document from db related to docker registry.
	GetDockerRegistry(url string) (map[string]interface{}, error)

	// DeleteDockerRegistry delete a specific docker registry information from db related to registry.
	DeleteDockerRegistry(registryId string) error

	// AddDockerImages add a specific docker image to the target registry.
	AddDockerImages(registryId string, images []map[string]interface{}) error

	// GetDockerImages returns all docker images which belong to the target registry.
	GetDockerImages(registryId string) ([]map[string]interface{}, error)

	// UpdateDockerImage update status of docker image which belong to the target registry.
	UpdateDockerImage(registryId string, image map[string]interface{}) error

	// DeleteDockerImage delete a specific docker image from the target registry.
	DeleteDockerImage(registryId string, image map[string]interface{})
}
