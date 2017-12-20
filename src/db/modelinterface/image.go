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

type ImageInterface interface {
	// AddDockerImage insert a new docker image.
	AddDockerImage(image map[string]interface{}) (map[string]interface{}, error)

	// DeleteDockerImage delete a specific image from db related to image.
	DeleteDockerImage(imageId string) error

	// UpdateDockerImage update status of docker image from db related to image.
	UpdateDockerImage(imageId string, image map[string]interface{}) error

	// GetDockerImage returns a single document from db related to image.
	GetDockerImage(imageId string) (map[string]interface{}, error)
}
